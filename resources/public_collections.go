package resource

import (
	"errors"
	"fmt"

	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/public_collections"
	"github.com/go-openapi/runtime/middleware"
)

// PublicCollections TBD
type PublicCollections struct {
	Context Context
}

// GetList TBD
func (d *PublicCollections) GetList(params public_collections.PostPublicCollectionsParams) middleware.Responder {
	logger := d.Context.Log
	logger.Debug("get public collections")

	var (
		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return public_collections.NewPostPublicCollectionsDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		DBM = d.Context.DBM
	)

	// collections
	query := params.RPublicCollections.Query
	userID := params.RPublicCollections.UserID
	paginator := dto.NewDefaultPagePaginator()
	if paramsPaginator := params.RPublicCollections.Paginator; paramsPaginator != nil {
		paginator.Cnt = *paramsPaginator.Cnt
		paginator.Page = *paramsPaginator.Page
	}

	// check root context
	var currentUserRootId *int64
	if userID != nil {
		if rootList, err := DBM.GetMainUserRoot(*userID); err != nil {
			return errorResponse(500, err.Error())
		} else if len(rootList) != 1 {
			return errorResponse(404, "cant find user root")
		} else {
			currentUserRootId = &rootList[0].ID
		}
	}

	totalCollectionsCnt, collections, err := DBM.GetPublicCollections(query, currentUserRootId, paginator)
	if err != nil {
		logger.Error("GetCollectionsByRootID", "err", err)
		return errorResponse(500, err.Error())
	}

	collectionsIDs := collections.GetIDs()

	// collections cnt
	collections2cnt, err := DBM.GetObjectsCntByCollections(collectionsIDs)
	if err != nil {
		logger.Error("GetObjectsCntByCollections", "err", err)
		return errorResponse(500, err.Error())
	}

	modelCollectionsList := models.NewModelCollectionList(collections)
	for _, modelCollection := range modelCollectionsList.List {
		collectionID := *modelCollection.ID
		modelCollection.WithObjectsCnt(collections2cnt[collectionID])
	}

	// remove user for anonymous collection
	for _, c := range collections {
		if c.IsAnonymous {
			c.RootID = 0
		}
	}

	userRootRefList, err := DBM.GetUserRootRefs(collections.RootsIDs())
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// fetch all users
	var usersIDs []int64
	usersIDs = append(usersIDs, userRootRefList.GetRootOwnerIDs()...)
	if userID != nil {
		usersIDs = append(usersIDs, *userID)
	}

	// collections owners
	userList, err := DBM.GetUsersByIDs(usersIDs)
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// fetch all medias
	var allMediasIDs []int64
	allMediasIDs = append(allMediasIDs, collections.GetMediasIDs()...)
	allMediasIDs = append(allMediasIDs, userList.GetAvatarsMediaIDs()...)

	// fetch all medias
	medias, err := DBM.GetMediasByIDs(allMediasIDs)
	if err != nil {
		logger.Error("GetMediasByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	return public_collections.NewPostPublicCollectionsOK().WithPayload(&models.APublicCollections{
		Collections:           modelCollectionsList,
		Medias:                models.NewModelMediaList(medias),
		Users:                 models.NewModelUserList(userList),
		UsersCollections:      models.NewModelUsersCollectionsList(collections, userRootRefList),
		CollectionsTotalCount: totalCollectionsCnt,
	})
}

