package resource

import (
	"errors"
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"gopkg.in/inconshreveable/log15.v2"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dal"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/collections"
)

// Collection TBD
type Collection struct {
	Context Context
}

// GetCollectionsObjects TBD
func (c *Collection) GetCollectionsObjects(params collections.PostCollectionsObjectsParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)
		DBM         = c.Context.DBM

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return collections.NewPostCollectionsObjectsDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}

		successResponse = func() *collections.PostCollectionsObjectsOK {
			return collections.NewPostCollectionsObjectsOK()
		}

		forbiddenResponse = collections.NewPostCollectionsObjectsForbidden

		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))
	)

	logger.Debug("GetCollectionObjects")

	// collections list
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

	// check collection
	ok, err := accessChecker.HasUserRightsForCollections(
		userContext.User.ID,
		dto.RightEntityLevelRead,
		collections.GetIDs(),
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant read collection %d", collections.GetIDs())
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	// get rights
	userCollectionRightList, err := DBM.GetUserRightsForCollections(userContext.User.ID, collections.GetIDs())
	if err != nil {
		logger.Error("GetUserRightsForCollections", "err", err)
		return errorResponse(500, err.Error())
	}
	collections2level := userCollectionRightList.GetEntityIDToLevel(dto.RightEntityTypeCollection)

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

	// badges
	badgeExtractor, err := newObjectBadgesExtractor(mhContext, objectsPreviewIDs)
	if err != nil {
		logger.Error("GetObjectsBadgeRefs", "err", err)
		return errorResponse(500, "can't get badges ref")
	}

	// statuses
	objectStatuses, err := DBM.GetCurrentObjectsStatusesRefs(objectsPreviewIDs)
	if err != nil {
		logger.Error("GetCurrentObjectsStatusesRefs", "err", err)
		return errorResponse(500, "can't get GetCurrentObjectsStatusesRefs")
	}
	objectStatusMap := objectStatuses.ObjectToOneStatusMap()

	// valuations
	valuationsExtractor, err := newObjectValuationsExtractor(mhContext, objectsPreviewIDs)
	if err != nil {
		logger.Error("GetValuations", "err", err)
		return errorResponse(500, "can't get object valuations")
	}

	objectValuationMap := valuationsExtractor.GetObjectValuations().ObjectToOneValuationMap()

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
			WithBadgeIDs(badgeExtractor.getBadgesIDsByObjectID(ID)).
			WithOriginLocationIDs(originalLocationsIDs).
			WithStatus(objectStatusMap[ID]).
			WithAccessLevel(collections2level[*modelObjectPreview.CollectionID]).
			WithValuations(models.NewModelValuationList(objectValuationMap[ID]))
	}

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

	// payload
	payload := &models.AObjectsPreview{
		TotalObjects:    int64(objectCnt),
		Medias:          models.NewModelMediaList(medias),
		Actors:          models.NewModelActorList(actorExtractor.getActors()),
		Badges:          models.NewModelBadgeList(badgeExtractor.getBadges()),
		OriginLocations: models.NewModelOriginLocationList(originalLocations),
		ObjectsPreview:  modelsObjectPreviews,
	}

	return successResponse().WithPayload(payload)
}

// GetDraftCollection TBD
func (c *Collection) GetDraftCollection(params collections.GetCollectionsDraftParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)

	logger.Debug("GetDraftCollection")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return collections.NewGetCollectionsDraftDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	return errorResponse(500, "deprecated")

	rootID := params.RootID

	DBM := c.Context.DBM
	collection, err := DBM.GetTypedCollection(rootID, dto.DraftCollectionTypo)
	if err != nil {
		logger.Error("GetDraftCollection", "err", err)
		return errorResponse(500, err.Error())
	}

	objectCnt, err := DBM.GetObjectsCnt([]int64{collection.ID})
	if err != nil {
		logger.Debug("GetObjectsCnt", "error", err)
		return errorResponse(500, err.Error())
	}

	return collections.
		NewGetCollectionsDraftOK().
		WithPayload(models.NewModelCollection(collection).WithObjectsCnt(objectCnt))
}

type mhContext struct {
	DBM dal.TrManager
	log log15.Logger
}

// media
type objectMediasExtractor struct {
	medias             dto.MediaList
	refs               dto.ObjectMediaRefList
	objectID2mediasIDs map[int64][]int64
	context            *mhContext
}

type objectMediasExtractorOpts struct {
	onlyPhoto bool
}

func newObjectMediasExtractor(context *mhContext, objectIDs []int64, opts *objectMediasExtractorOpts) (*objectMediasExtractor, error) {
	// medias
	var (
		mediaRefs dto.ObjectMediaRefList
		err       error
	)
	if opts.onlyPhoto {
		mediaRefs, err = context.DBM.GetObjectsMediaRefsPhoto(objectIDs)
	} else {
		mediaRefs, err = context.DBM.GetObjectsMediaRefs(objectIDs)
	}
	if err != nil {
		context.log.Error("GetObjectsMediaRefs", "err", err)
		return nil, err
	}

	return &objectMediasExtractor{
		refs:               mediaRefs,
		objectID2mediasIDs: mediaRefs.ObjectIDToMediasIDs(),
		context:            context,
	}, nil
}

