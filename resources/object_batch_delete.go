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

// ButchDeleteObject moves object to trash or deletes it if already in
func (o *Object) ButchDeleteObject(params objects.PostObjectsDeleteParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return objects.NewPostObjectsDeleteDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}

		objectLockedResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return objects.NewPostObjectsDeleteLocked().WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}

		notFoundResponder = func() middleware.Responder {
			return objects.NewPostObjectsDeleteNotFound()
		}
		forbiddenResponse = objects.NewPostObjectsDeleteForbidden

		DBM           = o.Context.DBM
		objectIDs     = params.RDeleteObjects.ObjectsIds
		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	)

	logger.Debug("batch delete object")

	// check rights
	ok, err := accessChecker.HasUserRightsForObjects(
		userContext.User.ID,
		dto.RightEntityLevelWrite,
		objectIDs,
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant delete objects %v", objectIDs)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	// get object
	objectList, err := DBM.GetObjectsByIDs(objectIDs)
	if err != nil {
		logger.Error("cant find object", "err", err)
		return errorResponse(500, err.Error())
	} else if len(objectList) != len(objectIDs) || len(objectList) == 0 {
		logger.Error("cant find object")
		return notFoundResponder()
	}

	collectionID := objectList[0].CollectionID
	for _, object := range objectList {
		if collectionID != object.CollectionID {
			logger.Error("objects don't belong same collection")
			objectLockedResponse(500, "objects don't belong same collection")
		}
	}
	forIndex := objectList.SearchObjects()

	reindexExternalVersion, err := DBM.GetCurrentTransaction()
	if err != nil {
		logger.Error("cant find transaction", "err", err)
		return errorResponse(500, err.Error())
	}

	docsForIndex := make(services.ObjectDocsForIndex, len(forIndex))
	for i, objectForIndex := range forIndex {
		objectForIndex.SetExternalVersion(reindexExternalVersion)
		docsForIndex[i] = objectForIndex
	}

	// get collection
	var objectCollection *dto.Collection
	if colllectionList, err := DBM.GetCollectionsByIDs([]int64{collectionID}); err != nil {
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
	finalDelete := (rootTrashCollection.ID == objectCollection.ID) || params.RDeleteObjects.Final
	logger.Debug("delete object", "final", finalDelete)
	if finalDelete {
		if err := tx.DeleteObjectsByIDs(objectIDs); err != nil {
			logger.Error("cant delete object", "err", err)
			return errorResponse(500, err.Error())
		}

		reindexJob = func() {
			err := o.Context.SearchClient.BulkObjectDelete(docsForIndex)
			if err != nil {
				logger.Error("delete object search", "err", err)
			}
		}
	} else {
		for _, o := range objectList {
			o.CollectionID = rootTrashCollection.ID
			if err := tx.UpdateObject(o); err != nil {
				logger.Error("cant move object to trash", "err", err)
				return errorResponse(500, err.Error())
			}
		}

		reindexJob = func() {
			err := o.Context.SearchClient.BulkObjectIndex(docsForIndex)
			if err != nil {
				logger.Error("update object search", "err", err)
			}
		}
	}

	events := make(dto.EventList, len(objectList))
	for i, object := range objectList {
		event := &dto.Event{
			UserID:       userContext.User.ID,
			Type:         dto.EventTypeDeletedObject,
			CreationTime: time.Now(),
			EventUnion: dto.EventUnion{
				DeletedObject: &dto.EventDeletedObject{
					ObjectID:             object.ID,
					OriginalCollectionID: collectionID,
				},
			},
		}

		if _, err := services.EmplaceEvent(tx, event); err != nil {
			logger.Error("delete object event", "error", err)
			return errorResponse(500, err.Error())
		}
		events[i] = event
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	o.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		reindexJob()
		o.Context.EventSender.Send(events...)
	}))

	return objects.NewDeleteObjectsIDDeleteNoContent()
}
