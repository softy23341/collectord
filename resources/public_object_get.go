package resource

import (
	"errors"

	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/objects"
	"git.softndit.com/collector/backend/restapi/operations/public_objects"
	"github.com/go-openapi/runtime/middleware"
)

// PublicCollections TBD
type PublicObjects struct {
	Context Context
}

// GetObjectByID TBD
func (po *PublicObjects) GetPublicObjectByID(params public_objects.GetPublicObjectsIDParams) middleware.Responder {
	var (
		logger = po.Context.Log

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

		DBM = po.Context.DBM

		forbiddenResponse = objects.NewGetObjectsIDForbidden
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
		logger.Error("access denied", "err", "collection is not public")
		return forbiddenResponse()
	}

	// remove public data
	if collection.IsAnonymous {
		collection.RootID = 0
	}

	// collections owners
	var user *dto.User
	userRootRefList, err := DBM.GetUserRootRefs(collections.RootsIDs())
	if err != nil {
		logger.Error("GetUserRootRefs", "err", err)
		return errorResponse(500)
	}
	userList, err := DBM.GetUsersByIDs(userRootRefList.GetRootOwnerIDs())
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500)
	}
	if len(userList) == 1 {
		user = userList[0]
		if user.AvatarMediaID != nil {
			additionalMediasIDs = append(additionalMediasIDs, *user.AvatarMediaID)
		}
	}

	objectCnt, err := DBM.GetObjectsCnt([]int64{collection.ID})
	if err != nil {
		logger.Debug("GetObjectsCnt", "error", err)
		return errorResponse(500)
	}

	if mediaID := collection.ImageMediaID; mediaID != nil {
		additionalMediasIDs = append(additionalMediasIDs, *mediaID)
	}

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

	// construct object
	modelObject := models.NewModelObject(object).
		WithMaterialIDs(materials.IDs()).
		WithMediaIDs(objectToMediasIDs[object.ID]).
		WithActorIDs(actorExtractor.getActorsIDsByObjectID(object.ID)).
		WithOriginLocationIDs(originLocationsIDs)

	// sanitize private data
	modelObject.PurchaseDate = nil
	modelObject.PurchasePrice = nil

	// medias
	allMedias := append(medias, additionalMedias...)

	// origin locations
	originLocationList := originLocationExtractor.getOriginLocations()

	payload := &models.AObject{
		Object:          modelObject,
		Actors:          models.NewModelActorList(actorExtractor.getActors()),
		OriginLocations: models.NewModelOriginLocationList(originLocationList),
		Materials:       models.NewModelMaterialList(materials),
		Medias:          models.NewModelMediaList(allMedias),
		DateIntervals:   models.NewModelDateIntervalList(namedDateIntervals),
		Collection:      models.NewModelCollection(collection).WithObjectsCnt(objectCnt),
		User:            models.NewModelUser(user),
	}

	return objects.NewGetObjectsIDOK().WithPayload(payload)
}