func (o *objectMediasExtractor) getMedias() (dto.MediaList, error) {
	medias, err := o.context.DBM.GetMediasByIDs(o.refs.UniqMediasIDs())
	if err != nil {
		o.context.log.Error("GetMediasByIDs", "err", err)
		return nil, err
	}
	o.medias = medias

	return o.medias, nil
}

func (o *objectMediasExtractor) getMediasIDs() []int64 {
	return o.refs.UniqMediasIDs()
}

func (o *objectMediasExtractor) getMediasIDsByObjectID(objectID int64) []int64 {
	return o.objectID2mediasIDs[objectID]
}

// actor
type objectActorsExtractor struct {
	actors             dto.ActorList
	refs               dto.ObjectActorRefList
	objectID2actorsIDs map[int64][]int64
	context            *mhContext
}

func newObjectActorsExtractor(context *mhContext, objectIDs []int64) (*objectActorsExtractor, error) {
	// actors
	actorRefs, err := context.DBM.GetObjectsActorRefs(objectIDs)
	if err != nil {
		context.log.Error("GetObjectsActorRefs", "err", err)
		return nil, err
	}

	actors, err := context.DBM.GetActorsByIDs(actorRefs.ActorsIDs())
	if err != nil {
		context.log.Error("GetActorsByIDs", "err", err)
		return nil, err
	}

	return &objectActorsExtractor{
		refs:               actorRefs,
		objectID2actorsIDs: actorRefs.ObjectIDToActorsIDs(),
		actors:             actors,
		context:            context,
	}, nil
}

func (o *objectActorsExtractor) getActors() dto.ActorList {
	return o.actors
}

func (o *objectActorsExtractor) getActorsIDsByObjectID(objectID int64) []int64 {
	return o.objectID2actorsIDs[objectID]
}

// origin location
type objectOriginLocationsExtractor struct {
	originLocations             dto.OriginLocationList
	refs                        dto.ObjectOriginLocationRefList
	objectID2originLocationsIDs map[int64][]int64
	context                     *mhContext
}

func newObjectOriginLocationsExtractor(context *mhContext, objectIDs []int64) (*objectOriginLocationsExtractor, error) {
	// originLocations
	originLocationRefs, err := context.DBM.GetObjectsOriginLocationRefs(objectIDs)
	if err != nil {
		context.log.Error("GetObjectsOriginLocationRefs", "err", err)
		return nil, err
	}
	originLocations, err := context.DBM.GetOriginLocationsByIDs(originLocationRefs.OriginLocationsIDs())
	if err != nil {
		context.log.Error("GetOriginLocationsByIDs", "err", err)
		return nil, err
	}

	return &objectOriginLocationsExtractor{
		refs:                        originLocationRefs,
		objectID2originLocationsIDs: originLocationRefs.ObjectIDToOriginLocationsIDs(),
		originLocations:             originLocations,
		context:                     context,
	}, nil
}

func (o *objectOriginLocationsExtractor) getOriginLocations() dto.OriginLocationList {
	return o.originLocations
}

func (o *objectOriginLocationsExtractor) getOriginLocationsIDsByObjectID(objectID int64) []int64 {
	return o.objectID2originLocationsIDs[objectID]
}

// badge
type objectBadgesExtractor struct {
	badges             dto.BadgeList
	refs               dto.ObjectBadgeRefList
	objectID2badgesIDs map[int64][]int64
	context            *mhContext
}

func newObjectBadgesExtractor(context *mhContext, objectIDs []int64) (*objectBadgesExtractor, error) {
	// badges
	badgeRefs, err := context.DBM.GetObjectsBadgeRefs(objectIDs)
	if err != nil {
		context.log.Error("GetObjectsBadgeRefs", "err", err)
		return nil, err
	}
	badges, err := context.DBM.GetBadgesByIDs(badgeRefs.BadgesIDs())
	if err != nil {
		context.log.Error("GetBadgesByIDs", "err", err)
		return nil, err
	}

	return &objectBadgesExtractor{
		refs:               badgeRefs,
		objectID2badgesIDs: badgeRefs.ObjectIDToBadgesIDs(),
		badges:             badges,
		context:            context,
	}, nil
}

func (o *objectBadgesExtractor) getBadges() dto.BadgeList {
	return o.badges
}

func (o *objectBadgesExtractor) getBadgesIDsByObjectID(objectID int64) []int64 {
	return o.objectID2badgesIDs[objectID]
}

// valuation
type objectValuationsExtractor struct {
	valuations dto.ValuationList
	context    *mhContext
}

func newObjectValuationsExtractor(context *mhContext, objectIDs []int64) (*objectValuationsExtractor, error) {
	// valuations
	valuations, err := context.DBM.GetObjectValuations(objectIDs)
	if err != nil {
		context.log.Error("GetObjectValuations", "err", err)
		return nil, err
	}

	return &objectValuationsExtractor{
		valuations: valuations,
		context:    context,
	}, nil
}

func (o *objectValuationsExtractor) GetObjectValuations() dto.ValuationList {
	return o.valuations
}
