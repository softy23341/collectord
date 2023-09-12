package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/dashboard"
	"github.com/go-openapi/runtime/middleware"
)

// Dashboard TBD
type Dashboard struct {
	Context Context
}

// GetDashboard TBD
func (d *Dashboard) GetDashboard(params dashboard.GetDashboardParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("get dashboard")

	var (
		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return dashboard.NewGetDashboardDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		forbiddenResponse = dashboard.NewGetDashboardForbidden

		DBM           = d.Context.DBM
		accessChecker = NewAccessRightsChecker(
			DBM,
			logger.New("service", "access rights"),
		)

		rootID                  = params.RootID
		objectsCntPerCollection = params.ObjectsCntPerCollection
	)

	if ok, err := accessChecker.IsUserRelatedToRoots(userContext.User.ID, []int64{rootID}); err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user is not related to root: %d", rootID)
		logger.Error("user is not allowed to be here", "err", err)
		return forbiddenResponse()
	}

	// groups
	groups, err := DBM.GetGroupsByRootID(rootID)
	if err != nil {
		logger.Error("GetGroupsByRootID", "err", err)
		return errorResponse(500, err.Error())
	}
	userGroupRightList, err := DBM.GetUserRightsForGroups(userContext.User.ID, groups.GetIDs())
	if err != nil {
		logger.Error("GetUserRightsForGroups", "err", err)
		return errorResponse(500, err.Error())
	}
	groups = groups.RemoveByIDs(
		userGroupRightList.UnderLevel(dto.RightEntityLevelRead).GetEntityIDs(),
	)
	groups2level := userGroupRightList.GetEntityIDToLevel(dto.RightEntityTypeGroup)

	// collections
	collections, err := DBM.GetCollectionsByRootID(rootID)
	if err != nil {
		logger.Error("GetCollectionsByRootID", "err", err)
		return errorResponse(500, err.Error())
	}
	userCollectionRightList, err := DBM.GetUserRightsForCollections(userContext.User.ID, collections.GetIDs())
	if err != nil {
		logger.Error("GetUserRightsForCollections", "err", err)
		return errorResponse(500, err.Error())
	}
	collections = collections.RemoveByIDs(
		userCollectionRightList.UnderLevel(dto.RightEntityLevelRead).GetEntityIDs(),
	)
	collections2level := userCollectionRightList.GetEntityIDToLevel(dto.RightEntityTypeCollection)
	collectionsIDs := collections.GetIDs()

	// collections cnt
	collections2cnt, err := DBM.GetObjectsCntByCollections(collectionsIDs)
	if err != nil {
		logger.Error("GetObjectsCntByCollections", "err", err)
		return errorResponse(500, err.Error())
	}

	// collections valuations
	collections2valuation, err := DBM.GetValuationsByCollectionIDs(collectionsIDs)
	if err != nil {
		logger.Error("GetValuationsByCollectionIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// objectPreviews
	var objectsPreviewsList dto.ObjectPreviewList

	paginator := dto.NewDefaultPagePaginator()
	if objectsCntPerCollection != nil {
		paginator.Cnt = *objectsCntPerCollection
	}
	orders := dto.NewDefaultObjectOrders()

	for _, collection := range collections {
		objectPreviews, err := DBM.GetCollectionsObjectPreviews([]int64{collection.ID}, orders, paginator)
		if err != nil {
			logger.Error("GetCollectionsObjectPreviews", "err", err)
			return errorResponse(500, err.Error())
		}

		objectsPreviewsList = append(objectsPreviewsList, objectPreviews...)
	}

	// collection groups ref
	groupRefs, err := DBM.GetCollectionsGroupRefs(collectionsIDs)
	if err != nil {
		logger.Error("GetCollectionsGroupRefs", "err", err)
		return errorResponse(500, "can't get groups refs")
	}

	objectsPreviewIDs := objectsPreviewsList.GetIDs()
	objectMeta, err := NewObjectPreviewExtractor(d.Context.DBM, logger.New("object extractor"), &ObjectPreviewExtractorOpts{
		mediaExtractorOpts: &objectMediasExtractorOpts{
			onlyPhoto: true,
		},
	}).
		SetObjectsIDs(objectsPreviewIDs).
		FetchAll().
		Result()

	if err != nil {
		logger.Error("NewObjectPreviewExtractor", "err", err)
		return errorResponse(500, err.Error())
	}

	// fetch all medias
	var allMediasIDs []int64
	allMediasIDs = append(collections.GetMediasIDs(), allMediasIDs...)

	objectStatusMap := objectMeta.ObjectStatuses.ObjectToOneStatusMap()
	// construct object json models
	modelsObjectPreviews := models.NewModelObjectPreviewList(objectsPreviewsList)
	for _, modelObjectPreview := range modelsObjectPreviews.List {
		ID := *modelObjectPreview.ID
		originalLocationsIDs := objectMeta.OriginLocationExtractor.getOriginLocationsIDsByObjectID(ID)
		objectMedias := objectMeta.MediaExtractor.getMediasIDsByObjectID(ID)
		if len(objectMedias) >= 1 {
			objectMedias = objectMedias[0:1]
			allMediasIDs = append(allMediasIDs, objectMedias...)
		}
		modelObjectPreview.
			WithMediaIDs(objectMedias).
			WithActorIDs(objectMeta.ActorExtractor.getActorsIDsByObjectID(ID)).
			WithOriginLocationIDs(originalLocationsIDs).
			WithBadgeIDs(objectMeta.BadgeExtractor.getBadgesIDsByObjectID(ID)).
			WithStatus(objectStatusMap[ID]).
			WithAccessLevel(collections2level[*modelObjectPreview.CollectionID])
	}

	// collections
	collection2groups := groupRefs.CollectionIDToGroupsIDs()

	modelcollectionsList := models.NewModelCollectionList(collections)
	for _, modelCollection := range modelcollectionsList.List {
		collectionID := *modelCollection.ID
		modelCollection.
			WithGroupsIDs(collection2groups[collectionID]).
			WithObjectsCnt(collections2cnt[collectionID]).
			WithAccessLevel(collections2level[collectionID]).
			WithValuation(&models.CollectionValuation{
				Price:      collections2valuation[collectionID],
				CurrencyID: 1,
			})
	}

	// groups
	groups2valuation := make(map[int64]int64)
	groups2collection := groupRefs.GroupIDToCollectionsIDs()
	for groupID := range groups2collection {
		for _, colID := range groups2collection[groupID] {
			groups2valuation[groupID] += collections2valuation[colID]
		}
	}

	modelgroupsList := models.NewModelGroupList(groups)
	for _, modelGroup := range modelgroupsList.List {
		groupID := *modelGroup.ID
		modelGroup.
			WithAccessLevel(groups2level[groupID]).
			WithValuation(&models.GroupValuation{
				Price:      groups2valuation[groupID],
				CurrencyID: 1,
			})
	}

	// fetch all medias
	medias, err := DBM.GetMediasByIDs(allMediasIDs)
	if err != nil {
		logger.Error("GetMediasByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	originLocations := objectMeta.OriginLocationExtractor.getOriginLocations()
	return dashboard.NewGetDashboardOK().WithPayload(&models.ADashboard{
		Collections:     modelcollectionsList,
		Groups:          modelgroupsList,
		Medias:          models.NewModelMediaList(medias),
		ObjectsPreview:  modelsObjectPreviews,
		Actors:          models.NewModelActorList(objectMeta.ActorExtractor.getActors()),
		Badges:          models.NewModelBadgeList(objectMeta.BadgeExtractor.getBadges()),
		OriginLocations: models.NewModelOriginLocationList(originLocations),
	})
}
