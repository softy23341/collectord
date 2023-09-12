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

// CreateObject TBD
func (o *Object) CreateObject(params objects.PostObjectsParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)

		userID = userContext.User.ID

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return objects.NewPostObjectsDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		successResponse = func(objectID int64) middleware.Responder {
			return objects.NewPostObjectsOK().WithPayload(&models.ANewObject{
				ID: &objectID,
			})
		}
		forbiddenResponse = objects.NewPostObjectsForbidden

		DBM           = o.Context.DBM
		clientUniqID  = params.RNewObject.ClientUniqID
		inputObject   = params.RNewObject.Object
		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	)

	// check rights
	ok, err := accessChecker.HasUserRightsForCollections(
		userContext.User.ID,
		dto.RightEntityLevelWrite,
		[]int64{*inputObject.CollectionID},
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant create object in %d", inputObject.CollectionID)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	objectID, err := DBM.GetObjectIDByUserUniqID(userID, *clientUniqID)
	if err != nil {
		logger.Error("get object by user id", "error", err)
		return errorResponse(500, err.Error())
	}
	if objectID != nil {
		logger.Debug("Object already find", "objectID", objectID)
		return successResponse(*objectID)
	}

	// collections
	var collection *dto.Collection
	if collections, err := DBM.GetCollectionsByIDs([]int64{*inputObject.CollectionID}); err != nil {
		logger.Error("get collection by ids", "error", err)
		return errorResponse(500, err.Error())
	} else if len(collections) != 1 {
		err := errors.New("Can't find collections by id")
		logger.Error("get collection by ids", "error", err)
		return errorResponse(500, err.Error())
	} else {
		collection = collections[0]
	}

	rootID := collection.RootID

	// medias
	var medias dto.MediaList
	if len(inputObject.MediasIds) > 0 {
		mediaManager := &mediaManager{DBM: DBM, log: o.Context.Log}
		medias, err = mediaManager.getMediasByIDs(inputObject.MediasIds)
		if err != nil {
			logger.Error("get medias by ids", "error", err)
			return errorResponse(errs.Error2HTTPCode(err), err.Error())
		}
		if len(medias) != len(inputObject.MediasIds) {
			err = errors.New("cant find all medias")
			logger.Error("get medias by ids", "error", err)
			return errorResponse(errs.Error2HTTPCode(err), err.Error())
		}

		// check access rights
		ok := accessChecker.IsUserOwnerOfMedias(userContext.User.ID, medias)
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return errorResponse(500, err.Error())
		} else if !ok {
			err := errors.New("user is not owner of media")
			logger.Error("user cant allow be here", "err", err)
			return forbiddenResponse()
		}
	}

	// production date
	var productionDateInterval *dto.NamedDateInterval
	if prDate := inputObject.ProductionDate; prDate != nil && prDate.DateIntervalID != nil {
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

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin transaction", "error", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// materials
	materialManager := &materialManager{DBM: tx, log: o.Context.Log, rootID: rootID}
	outMaterials, err := materialManager.createMaterialsByNormalNames(inputObject.Materials)
	if err != nil {
		logger.Error("create or find materials", "error", err)
		return errorResponse(500, err.Error())
	}

	// actors
	actorManager := &actorManager{DBM: tx, log: o.Context.Log, rootID: rootID}
	outActors, err := actorManager.createActorsByNormalNames(inputObject.Actors)
	if err != nil {
		logger.Error("create or find materials", "error", err)
		return errorResponse(500, err.Error())
	}

	// origin locations
	originLocationManager := &originLocationManager{DBM: tx, log: o.Context.Log, rootID: rootID}
	outOriginLocations, err := originLocationManager.createOriginLocationsByNormalNames(inputObject.OriginLocations)
	if err != nil {
		logger.Error("create or find materials", "error", err)
		return errorResponse(500, err.Error())
	}

	// badges
	badgeManager := &badgeManager{DBM: tx, log: o.Context.Log, rootID: rootID}
	outBadges, err := badgeManager.createBadgesByNormalNames(models.NewDtoInputBadgeList(inputObject.Badges))
	if err != nil {
		logger.Error("create or find badges", "error", err)
		return errorResponse(500, err.Error())
	}

	// currencies
	var currenciesIDs []int64
	if purchasePrice := inputObject.PurchasePrice; purchasePrice != nil && purchasePrice.CurrencyID != nil {
		currenciesIDs = append(currenciesIDs, *purchasePrice.CurrencyID)
	}
	if err := o.validateCurrencies(currenciesIDs); err != nil {
		logger.Error("GetCurrenciesByIDs error", "error", err)
		return errorResponse(500, err.Error())
	}

	// Object creation
	object := dto.NewObject()
	object.Name = inputObject.Name
	object.UserID = userID
	object.UserUniqID = *clientUniqID
	object.CollectionID = collection.ID

	inputObject.PurchasePrice.AddPurchasePriceToObject(object)
	inputObject.ProductionDate.AddProductionDateToObject(object)

	if _, err := inputObject.PurchaseDate.AddPurchaseDateToObject(object); err != nil {
		return errorResponse(422, err.Error())
	}

	if inputObject.Description != nil {
		object.Description = *inputObject.Description
	}

	if inputObject.Provenance != nil {
		object.Provenance = *inputObject.Provenance
	}

	if rootIDNumber := inputObject.RootIDNumber; rootIDNumber != nil {
		object.RootIDNumber = *rootIDNumber
	}

	err = tx.CreateObject(object)
	if err != nil {
		logger.Error("create object", "error", err)
		return errorResponse(500, err.Error())
	}
	forIndex := object.SearchObject().WithRootID(rootID)

	// object status
	if statusRef := inputObject.ObjectStatus; statusRef != nil {
		if list, err := tx.GetObjectStatusByIDs([]int64{*statusRef.StatusID}); err != nil {
			logger.Error("cant get object status")
			return errorResponse(500, err.Error())
		} else if len(list) != 1 {
			logger.Error("cant get object status")
			return errorResponse(500, err.Error())
		}

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

	// valuations
	var valuationsCurrenciesIDs []int64
	if inputObject.Valuations != nil && inputObject.Valuations.List != nil {
		for i := range inputObject.Valuations.List {
			valuationsCurrenciesIDs = append(valuationsCurrenciesIDs, *inputObject.Valuations.List[i].CurrencyID)
		}
	}
	err = o.validateCurrencies(valuationsCurrenciesIDs)
	if err != nil {
		logger.Error("Valuations GetCurrenciesByIDs error", "error", err)
		return errorResponse(500, err.Error())
	}

	if inputObject.Valuations != nil && inputObject.Valuations.List != nil {
		valuationsManager := &valuationsManager{DBM: tx, log: o.Context.Log}
		valuationList, err := models.NewDtoValuationList(inputObject.Valuations)
		if err != nil {
			logger.Error("create valuations parse date", "error", err)
			return errorResponse(422, err.Error())
		}
		_, err = valuationsManager.createValuations(object.ID, valuationList)
		if err != nil {
			logger.Error("create valuations", "error", err)
			return errorResponse(500, err.Error())
		}
	}

	if productionDateInterval != nil {
		if object.IsProductionDateIntervalPresent() {
			logger.Error("inconsistent production date")
			return errorResponse(500, "inconsistent production date")
		}

		forIndex.
			WithProductionIntervalFrom(&productionDateInterval.ProductionDateIntervalFrom).
			WithProductionIntervalTo(&productionDateInterval.ProductionDateIntervalTo)
	}

	if len(medias) > 0 {
		if err = tx.CreateObjectMediasRefs(object.ID, inputObject.MediasIds); err != nil {
			logger.Error("CreateObjectMediasRefs", "error", err)
			return errorResponse(500, err.Error())
		}
	}

	if len(outMaterials) > 0 {
		forIndex.WithMaterialsIDs(outMaterials.IDs())
		if err = tx.CreateObjectMaterialsRefs(object.ID, outMaterials.IDs()); err != nil {
			logger.Error("CreateObjectMaterialsRefs", "error", err)
			return errorResponse(500, err.Error())
		}
	}

	if len(outActors) > 0 {
		forIndex.
			WithActorsIDs(outActors.IDs()).
			WithActorsNames(outActors.NormalNames())

		if err = tx.CreateObjectActorsRefs(object.ID, outActors.IDs()); err != nil {
			logger.Error("CreateObjectActorsRefs", "error", err,
				"objectID", object.ID, "actors", outActors.IDs())
			return errorResponse(500, err.Error())
		}
	}

	if len(outOriginLocations) > 0 {
		forIndex.WithOriginLocationsIDs(outOriginLocations.IDs())
		if err = tx.CreateObjectOriginLocationsRefs(object.ID, outOriginLocations.IDs()); err != nil {
			logger.Error("CreateObjectOriginLocationsRefs", "error", err)
			return errorResponse(500, err.Error())
		}
	}

	if len(outBadges) > 0 {
		forIndex.WithBadgesIDs(outBadges.IDs())
		if err = tx.CreateObjectBadgesRefs(object.ID, outBadges.IDs()); err != nil {
			logger.Error("CreateObjectBadgesRefs", "error", err)
			return errorResponse(500, err.Error())
		}
	}

	version, err := tx.GetCurrentTransaction()
	if err != nil {
		logger.Error("GetCurrentTransaction", "error", err)
		return errorResponse(500, err.Error())
	}

	forIndex.SetExternalVersion(version)

	newObjectEvent := &dto.Event{
		UserID:       userContext.User.ID,
		Type:         dto.EventTypeNewObject,
		CreationTime: time.Now(),
		EventUnion: dto.EventUnion{
			NewObject: &dto.EventNewObject{
				ObjectID:     object.ID,
				CollectionID: object.CollectionID,
			},
		},
	}
	if _, err := services.EmplaceEvent(tx, newObjectEvent); err != nil {
		logger.Error("new object event", "error", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	o.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		err := o.Context.SearchClient.IndexObject(forIndex)
		if err != nil {
			logger.Error("insert into search", "err", err)
		}

		o.Context.EventSender.Send(newObjectEvent)
	}))

	return successResponse(object.ID)
}
