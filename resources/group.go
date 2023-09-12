package resource

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/groups"
)

// Group TBD
type Group struct {
	Context Context
}

// GetGroupByID TBD
func (g *Group) GetGroupByID(params groups.GetGroupsIDParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return groups.NewGetGroupsIDDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
		forbiddenResponse = groups.NewGetGroupsIDForbidden

		DBM           = g.Context.DBM
		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))

		short = false
	)

	if sh := params.Short; sh != nil {
		short = *sh
	}

	groupsList, err := DBM.GetGroupsByIDs([]int64{params.ID})
	if err != nil {
		logger.Error("GetGroupsByIDs", "err", err)
		return errorResponse(500, err.Error())
	}
	if len(groupsList) != 1 {
		logger.Error("GetGroupsByIDs", "err", "group not found")
		return errorResponse(404, fmt.Sprintf("group not find"))
	}
	group := groupsList[0]

	// check group
	ok, err := accessChecker.HasUserRightsForGroups(
		userContext.User.ID,
		dto.RightEntityLevelRead,
		[]int64{group.ID},
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant access group %d", group.ID)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	// collections
	collections, err := DBM.GetCollectionsByGroupIDs([]int64{group.ID})
	if err != nil {
		logger.Error("GetCollectionsByGroupID", "err", err)
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
	collections2cnt, err := DBM.GetObjectsCntByCollections(collectionsIDs)
	if err != nil {
		logger.Error("GetObjectsCntByCollections", "err", err)
		return errorResponse(500, err.Error())
	}

	// objectPreviews
	paginator := dto.NewDefaultPagePaginator()
	if short {
		paginator.Cnt = 0
	}
	orders := dto.NewDefaultObjectOrders()

	objectPreviews, err := DBM.GetCollectionsObjectPreviews(collectionsIDs, orders, paginator)
	if err != nil {
		logger.Error("GetCollectionsObjectPreviews", "err", err)
		return errorResponse(500, err.Error())
	}
	objectsPreviewIDs := objectPreviews.GetIDs()

	// objects medias
	mhContext := &mhContext{
		DBM: DBM,
		log: logger,
	}

	// objects medias
	mediaExtractor, err := newObjectMediasExtractor(mhContext, objectsPreviewIDs, &objectMediasExtractorOpts{
		onlyPhoto: true,
	})
	if err != nil {
		logger.Error("newObjectMediasExtractor", "err", err)
		return errorResponse(500, "can't get medias ref")
	}

	// all medias
	allMediasIDs := append(mediaExtractor.getMediasIDs(), collections.GetMediasIDs()...)
	medias, err := DBM.GetMediasByIDs(allMediasIDs)
	if err != nil {
		logger.Error("GetMediasByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// actors
	actorExtractor, err := newObjectActorsExtractor(mhContext, objectsPreviewIDs)
	if err != nil {
		logger.Error("GetObjectsActorRefs", "err", err)
		return errorResponse(500, "can't get actors ref")
	}

	// badges
	badgeExtractor, err := newObjectBadgesExtractor(mhContext, objectsPreviewIDs)
	if err != nil {
		logger.Error("GetObjectsBadgeRefs", "err", err)
		return errorResponse(500, "can't get badges ref")
	}

	// originLocations
	originLocationExtractor, err := newObjectOriginLocationsExtractor(mhContext, objectsPreviewIDs)
	if err != nil {
		logger.Error("GetObjectsOriginLocationRefs", "err", err)
		return errorResponse(500, "can't get originLocations ref")
	}

	// construct object json models
	modelsObjectPreviews := models.NewModelObjectPreviewList(objectPreviews)
	for _, modelObjectPreview := range modelsObjectPreviews.List {
		ID := *modelObjectPreview.ID
		originalLocationsIDs := originLocationExtractor.getOriginLocationsIDsByObjectID(ID)
		modelObjectPreview.
			WithMediaIDs(mediaExtractor.getMediasIDsByObjectID(ID)).
			WithActorIDs(actorExtractor.getActorsIDsByObjectID(ID)).
			WithActorIDs(badgeExtractor.getBadgesIDsByObjectID(ID)).
			WithAccessLevel(collections2level[*modelObjectPreview.CollectionID]).
			WithOriginLocationIDs(originalLocationsIDs)
	}

	// collections
	collectionModelList := models.NewModelCollectionList(collections)
	for _, collectionModel := range collectionModelList.List {
		collectionModel.
			WithObjectsCnt(collections2cnt[*collectionModel.ID]).
			WithAccessLevel(collections2level[*collectionModel.ID])
	}

	// original locations
	originalLocations := originLocationExtractor.getOriginLocations()

	return groups.NewGetGroupsIDOK().WithPayload(&models.AGroup{
		Actors:          models.NewModelActorList(actorExtractor.getActors()),
		OriginLocations: models.NewModelOriginLocationList(originalLocations),
		Badges:          models.NewModelBadgeList(badgeExtractor.getBadges()),
		Medias:          models.NewModelMediaList(medias),
		Collections:     collectionModelList,
		Group:           models.NewModelGroup(group),
		ObjectsPreview:  modelsObjectPreviews,
	})
}
