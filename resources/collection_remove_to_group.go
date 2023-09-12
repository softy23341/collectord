package resource

import (
	"errors"
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/collections"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// RemoveCollectionsFromGroup TBD
func (c *Collection) RemoveCollectionsFromGroup(params collections.PostCollectionsRemoveFromGroupParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)
		DBM         = c.Context.DBM

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return collections.NewPostCollectionsAddToGroupDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		forbiddenResponse = collections.NewPostCollectionsAddToGroupForbidden

		moveParams = params.RRemoveCollectionsFromeGroup

		collectionsIDs = moveParams.CollectionsIds
		groupID        = *moveParams.TargetGroupID

		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	)
	logger.Debug("rm collections")

	// find group
	var updatedGroup *dto.Group
	if groupList, err := DBM.GetGroupsByIDs([]int64{groupID}); err != nil {
		logger.Error("cant find group", "err", err)
		return errorResponse(500, err.Error())
	} else if len(groupList) != 1 {
		err := fmt.Errorf("cant find group with id: %d", groupID)
		logger.Error("cant find group", "err", err)
		return errorResponse(500, err.Error())
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

	// check collections
	groupCollectionsList, err := DBM.GetCollectionsByIDs(collectionsIDs)
	if err != nil {
		logger.Error("GetCollectionsByIDs", "error", err)
		return errorResponse(500, err.Error())
	} else if len(groupCollectionsList) != len(collectionsIDs) {
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
		dto.RightEntityLevelRead,
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

	tx, err := c.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("cant begin tx", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.DeleteCollectionsGroupRefsByGroupAndCollections(moveParams.CollectionsIds, *moveParams.TargetGroupID); err != nil {
		logger.Error("cant delete collections from group", "err", err)
		return errorResponse(500, err.Error())
	}

	var events dto.EventList
	for _, collectionID := range moveParams.CollectionsIds {
		removedCollectionEvent := &dto.Event{
			UserID:       userContext.User.ID,
			Type:         dto.EventTypeRemovedFromGroupCollection,
			CreationTime: time.Now(),
			EventUnion: dto.EventUnion{
				RemovedFromGroupCollection: &dto.EventRemovedFromGroupCollection{
					GroupID:      *moveParams.TargetGroupID,
					CollectionID: collectionID,
				},
			},
		}

		if _, err := services.EmplaceEvent(tx, removedCollectionEvent); err != nil {
			logger.Error("removed collection event", "error", err)
			return errorResponse(500, err.Error())
		}

		events = append(events, removedCollectionEvent)
	}

	if err := tx.Commit(); err != nil {
		logger.Error("cant commit", "err", err)
		return errorResponse(500, err.Error())
	}

	c.Context.EventSender.Send(events...)

	return collections.NewPostCollectionsRemoveFromGroupNoContent()

}
