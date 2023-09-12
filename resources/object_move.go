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

// MoveObjects TBD
func (o *Object) MoveObjects(params objects.PostObjectsMoveParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)

		DBM = o.Context.DBM

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return objects.NewPostObjectsMoveDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		forbiddenResponse = objects.NewPostObjectsMoveForbidden

		moveParams    = params.RMoveObjects
		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	)

	logger.Debug("move object")

	var objectList dto.ObjectList
	objectList, err := DBM.GetObjectsByIDs(moveParams.ObjectsIds)
	if err != nil {
		logger.Error("cant get objects", "err", err)
		return errorResponse(500, err.Error())
	} else if len(objectList) != len(moveParams.ObjectsIds) {
		err := fmt.Errorf("cant find all objects")
		logger.Error("cant get objects", "err", err)
		return errorResponse(500, err.Error())
	}

	// check rights
	ok, err := accessChecker.HasUserRightsForObjects(
		userContext.User.ID,
		dto.RightEntityLevelWrite,
		moveParams.ObjectsIds,
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant edit object in %v", moveParams.ObjectsIds)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	// check rights
	ok, err = accessChecker.HasUserRightsForCollections(
		userContext.User.ID,
		dto.RightEntityLevelWrite,
		[]int64{*moveParams.TargetCollectionID},
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant edit object in %d", *moveParams.TargetCollectionID)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	err = tx.ChangeObjectsCollectionByIDs(moveParams.ObjectsIds, *moveParams.TargetCollectionID)
	if err != nil {
		logger.Error("cant ChangeObjectsCollectionByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	var events dto.EventList
	for _, object := range objectList {
		deletedObjectEvent := &dto.Event{
			UserID:       userContext.User.ID,
			Type:         dto.EventTypeMovedObject,
			CreationTime: time.Now(),
			EventUnion: dto.EventUnion{
				MovedObject: &dto.EventMovedObject{
					ObjectID:             object.ID,
					OriginalCollectionID: object.CollectionID,
					NewCollectionID:      *moveParams.TargetCollectionID,
				},
			},
		}

		if _, err := services.EmplaceEvent(tx, deletedObjectEvent); err != nil {
			logger.Error("new object event", "error", err)
			return errorResponse(500, err.Error())
		}
		events = append(events, deletedObjectEvent)
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	o.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		err := o.Context.SearchClient.ScrollThrought(&services.ScrollSearchQuery{
			ObjectIDs: moveParams.ObjectsIds,
		}, o.Context.ReindexObjects)

		if err != nil {
			logger.Error("cant update object", "err", err)
		}

		o.Context.EventSender.Send(events...)
	}))

	return objects.NewPostObjectsMoveNoContent()
}
