package resource

import (
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/objects"
	"git.softndit.com/collector/backend/services"
	"github.com/go-openapi/runtime/middleware"
)

// DeleteObject moves object to trash or deletes it if already in
func (o *Object) DeleteObject(params objects.DeleteObjectsIDDeleteParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return objects.NewDeleteObjectsIDDeleteDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}

		objectLockedResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return objects.NewDeleteObjectsIDDeleteLocked().WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		_ = objectLockedResponse // reserved

		notFoundResponder = func() middleware.Responder {
			return objects.NewDeleteObjectsIDDeleteNotFound()
		}
		forbiddenResponse = objects.NewDeleteObjectsIDDeleteForbidden

		DBM           = o.Context.DBM
		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	)

	logger.Debug("delete object")

	// check rights
	ok, err := accessChecker.HasUserRightsForObjects(
		userContext.User.ID,
		dto.RightEntityLevelWrite,
		[]int64{params.ID},
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant delete object %d", params.ID)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	// get object
	var objectToDelete *dto.Object
	if objectList, err := DBM.GetObjectsByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant find object", "err", err)
		return errorResponse(500, err.Error())
	} else if len(objectList) != 1 {
		logger.Error("cant find object")
		return notFoundResponder()
	} else {
		objectToDelete = objectList[0]
	}
	forIndex := objectToDelete.SearchObject()

	reindexExternalVersion, err := DBM.GetCurrentTransaction()
	if err != nil {
		logger.Error("cant find transaction", "err", err)
		return errorResponse(500, err.Error())
	}

	forIndex.SetExternalVersion(reindexExternalVersion)

	// get collection
	var objectCollection *dto.Collection
	if colllectionList, err := DBM.GetCollectionsByIDs([]int64{objectToDelete.CollectionID}); err != nil {
		logger.Error("cant find collection", "err", err)
		return errorResponse(500, err.Error())
	} else if len(colllectionList) != 1 {
		logger.Error("cant find collection", "err", err)
		return errorResponse(500, err.Error())
	} else {
		objectCollection = colllectionList[0]
	}

	// get trash collection
	rootTrashCollection, err := DBM.GetTypedCollection(objectCollection.RootID, dto.TrashCollectionTypo)
	if err != nil {
		logger.Error("cant find trash collection", "err", err)
		return errorResponse(500, err.Error())
	}
	if rootTrashCollection == nil {
		err := fmt.Errorf("cant find trash collection %d", objectCollection.RootID)
		logger.Error("cant find trash collection", "err", err)
		return errorResponse(500, err.Error())
	}

	var reindexJob func()

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// check delete
	finalDelete := (rootTrashCollection.ID == objectCollection.ID) || params.Final
	logger.Debug("delete object", "final", finalDelete)
	if finalDelete {
		if err := tx.DeleteObjectsByIDs([]int64{objectToDelete.ID}); err != nil {
			logger.Error("cant delete object", "err", err)
			return errorResponse(500, err.Error())
		}

		reindexJob = func() {
			err := o.Context.SearchClient.BulkObjectDelete([]services.ObjectDocForIndex{forIndex})
			if err != nil {
				logger.Error("delete object search", "err", err)
			}
		}
	} else {
		objectToDelete.CollectionID = rootTrashCollection.ID
		if err := tx.UpdateObject(objectToDelete); err != nil {
			logger.Error("cant move object to trash", "err", err)
			return errorResponse(500, err.Error())
		}

		reindexJob = func() {
			err := o.Context.SearchClient.IndexObject(forIndex)
			if err != nil {
				logger.Error("update object search", "err", err)
			}
		}
	}

	deletedObjectEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeDeletedObject,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			DeletedObject: &dto.EventDeletedObject{
				ObjectID:             objectToDelete.ID,
				OriginalCollectionID: objectToDelete.CollectionID,
			},
		},
	}

	if _, err := services.EmplaceEvent(tx, deletedObjectEvent); err != nil {
		logger.Error("delete object event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	o.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		reindexJob()
		o.Context.EventSender.Send(deletedObjectEvent)
	}))

	return objects.NewDeleteObjectsIDDeleteNoContent()
}