// GetCollection TBD
func (d *PublicCollections) GetCollection(params public_collections.GetPublicCollectionsIDParams) middleware.Responder {
	var (
		logger = d.Context.Log
		DBM    = d.Context.DBM

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return public_collections.NewGetPublicCollectionsIDDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}

		forbiddenResponse = public_collections.NewGetPublicCollectionsIDForbidden

		collectionsIDs = []int64{params.ID}
	)

	logger.Debug("GetCollection")

	// collections
	var collection *dto.Collection
	if collectionsList, err := DBM.GetCollectionsByIDs(collectionsIDs); err != nil {
		logger.Error("get collections by ids", "err", err)
		return errorResponse(500, err.Error())
	} else if len(collectionsList) == 0 {
		err := fmt.Errorf("cant find collection: %d", params.ID)
		logger.Error("err", "err", err)
		return public_collections.NewGetPublicCollectionsIDNotFound()
	} else {
		collection = collectionsList[0]
	}

	if !collection.Public {
		forbiddenResponse()
	}

	// remove user for anonymous collection
	if collection.IsAnonymous {
		collection.RootID = 0
	}

	// collections cnt
	collections2cnt, err := DBM.GetObjectsCntByCollections(collectionsIDs)
	if err != nil {
		logger.Error("GetObjectsCntByCollections", "err", err)
		return errorResponse(500, err.Error())
	}

	var allMediasIDs []int64
	if img := collection.ImageMediaID; img != nil {
		allMediasIDs = append(allMediasIDs, *img)
	}

	// users
	userRootRefList, err := DBM.GetUserRootRefs([]int64{collection.RootID})
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// collections owners
	userList, err := DBM.GetUsersByIDs(userRootRefList.GetRootOwnerIDs())
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500, err.Error())
	}
	allMediasIDs = append(allMediasIDs, userList.GetAvatarsMediaIDs()...)

	// medias
	medias, err := DBM.GetMediasByIDs(allMediasIDs)
	if err != nil {
		logger.Error("GetMediasByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// collections
	modelCollection := models.NewModelCollection(collection)
	modelCollection.
		WithObjectsCnt(collections2cnt[params.ID])

	return public_collections.NewGetPublicCollectionsIDOK().WithPayload(&models.APublicCollection{
		Collection: modelCollection,
		Medias:     models.NewModelMediaList(medias),
		Users:      models.NewModelUserList(userList),
	})
}

// GetObjects TBD
func (d *PublicCollections) GetObjects(params public_collections.PostPublicCollectionsObjectsParams) middleware.Responder {
	var (
		logger = d.Context.Log
		DBM    = d.Context.DBM

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return public_collections.NewPostPublicCollectionsObjectsDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}

		successResponse = func() *public_collections.PostPublicCollectionsObjectsOK {
			return public_collections.NewPostPublicCollectionsObjectsOK()
		}

		forbiddenResponse = public_collections.NewPostPublicCollectionsObjectsForbidden
	)

	logger.Debug("GetCollectionObjects")

	// public_collections list
	collectionsIDs := params.RGetCollectionsObjects.CollectionsIds
	collections, err := DBM.GetCollectionsByIDs(collectionsIDs)
	if err != nil {
		logger.Error("GetCollectionsByIDs", "error", err)
		return errorResponse(500, err.Error())
	}
	if len(collections) != len(collectionsIDs) {
		err = errors.New("collection not found")
		logger.Debug("GetCollectionsByIDs", "error", err)
		return errorResponse(404, err.Error())
	}

	// check rights
	for _, c := range collections {
		if !c.Public {
			return forbiddenResponse()
		}
	}

	// remove user for anonymous collection
	for _, c := range collections {
		if c.IsAnonymous {
			c.RootID = 0
		}
	}

	// collections owners
	userRootRefList, err := DBM.GetUserRootRefs(collections.RootsIDs())
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500, err.Error())
	}
	userList, err := DBM.GetUsersByIDs(userRootRefList.GetRootOwnerIDs())
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// get objects cnt
	objectCnt, err := DBM.GetObjectsCnt(collectionsIDs)
	if err != nil {
		logger.Debug("GetObjectsCnt", "error", err)
		return errorResponse(500, err.Error())
	}

	// objects
	orders := models.NewDtoObjectOrder(params.RGetCollectionsObjects.Orders)

	paginator := dto.NewDefaultPagePaginator()
	if paramsPaginator := params.RGetCollectionsObjects.Paginator; paramsPaginator != nil {
		paginator.Cnt = *paramsPaginator.Cnt
		paginator.Page = *paramsPaginator.Page
	}

	objectPreviews, err := DBM.GetCollectionsObjectPreviews(collectionsIDs, orders, paginator)
	if err != nil {
		logger.Debug("GetCollectionObjectPreviews", "err", err)
		return errorResponse(500, "can't get collection")
	}
	objectsPreviewIDs := objectPreviews.GetIDs()

	mhContext := &mhContext{
		DBM: DBM,
		log: logger,
	}

	// medias
	mediaExtractor, err := newObjectMediasExtractor(mhContext, objectsPreviewIDs, &objectMediasExtractorOpts{
		onlyPhoto: true,
	})
	if err != nil {
		logger.Error("GetObjectsMediaRefs", "err", err)
		return errorResponse(500, "can't get medias ref")
	}

	// actors
	actorExtractor, err := newObjectActorsExtractor(mhContext, objectsPreviewIDs)
	if err != nil {
		logger.Error("GetObjectsActorRefs", "err", err)
		return errorResponse(500, "can't get actors ref")
	}

	// originLocations
	originLocationExtractor, err := newObjectOriginLocationsExtractor(mhContext, objectsPreviewIDs)
	if err != nil {
		logger.Error("GetObjectsOriginLocationRefs", "err", err)
		return errorResponse(500, "can't get originLocations ref")
	}

	// mediasIDs to fetch
	var mediasIDs []int64

	// construct object json models
	modelsObjectPreviews := models.NewModelObjectPreviewList(objectPreviews)
	for _, modelObjectPreview := range modelsObjectPreviews.List {
		ID := *modelObjectPreview.ID
		originalLocationsIDs := originLocationExtractor.getOriginLocationsIDsByObjectID(ID)
		objectMedias := mediaExtractor.getMediasIDsByObjectID(ID)
		if len(objectMedias) >= 1 {
			objectMedias = objectMedias[0:1]
			mediasIDs = append(mediasIDs, objectMedias...)
		}
		modelObjectPreview.
			WithMediaIDs(objectMedias).
			WithActorIDs(actorExtractor.getActorsIDsByObjectID(ID)).
			WithOriginLocationIDs(originalLocationsIDs)
	}

	// add users media
	mediasIDs = append(mediasIDs, userList.GetAvatarsMediaIDs()...)

	// add collection media
	mediasIDs = append(mediasIDs, collections.GetMediasIDs()...)

	// fetch medias
	var medias dto.MediaList
	if len(mediasIDs) > 0 {
		medias, err = DBM.GetMediasByIDs(mediasIDs)
		if err != nil {
			logger.Error("get medias", "err", err)
			return errorResponse(500, "can't get additionalMediasIDs")
		}
	}

	// original locations
	originalLocations := originLocationExtractor.getOriginLocations()

	// collections
	modelCollectionsList := models.NewModelCollectionList(collections)

	// payload
	payload := &models.AObjectsPreview{
		TotalObjects:     int64(objectCnt),
		Medias:           models.NewModelMediaList(medias),
		Actors:           models.NewModelActorList(actorExtractor.getActors()),
		OriginLocations:  models.NewModelOriginLocationList(originalLocations),
		ObjectsPreview:   modelsObjectPreviews,
		Collections:      modelCollectionsList,
		Users:            models.NewModelUserList(userList),
		UsersCollections: models.NewModelUsersCollectionsList(collections, userRootRefList),
	}

	return successResponse().WithPayload(payload)
}
