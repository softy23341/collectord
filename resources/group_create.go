package resource

import (
	"errors"
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/groups"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// CreateGroup TBD
func (g *Group) CreateGroup(params groups.PostGroupsParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)

	var (
		logger = userContext.Logger(params.HTTPRequest)
		DBM    = g.Context.DBM

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return groups.NewPostGroupsDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		forbiddenResponse = groups.NewPostGroupsForbidden

		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	)

	// check access rights
	ok, err := accessChecker.IsUserRelatedToRoots(userContext.User.ID, []int64{*params.RNewGroup.RootID})
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root ids: %d, %d", userContext.User.ID, *params.RNewGroup.RootID)
		logger.Error("user cant create entity here", "err", err)
		return forbiddenResponse()
	}

	// collections
	inputCollectionsIDs := params.RNewGroup.CollectionsIds
	groupCollectionsList, err := DBM.GetCollectionsByIDs(inputCollectionsIDs)
	if err != nil {
		logger.Error("GetCollectionsByIDs", "error", err)
		return errorResponse(500, err.Error())
	} else if len(groupCollectionsList) != len(inputCollectionsIDs) {
		err := errors.New("cant find all collections")
		logger.Error("cant find all collections", "error", err)
		return errorResponse(500, err.Error())
	} else if !util.Int64Slice(groupCollectionsList.RootsIDs()).All(*params.RNewGroup.RootID) {
		err := errors.New("cant mix collections form different roots")
		logger.Error("cant find all collections", "error", err)
		return errorResponse(500, err.Error())
	}

	ok, err = accessChecker.HasUserRightsForCollections(
		userContext.User.ID,
		dto.RightEntityLevelWrite,
		groupCollectionsList.GetIDs(),
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user cant touch collections: %d => %v",
			userContext.User.ID,
			groupCollectionsList.GetIDs(),
		)
		logger.Error("user cant be here", "err", err)
		return forbiddenResponse()
	}

	// validate on unique root_id & name
	if group, err := DBM.GetGroupByRootIDAndName(*params.RNewGroup.RootID, *params.RNewGroup.Name); err != nil {
		logger.Error("GetGroupByRootIDAndName", "error", err)
		return errorResponse(500, err.Error())
	} else if group != nil {
		err := fmt.Errorf("group '%s' already exist with id: %d", *params.RNewGroup.Name, group.ID)
		logger.Error("group already exist", "group", group.ID)
		return errorResponse(422, err.Error())
	}

	newGroup := &dto.Group{
		Name:       *params.RNewGroup.Name,
		RootID:     *params.RNewGroup.RootID,
		UserUniqID: params.RNewGroup.ClientUniqID,
		UserID:     &userContext.User.ID,
	}

	tx, err := g.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begin tx", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.CreateGroup(newGroup); err != nil {
		logger.Error("create group", "err", err)
		return errorResponse(500, err.Error())
	}

	if len(inputCollectionsIDs) != 0 {
		if err := tx.CreateGroupCollectionsRefs(newGroup.ID, inputCollectionsIDs); err != nil {
			logger.Error("create group-collections refs", "err", err)
			return errorResponse(500, err.Error())
		}
	}

	newGroupEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeNewEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			NewEntity: &dto.EntityRef{
				Typo: dto.GroupEntityType,
				ID:   newGroup.ID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, newGroupEvent); err != nil {
		logger.Error("new group event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "err", err)
		return errorResponse(500, err.Error())
	}

	g.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		g.Context.EventSender.Send(newGroupEvent)
	}))

	return groups.NewPostGroupsOK().WithPayload(&models.ANewGroup{
		ID: &newGroup.ID,
	})
}
