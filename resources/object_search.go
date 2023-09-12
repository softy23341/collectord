package resource

import (
	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/objects"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// SearchObjects TBD
func (o *Object) SearchObjects(params objects.PostObjectsSearchParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)
		DBM         = o.Context.DBM

		errorResponse = func(code int) middleware.Responder {
			return objects.NewPostObjectsSearchDefault(code)
		}
		forbiddenResponse = objects.NewPostObjectsSearchForbidden

		// models filters to my filter
		filters = models.NewDtoObjectSearchFilter(params.RSearchObjects.Filters)
		orders  = models.NewDtoObjectOrder(params.RSearchObjects.Orders)

		paginator = &dto.PagePaginator{Page: 0, Cnt: 10}
	)

	if params.RSearchObjects.RootID == 0 {
		return forbiddenResponse()
	}

	logger.Debug("search objects")

	if iPaginator := params.RSearchObjects.Paginator; iPaginator != nil {
		paginator.Page = *iPaginator.Page
		paginator.Cnt = *iPaginator.Cnt
	}

	// remove secret collections
	userCollectionRightList, err := DBM.GetUserRightsForCollectionsInRoot(userContext.User.ID, params.RSearchObjects.RootID)
	if err != nil {
		logger.Error("GetUserRightsForCollections", "err", err)
		return errorResponse(500)
	}
	collectionsIDsToExclude := userCollectionRightList.UnderLevel(dto.RightEntityLevelRead).GetEntityIDs()
	filters.CollectionsToExclude = collectionsIDsToExclude

	var collections2level map[int64]dto.RightEntityLevel
	if filters != nil {
		collectionsIDsToCheck := filters.Collections

		if groupsIDs := params.RSearchObjects.Filters.Groups; len(groupsIDs) != 0 {
			inGroupRefList, err := o.Context.DBM.GetCollectionsGroupRefsByGroups(groupsIDs)
			if err != nil {
				logger.Error("GetCollectionsGroupRefsByGroups", "err", err)
				return errorResponse(500)
			}

			groupID2CollectionsIDs := inGroupRefList.GroupIDToCollectionsIDs()
			for _, collectionsGroup := range groupID2CollectionsIDs {
				filters.CollectionsGroups = append(filters.CollectionsGroups, collectionsGroup)
				collectionsIDsToCheck = append(collectionsIDsToCheck, collectionsGroup...)
			}
		}

		filters.Collections = util.Int64Slice(filters.Collections).DeleteAll(collectionsIDsToExclude)
		collections2level = userCollectionRightList.GetEntityIDToLevel(dto.RightEntityTypeCollection)

		if len(filters.CollectionsGroups) > 0 {
			for _, collections := range filters.CollectionsGroups {
				collections = util.Int64Slice(collections).DeleteAll(collectionsIDsToExclude)
			}
		}
	}

	sr, err := o.Context.SearchClient.SearchObjects(&services.SearchQuery{
		Query:     params.RSearchObjects.Query,
		RootID:    params.RSearchObjects.RootID,
		Filters:   filters,
		Orders:    orders,
		Paginator: paginator,
	})

	if err != nil {
		logger.Error("search objects", "err", err)
		return errorResponse(500)
	}

	objectsPreviewIDs := sr.Objects.IDs()

	// get objects
	objectsPreviews, err := o.Context.DBM.GetObjectsPreviewByIDs(objectsPreviewIDs)
	if err != nil {
		logger.Error("GetObjectsPreviewByIDs", "err", err)
		return errorResponse(500)
	}

	// extractors context
	mhContext := &mhContext{
		DBM: o.Context.DBM,
		log: logger,
	}

	// get actors
	actorRefs, err := o.Context.DBM.GetObjectsActorRefs(objectsPreviewIDs)
	if err != nil {
		o.Context.Log.Error("GetObjectsActorRefs", "err", err)
		return errorResponse(500)
	}
	objectToActorsIDs := actorRefs.ObjectIDToActorsIDs()

	actorsList, err := o.Context.DBM.GetActorsByIDs(sr.Filters.Actors.IDs())
	if err != nil {
		logger.Error("GetActorsByIDs", "err", err)
		return errorResponse(500)
	}

	// originLocations
	originLocationRefs, err := o.Context.DBM.GetObjectsOriginLocationRefs(objectsPreviewIDs)
	if err != nil {
		logger.Error("GetObjectsOriginLocationRefs", "err", err)
		return errorResponse(500)
	}
	objectToOriginLocationsIDs := originLocationRefs.ObjectIDToOriginLocationsIDs()

	originLocationsList, err := o.Context.DBM.GetOriginLocationsByIDs(sr.Filters.OriginLocations.IDs())
	if err != nil {
		logger.Error("GetOriginLocationsByIDs", "err", err)
		return errorResponse(500)
	}

	// badges
	badgeRefs, err := o.Context.DBM.GetObjectsBadgeRefs(objectsPreviewIDs)
	if err != nil {
		logger.Error("GetObjectsBadgeRefs", "err", err)
		return errorResponse(500)
	}
	objectIDTobadgesIDs := badgeRefs.ObjectIDToBadgesIDs()

	badgeList, err := o.Context.DBM.GetBadgesByIDs(sr.Filters.Badges.IDs())
	if err != nil {
		logger.Error("GetBadgesByIDs", "err", err)
		return errorResponse(500)
	}

	// materials
	materials, err := o.Context.DBM.GetMaterialsByIDs(sr.Filters.Materials.IDs())
	if err != nil {
		logger.Error("get materials by ids", "err", err)
		return errorResponse(500)
	}

	// collections
	collections, err := o.Context.DBM.GetCollectionsByIDs(sr.Filters.Collections.IDs())
	if err != nil {
		logger.Error("GetCollectionsByIDs", "err", err)
		return errorResponse(500)
	}
	collectionsIDs := collections.GetIDs()

	collections2cnt, err := o.Context.DBM.GetObjectsCntByCollections(collectionsIDs)
	if err != nil {
		logger.Error("GetObjectsCntByCollections", "err", err)
		return errorResponse(500)
	}

	// collection groups ref
	groupRefs, err := o.Context.DBM.GetCollectionsGroupRefs(collectionsIDs)
	if err != nil {
		logger.Error("GetCollectionsGroupRefs", "err", err)
		return errorResponse(500)
	}

	// groups
	groupList, err := o.Context.DBM.GetGroupsByIDs(groupRefs.GroupsIDs())
	if err != nil {
		logger.Error("GetGroupsByIDs", "err", err)
		return errorResponse(500)
	}

	// objects medias
	mediaExtractor, err := newObjectMediasExtractor(mhContext, objectsPreviewIDs, &objectMediasExtractorOpts{
		onlyPhoto: true,
	})
	if err != nil {
		logger.Error("newObjectMediasExtractor", "err", err)
		return errorResponse(500)
	}

	// statuses
	objectStatuses, err := o.Context.DBM.GetCurrentObjectsStatusesRefs(objectsPreviewIDs)
	if err != nil {
		logger.Error("GetCurrentObjectsStatusesRefs", "err", err)
		return errorResponse(500)
	}
	objectStatusMap := objectStatuses.ObjectToOneStatusMap()

	// medias ids
	var mediasIDs []int64
	mediasIDs = append(mediasIDs, collections.GetMediasIDs()...)

	// named date intervals
	var namedDateIntervalsIDs []int64
	namedDateIntervalsIDs = append(namedDateIntervalsIDs, sr.Filters.ProdutcionNamedIntervals.IDs()...)
	namedDateIntervals, err := o.Context.DBM.GetNamedDateIntervalsByIDs(namedDateIntervalsIDs)
	if err != nil {
		logger.Error("GetNamedDateIntervalsByIDs", "err", err)
		return errorResponse(500)
	}

	// objects order
	// TODO remove silly resort after fetch objects from DB
	properOrderedPreviews := make(dto.ObjectPreviewList, 0, len(objectsPreviewIDs))
	for _, sObjectID := range objectsPreviewIDs {
		for _, object := range objectsPreviews {
			if object.ID == sObjectID {
				properOrderedPreviews = append(properOrderedPreviews, object)
				break
			}
		}

	}

	// valuations
	valuationsExtractor, err := newObjectValuationsExtractor(mhContext, objectsPreviewIDs)
	if err != nil {
		logger.Error("GetValuations", "err", err)
		return errorResponse(500)
	}

	objectValuationMap := valuationsExtractor.GetObjectValuations().ObjectToOneValuationMap()

	// construct object json models
	modelsObjectPreviews := models.NewModelObjectPreviewList(properOrderedPreviews)
	for _, modelObjectPreview := range modelsObjectPreviews.List {
		ID := *modelObjectPreview.ID
		originalLocationsIDs := objectToOriginLocationsIDs[ID]
		objectMedias := mediaExtractor.getMediasIDsByObjectID(ID)
		if len(objectMedias) >= 1 {
			objectMedias = objectMedias[0:1]
			mediasIDs = append(mediasIDs, objectMedias...)
		}
		modelObjectPreview.
			WithMediaIDs(objectMedias).
			WithActorIDs(objectToActorsIDs[ID]).
			WithOriginLocationIDs(originalLocationsIDs).
			WithBadgeIDs(objectIDTobadgesIDs[ID]).
			WithStatus(objectStatusMap[ID]).
			WithAccessLevel(collections2level[*modelObjectPreview.CollectionID]).
			WithValuations(models.NewModelValuationList(objectValuationMap[ID]))
	}

	// medias
	medias, err := o.Context.DBM.GetMediasByIDs(mediasIDs)
	if err != nil {
		logger.Error("GetMediasByIDs", "err", err)
		return errorResponse(500)
	}

	// construct collections
	modelsCollections := models.NewModelCollectionList(collections)
	collection2groups := groupRefs.CollectionIDToGroupsIDs()
	for _, modelCollection := range modelsCollections.List {
		modelCollection.
			WithGroupsIDs(collection2groups[*modelCollection.ID]).
			WithObjectsCnt(collections2cnt[*modelCollection.ID])
	}

	// manual construct stats by groups
	groupID2cnt := make(map[int64]int64, len(groupList))
	for _, collectionStat := range sr.Filters.Collections {
		groupIDs, find := collection2groups[collectionStat.PropertyID]
		if !find {
			continue
		}
		for _, groupID := range groupIDs {
			groupID2cnt[groupID] += collectionStat.Cnt
		}
	}

	groupsStat := make(dto.ObjectSearchFiltersResultList, 0, len(groupID2cnt))
	for groupID, cnt := range groupID2cnt {
		groupsStat = append(groupsStat, &dto.ObjectSearchFiltersResult{
			PropertyID: groupID,
			Cnt:        cnt,
		})
	}

	return objects.NewPostObjectsSearchOK().WithPayload(&models.ASearchObjects{
		Filters:         models.NewModelObjectSearchFilter(sr.Filters).WithGroups(groupsStat),
		Actors:          models.NewModelActorList(actorsList),
		OriginLocations: models.NewModelOriginLocationList(originLocationsList),
		Badges:          models.NewModelBadgeList(badgeList),
		Materials:       models.NewModelMaterialList(materials),
		Collections:     modelsCollections,
		Medias:          models.NewModelMediaList(medias),
		DateIntervals:   models.NewModelDateIntervalList(namedDateIntervals),
		ObjectsPreview:  modelsObjectPreviews,
		Groups:          models.NewModelGroupList(groupList),
		TotalHits:       sr.TotalHits,
	})
}
