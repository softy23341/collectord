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

// UpdateGroup TBD
func (g *Group) UpdateGroup(params groups.PutGroupsIDUpdateParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)
		DBM         = g.Context.DBM

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return groups.NewPutGroupsIDUpdateDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		notFoundResponse = func() middleware.Responder {
			return groups.NewPutGroupsIDUpdateNotFound()
		}
		forbiddenResponse = groups.NewPutGroupsIDUpdateForbidden

		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	)

	// find group
	var updatedGroup *dto.Group
	if groupList, err := g.Context.DBM.GetGroupsByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant find group", "err", err)
		return errorResponse(500, err.Error())
	} else if len(groupList) != 1 {
		err := fmt.Errorf("cant find group with id: %d", params.ID)
		logger.Error("cant find group", "err", err)
		return notFoundResponse()
	} else {
		updatedGroup = groupList[0]
	}

	// check group
	ok, err := accessChecker.HasUserRightsForGroups(
		userContext.User.ID,
		dto.RightEntityLevelWrite,
		[]int64{updatedGroup.ID},
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant update group %d", updatedGroup.ID)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	// updates
	groupParams := params.RUpdateGroup

	// check collections
	inputCollectionsIDs := groupParams.CollectionsIds
	groupCollectionsList, err := DBM.GetCollectionsByIDs(inputCollectionsIDs)
	if err != nil {
		logger.Error("GetCollectionsByIDs", "error", err)
		return errorResponse(500, err.Error())
	} else if len(groupCollectionsList) != len(inputCollectionsIDs) {
		err := errors.New("cant find all collections")
		logger.Error("cant find all collections", "error", err)
		return errorResponse(500, err.Error())
	} else if !util.Int64Slice(groupCollectionsList.RootsIDs()).All(updatedGroup.RootID) {
		err := errors.New("cant mix root entities")
		logger.Error("cant use these collections", "error", err)
		return forbiddenResponse()
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

	// set name
	if newName := groupParams.Name; newName != "" {
		// validate on unique root_id & name
		if group, err := DBM.GetGroupByRootIDAndName(updatedGroup.RootID, newName); err != nil {
			logger.Error("GetGroupByRootIDAndName", "error", err)
			return errorResponse(500, err.Error())
		} else if group != nil {
			err := fmt.Errorf("group '%s' already exist with id: %d", newName, group.ID)
			logger.Error("group already exist", "group", group.ID)
			return errorResponse(422, err.Error())
		}

		updatedGroup.Name = newName
	}

	tx, err := g.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("cant start transaction", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// process collections refs
	if inputCollectionsIDs != nil {
		logger.Error("collections array is deprecated")
		return forbiddenResponse()

		if err := tx.DeleteCollectionsGroupRefs(updatedGroup.ID); err != nil {
			logger.Error("cant delete group refers", "err", err)
			return errorResponse(500, err.Error())
		}

		if err := tx.CreateGroupCollectionsRefs(updatedGroup.ID, inputCollectionsIDs); err != nil {
			logger.Error("create group-collections refs", "err", err)
			return errorResponse(500, err.Error())
		}
	}

	// update group
	if err := tx.UpdateGroup(updatedGroup); err != nil {
		logger.Error("cant update group", "err", err)
		return errorResponse(500, err.Error())
	}

	editedGroupEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeEditedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			EditedEntity: &dto.EntityRef{
				Typo: dto.GroupEntityType,
				ID:   updatedGroup.ID,
			},
		},
	}
	if _, err := services.EmplaceEvent(tx, editedGroupEvent); err != nil {
		logger.Error("delete group event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("cant commit", "err", err)
		return errorResponse(500, err.Error())
	}

	g.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		g.Context.EventSender.Send(editedGroupEvent)
	}))

	return groups.NewPutGroupsIDUpdateNoContent()
}
