package resource

import (
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/groups"
	"git.softndit.com/collector/backend/services"
	"github.com/go-openapi/runtime/middleware"
)

// DeleteGroup TBD
func (g *Group) DeleteGroup(params groups.DeleteGroupsIDDeleteParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)

		logger = userContext.Logger(params.HTTPRequest)
		DBM    = g.Context.DBM

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return groups.NewDeleteGroupsIDDeleteDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}

		notFoundResponse = func() middleware.Responder {
			return groups.NewDeleteGroupsIDDeleteNotFound()
		}
		forbiddenResponse = groups.NewDeleteGroupsIDDeleteForbidden
		accessChecker     = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	)

	// find group and check rights
	var groupToDelete *dto.Group
	if groupList, err := DBM.GetGroupsByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant find group", "err", err)
		return errorResponse(500, err.Error())
	} else if len(groupList) != 1 {
		err := fmt.Errorf("cant find group with id: %d", params.ID)
		logger.Error("cant find group", "err", err)
		return notFoundResponse()
	} else {
		groupToDelete = groupList[0]
	}

	// check rights
	ok, err := accessChecker.HasUserRightsForGroups(
		userContext.User.ID,
		dto.RightEntityLevelWrite,
		[]int64{groupToDelete.ID},
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant delete group %d", groupToDelete.ID)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin tx", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// delete group (refs will be deleted cascade)
	if err := tx.DeleteGroups([]int64{groupToDelete.ID}); err != nil {
		logger.Error("cant delete group", "err", err)
		return errorResponse(500, err.Error())
	}

	deletedGroupEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeDeletedEntity,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			DeletedEntity: &dto.EntityRef{
				Typo: dto.GroupEntityType,
				ID:   groupToDelete.ID,
			},
		},
	}
	if _, err := services.EmplaceEvent(tx, deletedGroupEvent); err != nil {
		logger.Error("delete group event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "err", err)
		return errorResponse(500, err.Error())
	}

	g.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		g.Context.EventSender.Send(deletedGroupEvent)
	}))

	return groups.NewDeleteGroupsIDDeleteNoContent()
}
