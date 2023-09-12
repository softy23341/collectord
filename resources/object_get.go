package resource

import (
	"errors"
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/objects"
	"github.com/go-openapi/runtime/middleware"
)

// GetObjectByID TBD
func (o *Object) GetObjectByID(params objects.GetObjectsIDParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)

		errorResponse = func(code int) middleware.Responder {
			c := int32(code)
			return objects.NewGetObjectsIDDefault(code).WithPayload(
				&models.Error{Code: &c},
			)
		}

		notFoundResponse = func() middleware.Responder {
			c := int32(404)
			return objects.NewGetObjectsIDNotFound().WithPayload(
				&models.Error{Code: &c},
			)
		}

		DBM = o.Context.DBM

		forbiddenResponse = objects.NewGetObjectsIDForbidden
		accessChecker     = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	)

	var additionalMediasIDs []int64

	// objects medias
	mhContext := &mhContext{
		DBM: DBM,
		log: logger,
	}

	// Get object
	objectsList, err := DBM.GetObjectsByIDs([]int64{params.ID})
	if err != nil {
		logger.Error("get object", "err", err)
		return errorResponse(500)
	}
	if len(objectsList) != 1 {
		return notFoundResponse()
	}
	object := objectsList[0]

	// Get group
	collections, err := DBM.GetCollectionsByIDs([]int64{object.CollectionID})
	if err != nil {
		logger.Error("get collections", "err", err)
		return errorResponse(500)
	}
	if len(collections) != 1 {
		err = errors.New("collection not found")
		logger.Debug("GetCollectionsByIDs", "error", err)
		return errorResponse(404)
	}
	collection := collections[0]

	// check rights
	if !collection.Public {
		ok, err := accessChecker.HasUserRightsForObjects(
			userContext.User.ID,
			dto.RightEntityLevelRead,
			[]int64{object.ID},
		)
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return errorResponse(500)
		}
		if !ok {
			if params.MessageID != nil {
				allow, err := DBM.CheckUserTmpAccessToObject(userContext.User.ID, *params.MessageID, object.ID)
				if err != nil {
					logger.Error("CheckUserTmpAccessToObject", "err", err)
					return errorResponse(500)
				} else if !allow {
					logger.Error("tmp access denied")
					return forbiddenResponse()
				}
			} else {
				err := fmt.Errorf("cant edit object in %d", params.ID)
				logger.Error("access denied", "err", err.Error())
				return forbiddenResponse()
			}
		}
	}

	// get rights
	userCollectionRightList, err := DBM.GetUserRightsForCollections(userContext.User.ID, []int64{collection.ID})
	if err != nil {
		logger.Error("GetUserRightsForCollections", "err", err)
		return errorResponse(500)
	}
	collections2level := userCollectionRightList.GetEntityIDToLevel(dto.RightEntityTypeCollection)

	objectCnt, err := DBM.GetObjectsCnt([]int64{collection.ID})
	if err != nil {
		logger.Debug("GetObjectsCnt", "error", err)
		return errorResponse(500)
	}

	if mediaID := collection.ImageMediaID; mediaID != nil {
		additionalMediasIDs = append(additionalMediasIDs, *mediaID)
	}

	// statuses
	objectStatuses, err := DBM.GetCurrentObjectsStatusesRefs([]int64{object.ID})
	if err != nil {
		logger.Error("GetCurrentObjectsStatusesRefs", "err", err)
		return errorResponse(500)
	}
	objectStatusMap := objectStatuses.ObjectToOneStatusMap()

	// Get materials
	materials, err := DBM.GetMaterialsByObjectID(object.ID)
	if err != nil {
		logger.Error("get materials", "err", err)
		return errorResponse(500)
	}

	// Get medias
	mediaManager := &mediaManager{DBM: DBM, log: logger}
	medias, objectToMediasIDs, err := mediaManager.MediasByObjectEntities([]int64{object.GetID()})
	if err != nil {
		logger.Error("get medias", "err", err)
		return errorResponse(500)
	}

	// actors
	actorExtractor, err := newObjectActorsExtractor(mhContext, []int64{object.ID})
	if err != nil {
		logger.Error("GetObjectsActorRefs", "err", err)
		return errorResponse(500)
	}

	// originLocations
	originLocationExtractor, err := newObjectOriginLocationsExtractor(mhContext, []int64{object.ID})
	if err != nil {
		logger.Error("GetObjectsOriginLocationRefs", "err", err)
		return errorResponse(500)
	}

	// badges
	badgeExtractor, err := newObjectBadgesExtractor(mhContext, []int64{object.ID})
	if err != nil {
		logger.Error("GetObjectsBadgeRefs", "err", err)
		return errorResponse(500)
	}

	// date intervals
	var namedDateIntervals dto.NamedDateIntervalList
	if object.ProductionDateIntervalID != nil {
		ids := []int64{*object.ProductionDateIntervalID}
		namedDateIntervals, err = DBM.GetNamedDateIntervalsByIDs(ids)
		if err != nil {
			logger.Error("GetNamedDateIntervalsByIDs", "err", err)
			return errorResponse(500)
		}
	}

	// fetch additional medias
	var additionalMedias dto.MediaList
	if len(additionalMediasIDs) > 0 {
		additionalMedias, err = DBM.GetMediasByIDs(additionalMediasIDs)
		if err != nil {
			logger.Error("GetMediasByIDs", "err", err)
			return errorResponse(500)
		}
	}

	originLocationsIDs := originLocationExtractor.getOriginLocationsIDsByObjectID(object.ID)

	// valuations
	valuationsExtractor, err := newObjectValuationsExtractor(mhContext, []int64{object.ID})
	if err != nil {
		logger.Error("GetValuations", "err", err)
		return errorResponse(500)
	}

	// construct object
	modelObject := models.NewModelObject(object).
		WithMaterialIDs(materials.IDs()).
		WithMediaIDs(objectToMediasIDs[object.ID]).
		WithActorIDs(actorExtractor.getActorsIDsByObjectID(object.ID)).
		WithOriginLocationIDs(originLocationsIDs).
		WithBadgeIDs(badgeExtractor.getBadgesIDsByObjectID(object.ID)).
		WithStatus(objectStatusMap[object.ID]).
		WithValuations(models.NewModelValuationList(valuationsExtractor.GetObjectValuations())).
		WithAccessLevel(collections2level[collection.ID])

	// medias
	allMedias := append(medias, additionalMedias...)

	originLocationList := originLocationExtractor.getOriginLocations()

	payload := &models.AObject{
		Object:          modelObject,
		Actors:          models.NewModelActorList(actorExtractor.getActors()),
		OriginLocations: models.NewModelOriginLocationList(originLocationList),
		Materials:       models.NewModelMaterialList(materials),
		Medias:          models.NewModelMediaList(allMedias),
		Badges:          models.NewModelBadgeList(badgeExtractor.getBadges()),
		DateIntervals:   models.NewModelDateIntervalList(namedDateIntervals),
		Collection:      models.NewModelCollection(collection).WithObjectsCnt(objectCnt),
	}

	return objects.NewGetObjectsIDOK().WithPayload(payload)
}
