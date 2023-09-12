package resource

import (
	"errors"
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/errs"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/objects"
	"git.softndit.com/collector/backend/services"
	"github.com/go-openapi/runtime/middleware"
)

// UpdateObject TBD
func (o *Object) UpdateObject(params objects.PutObjectsIDParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)

		DBM = o.Context.DBM

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return objects.NewPutObjectsIDDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		successResponse = func() middleware.Responder {
			return objects.NewPutObjectsIDNoContent()
		}
		forbiddenResponse = objects.NewPutObjectsIDForbidden

		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	)

	logger.Debug("update object")

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
		err := fmt.Errorf("cant edit object in %d", params.ID)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	// Get object
	var object *dto.Object
	if objectsList, err := DBM.GetObjectsByIDs([]int64{params.ID}); err != nil {
		logger.Error("get object", "err", err)
		return errorResponse(500, err.Error())
	} else if len(objectsList) != 1 {
		return errorResponse(404, "object not found")
	} else {
		object = objectsList[0]
	}

	objectParams := params.REditObject.Object
	if objectParams == nil {
		return errorResponse(422, "empty body")
	}

	var targetCollection, sourceCollection *dto.Collection
	if collections, err := DBM.GetCollectionsByIDs([]int64{
		object.CollectionID, *objectParams.CollectionID,
	}); err != nil {
		logger.Error("can't get collection GetCollectionsByIDs", "err", err)
		return errorResponse(422, err.Error())
	} else if len(collections) == 0 {
		logger.Error("can't get collection GetCollectionsByIDs", "err", "404")
		return errorResponse(422, "can't find parent collection")
	} else {
		id2collection := collections.IDToCollection()
		targetCollection = id2collection[*objectParams.CollectionID]
		sourceCollection = id2collection[object.CollectionID]

		if targetCollection == nil {
			logger.Error("cant find target collection", "id", *objectParams.CollectionID)
			return errorResponse(422, "cant find target collection")
		}
		if sourceCollection == nil {
			logger.Error("cant find source collection", "id", object.CollectionID)
			return errorResponse(422, "cant find source collection")
		}

		if targetCollection.RootID != sourceCollection.RootID {
			logger.Error("cant tranfer objects between roots collections")
			return forbiddenResponse()
		}
	}
	rootID := targetCollection.RootID

	// check rights
	ok, err = accessChecker.HasUserRightsForCollections(
		userContext.User.ID,
		dto.RightEntityLevelWrite,
		[]int64{targetCollection.ID},
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant edit object in %d", params.ID)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	// name set
	if objectParams.Name != "" {
		object.Name = objectParams.Name
	}

	objectParams.ProductionDate.AddProductionDateToObject(object)
	if _, err := objectParams.PurchaseDate.AddPurchaseDateToObject(object); err != nil {
		return errorResponse(422, err.Error())
	}
	objectParams.PurchasePrice.AddPurchasePriceToObject(object)

	if objectParams.Description != nil {
		object.Description = *objectParams.Description
	}

	if objectParams.Provenance != nil {
		object.Provenance = *objectParams.Provenance
	}

	if rootIDNumber := objectParams.RootIDNumber; rootIDNumber != nil {
		object.RootIDNumber = *rootIDNumber
	}

	if targetCollectionID := objectParams.CollectionID; targetCollectionID != nil {
		object.CollectionID = *targetCollectionID
	}

	// production date
	var productionDateInterval *dto.NamedDateInterval
	if prDate := objectParams.ProductionDate; prDate != nil && prDate.DateIntervalID != nil {
		intervals, err := DBM.GetNamedDateIntervalsByIDs([]int64{*prDate.DateIntervalID})
		if err != nil {
			logger.Error("GetNamedDateIntervalsByIDs", "err", err)
			return errorResponse(500, err.Error())
		}
		if len(intervals) != 1 {
			logger.Error("wrong production date interval", "int", *prDate.DateIntervalID)
			return errorResponse(422, "wrong production date interval")
		}
		productionDateInterval = intervals[0]
	}

	if productionDateInterval != nil {
		if object.IsProductionDateIntervalPresent() {
			logger.Error("inconsistent production date")
			return errorResponse(500, "inconsistent production date")
		}
	}

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// materials handler
	if objectParams.Materials != nil {
		if err := tx.DeleteObjectMaterialsRefs(object.ID); err != nil {
			logger.Error("delete object material refs", "error", err)
			return errorResponse(500, err.Error())
		}

		materialManager := &materialManager{DBM: tx, log: o.Context.Log, rootID: rootID}
		outMaterials, err := materialManager.createMaterialsByNormalNames(objectParams.Materials)
		if err != nil {
			logger.Error("create or find materials", "error", err)
			return errorResponse(500, err.Error())
		}

		if err = tx.CreateObjectMaterialsRefs(object.ID, outMaterials.IDs()); err != nil {
			logger.Error("CreateObjectMaterialsRefs", "error", err)
			return errorResponse(500, err.Error())
		}
	}

	// actors handler
	if objectParams.Actors != nil {
		if err := tx.DeleteObjectActorsRefs(object.ID); err != nil {
			logger.Error("delete object actor refs", "error", err)
			return errorResponse(500, err.Error())
		}

		actorManager := &actorManager{DBM: tx, log: o.Context.Log, rootID: rootID}
		outActors, err := actorManager.createActorsByNormalNames(objectParams.Actors)
		if err != nil {
			logger.Error("create or find actors", "error", err)
			return errorResponse(500, err.Error())
		}

		if err = tx.CreateObjectActorsRefs(object.ID, outActors.IDs()); err != nil {
			logger.Error("CreateObjectActorsRefs", "error", err)
			return errorResponse(500, err.Error())
		}
	}

	// originLocations handler
	if objectParams.OriginLocations != nil {
		if err := tx.DeleteObjectOriginLocationsRefs(object.ID); err != nil {
			logger.Error("delete object originLocation refs", "error", err)
			return errorResponse(500, err.Error())
		}

		originLocationManager := &originLocationManager{DBM: tx, log: o.Context.Log, rootID: rootID}
		outOriginLocations, err := originLocationManager.createOriginLocationsByNormalNames(objectParams.OriginLocations)
		if err != nil {
			logger.Error("create or find originLocations", "error", err)
			return errorResponse(500, err.Error())
		}

		if err = tx.CreateObjectOriginLocationsRefs(object.ID, outOriginLocations.IDs()); err != nil {
			logger.Error("CreateObjectOriginLocationsRefs", "error", err)
			return errorResponse(500, err.Error())
		}
	}

	// badges handler
	if objectParams.Badges != nil {
		if err := tx.DeleteObjectBadgesRefs(object.ID); err != nil {
			logger.Error("delete object badge refs", "error", err)
			return errorResponse(500, err.Error())
		}

		badgeManager := &badgeManager{DBM: tx, log: o.Context.Log, rootID: rootID}
		outBadges, err := badgeManager.createBadgesByNormalNames(models.NewDtoInputBadgeList(objectParams.Badges))
		if err != nil {
			logger.Error("create or find badges", "error", err)
			return errorResponse(500, err.Error())
		}

		if err = tx.CreateObjectBadgesRefs(object.ID, outBadges.IDs()); err != nil {
			logger.Error("CreateObjectBadgesRefs", "error", err)
			return errorResponse(500, err.Error())
		}
	}

	// Medias handler
	if objectParams.MediasIds != nil {
		if err := tx.DeleteObjectMediasRefs(object.ID); err != nil {
			logger.Error("delete object media refs", "error", err)
			return errorResponse(500, err.Error())
		}

		mediaManager := &mediaManager{DBM: tx, log: o.Context.Log}
		mediaList, err := mediaManager.getMediasByIDs(objectParams.MediasIds)
		if err != nil {
			logger.Error("get medias by ids", "error", err)
			return errorResponse(errs.Error2HTTPCode(err), err.Error())
		}

		if len(mediaList) != len(objectParams.MediasIds) {
			err = errors.New("cant find all new media")
			logger.Error("get medias by ids", "error", err)
			return errorResponse(errs.Error2HTTPCode(err), err.Error())
		}

		// check access rights
		ok := accessChecker.IsUserOwnerOfMedias(userContext.User.ID, mediaList)
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return errorResponse(500, err.Error())
		} else if !ok {
			err := errors.New("user is not owner of media")
			logger.Error("user cant allow be here", "err", err)
			return forbiddenResponse()
		}

		if err = tx.CreateObjectMediasRefs(object.ID, objectParams.MediasIds); err != nil {
			logger.Error("CreateObjectMediasRefs", "error", err)
			return errorResponse(500, err.Error())
		}
	}

	// currencies
	var currenciesIDs []int64
	if purchasePrice := objectParams.PurchasePrice; purchasePrice != nil && purchasePrice.CurrencyID != nil {
		currenciesIDs = append(currenciesIDs, *purchasePrice.CurrencyID)
	}
	if err := o.validateCurrencies(currenciesIDs); err != nil {
		logger.Error("GetCurrenciesByIDs error", "error", err)
		return errorResponse(500, err.Error())
	}

	// object status
	if statusRef := objectParams.ObjectStatus; statusRef != nil {
		if list, err := tx.GetObjectStatusByIDs([]int64{*statusRef.StatusID}); err != nil {
			logger.Error("cant get object status")
			return errorResponse(500, err.Error())
		} else if len(list) != 1 {
			logger.Error("cant get object status")
			return errorResponse(500, err.Error())
		}

		// compare with current status
		refs, err := tx.GetCurrentObjectsStatusesRefs([]int64{object.ID})
		if err != nil {
			logger.Error("GetCurrentObjectsStatusesRefs", "err", err)
			return errorResponse(500, err.Error())
		}

		currentObjectStatus := refs.ObjectToOneStatusMap()[object.ID]
		needUpdate := (currentObjectStatus == nil) ||
			!(currentObjectStatus.Description == *statusRef.Description &&
				currentObjectStatus.ObjectStatusID == *statusRef.StatusID)

		if needUpdate {
			ref := &dto.ObjectStatusRef{
				ObjectID:       object.ID,
				ObjectStatusID: *statusRef.StatusID,
				Description:    *statusRef.Description,
			}
			if err := tx.CreateObjectStatusRef(ref); err != nil {
				logger.Error("cant create object ref", "err", err)
				return errorResponse(500, err.Error())
			}
		}
	}

	// valuations
	var valuationsCurrenciesIDs []int64
	if objectParams.Valuations != nil && objectParams.Valuations.List != nil {
		for i := range objectParams.Valuations.List {
			valuationsCurrenciesIDs = append(valuationsCurrenciesIDs, *objectParams.Valuations.List[i].CurrencyID)
		}
	}
	err = o.validateCurrencies(valuationsCurrenciesIDs)
	if err != nil {
		logger.Error("Valuations GetCurrenciesByIDs error", "error", err)
		return errorResponse(500, err.Error())
	}

	valuationsManager := &valuationsManager{DBM: tx, log: o.Context.Log}
	var valuationList dto.ValuationList
	if objectParams.Valuations != nil {
		valuationList, err = models.NewDtoValuationList(objectParams.Valuations)
		if err != nil {
			logger.Error("create valuations parse date", "error", err)
			return errorResponse(422, err.Error())
		}
	} else {
		valuationList = nil
	}
	_, err = valuationsManager.updateValuations(object.ID, valuationList)
	if err != nil {
		logger.Error("create valuations", "error", err)
		return errorResponse(500, err.Error())
	}

	// update object
	object.UpdateTime = time.Now()
	if err := tx.UpdateObject(object); err != nil {
		logger.Error("update object", "error", err)
		return errorResponse(500, err.Error())
	}

	var moveObjectEvent *dto.Event
	if sourceCollection.ID != targetCollection.ID {
		moveObjectEvent = &dto.Event{
			UserID:       userContext.User.ID,
			Type:         dto.EventTypeMovedObject,
			CreationTime: time.Now(),
			EventUnion: dto.EventUnion{
				MovedObject: &dto.EventMovedObject{
					ObjectID:             object.ID,
					OriginalCollectionID: sourceCollection.ID,
					NewCollectionID:      targetCollection.ID,
				},
			},
		}
		if _, err := services.EmplaceEvent(tx, moveObjectEvent); err != nil {
			logger.Error("move object event", "error", err)
			return errorResponse(500, err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	o.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		err := o.Context.SearchClient.ScrollThrought(&services.ScrollSearchQuery{
			RootID:    &rootID,
			ObjectIDs: []int64{object.ID},
		}, o.Context.ReindexObjects)
		if err != nil {
			logger.Error("cant update object", "err", err)
		}

		if moveObjectEvent != nil {
			o.Context.EventSender.Send(moveObjectEvent)
		}
	}))

	return successResponse()
}
