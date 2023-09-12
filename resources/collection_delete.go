package resource

import (
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/collections"
	"git.softndit.com/collector/backend/services"
	"github.com/go-openapi/runtime/middleware"
)

// DeleteCollection TBD
func (c *Collection) DeleteCollection(params collections.DeleteCollectionsIDDeleteParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)

		logger = userContext.Logger(params.HTTPRequest)
		DBM    = c.Context.DBM

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return collections.NewDeleteCollectionsIDDeleteDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}

		notFoundResponse = func() middleware.Responder {
			return collections.NewDeleteCollectionsIDDeleteNotFound()
		}

		collectionLocked = func() middleware.Responder {
			return collections.NewDeleteCollectionsIDDeleteLocked()
		}
		forbiddenResponse = collections.NewDeleteCollectionsIDDeleteForbidden
		accessChecker     = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	)
	_ = collectionLocked
	logger.Debug("delete collection")

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin tx", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// get collection
	var collectionToDelete *dto.Collection
	if collectionsList, err := tx.GetCollectionsByIDsForUpdate([]int64{params.ID}); err != nil {
		logger.Error("GetCollectionsByIDs", "err", err)
		return errorResponse(500, err.Error())
	} else if len(collectionsList) != 1 {
		logger.Error("cant find collection")
		return notFoundResponse()
	} else {
		collectionToDelete = collectionsList[0]
	}

	// check rights
	ok, err := accessChecker.HasUserRightsForCollections(
		userContext.User.ID,
		dto.RightEntityLevelWrite,
		[]int64{collectionToDelete.ID},
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant delete collection %d", collectionToDelete.ID)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	if err := tx.DeleteEntityRights(dto.RightEntityTypeCollection, collectionToDelete.ID); err != nil {
		logger.Error("cant delete rights to collection", "err", err)
		return errorResponse(500, err.Error())
	}

	// get trash collection
	rootTrashCollection, err := tx.GetTypedCollection(collectionToDelete.RootID, dto.TrashCollectionTypo)
	if err != nil {
		logger.Error("cant find trash collection", "err", err)
		return errorResponse(500, err.Error())
	}
	if rootTrashCollection == nil {
		err := fmt.Errorf("cant find trash collection %d", collectionToDelete.RootID)
		logger.Error("cant find trash collection", "err", err)
		return errorResponse(500, err.Error())
	}

	if !collectionToDelete.IsTrash() {
		deleteAll := params.Final

		if deleteAll { // delete collection and objects
			if err := tx.DeleteObjectsByCollectionsIDs([]int64{collectionToDelete.ID}); err != nil {
				logger.Error("cant delete object from collections", "id", collectionToDelete.ID)
				return errorResponse(500, err.Error())
			}
		} else { // delete collection, mv objects to trash
			from, to := collectionToDelete.ID, rootTrashCollection.ID
			if err := tx.ChangeObjectsCollection(from, to); err != nil {
				logger.Error("cant mv objects",
					"from", collectionToDelete.ID,
					"to", rootTrashCollection.ID,
				)
			}
		}
		if collectionToDelete.IsRegular() {
			if err := tx.DeleteCollections([]int64{collectionToDelete.ID}); err != nil {
				logger.Error("cant delete collections", "id", collectionToDelete.ID)
				return errorResponse(500, err.Error())
			}
		}
	} else if collectionToDelete.IsTrash() { // collection is trash // explicit
		if err := tx.DeleteObjectsByCollectionsIDs([]int64{collectionToDelete.ID}); err != nil {
			logger.Error("cant delete object from collections", "id", collectionToDelete.ID)
			return errorResponse(500, err.Error())
		}
	}

	var deletedCollectionEvent *dto.Event
	if !collectionToDelete.IsTrash() {
		deletedCollectionEvent = &dto.Event{
			UserID:       userContext.User.ID,
			Type:         dto.EventTypeDeletedEntity,
			CreationTime: time.Now(),
			EventUnion: dto.EventUnion{
				DeletedEntity: &dto.EntityRef{
					Typo: dto.CollectionEntityType,
					ID:   collectionToDelete.ID,
				},
			},
		}
		if _, err := services.EmplaceEvent(tx, deletedCollectionEvent); err != nil {
			logger.Error("delete collection event", "error", err)
			return errorResponse(500, err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "err", err)
		return errorResponse(500, err.Error())
	}

	reindexJob := func() error {
		err := c.Context.SearchClient.ScrollThrought(&services.ScrollSearchQuery{
			RootID: &collectionToDelete.RootID,
			Filters: &dto.ObjectSearchFilters{
				Collections: []int64{collectionToDelete.ID},
			},
		}, c.Context.ReindexObjects)
		if err != nil {
			logger.Error("cant update object", "err", err)
			return err
		}

		if deletedCollectionEvent != nil {
			c.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
				c.Context.EventSender.Send(deletedCollectionEvent)
			}))
		}
		return nil
	}

	c.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		if err := reindexJob(); err != nil {
			logger.Error("reindexJob error", "err", err)
		}
	}))

	return collections.NewDeleteCollectionsIDDeleteNoContent()
}
