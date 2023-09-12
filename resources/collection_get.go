package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/collections"
	"github.com/go-openapi/runtime/middleware"
)

// GetCollection TBD
func (c *Collection) GetCollection(params collections.GetCollectionsIDParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)
		DBM         = c.Context.DBM

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return collections.NewGetCollectionsIDDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}

		forbiddenResponse = collections.NewGetCollectionsIDForbidden

		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))

		collectionsIDs = []int64{params.ID}
	)

	logger.Debug("GetCollection")

	// collections valuations
	collections2valuation, err := c.Context.DBM.GetValuationsByCollectionIDs(collectionsIDs)
	if err != nil {
		logger.Error("GetValuationsByCollectionIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// collections
	var collection *dto.Collection
	if collectionsList, err := c.Context.DBM.GetCollectionsByIDs(collectionsIDs); err != nil {
		logger.Error("get collections by ids", "err", err)
		return errorResponse(500, err.Error())
	} else if len(collectionsList) == 0 {
		err := fmt.Errorf("cant find collection: %d", params.ID)
		logger.Error("err", "err", err)
		return collections.NewGetCollectionsIDNotFound()
	} else {
		collection = collectionsList[0]
	}

	// check rights
	ok, err := accessChecker.HasUserRightsForCollections(
		userContext.User.ID,
		dto.RightEntityLevelRead,
		[]int64{collection.ID},
	)
	if err != nil {
		logger.Error("cant check rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("cant read collection %d", collection.ID)
		logger.Error("access denied", "err", err.Error())
		return forbiddenResponse()
	}

	// users rights to collection
	rightsList, err := DBM.GetUsersRightsForCollection(collection.ID)
	if err != nil {
		logger.Error("cant get rights for collection", "err", err)
		return errorResponse(500, err.Error())
	}
	userID2level := rightsList.GetUserIDToLevel()

	// collections cnt
	collections2cnt, err := c.Context.DBM.GetObjectsCntByCollections(collectionsIDs)
	if err != nil {
		logger.Error("GetObjectsCntByCollections", "err", err)
		return errorResponse(500, err.Error())
	}

	// collection groups ref
	groupRefs, err := c.Context.DBM.GetCollectionsGroupRefs(collectionsIDs)
	if err != nil {
		logger.Error("GetCollectionsGroupRefs", "err", err)
		return errorResponse(500, "can't get groups refs")
	}

	// groups
	groupsList, err := c.Context.DBM.GetGroupsByIDs(groupRefs.GroupsIDs())
	if err != nil {
		logger.Error("get groups", "err", err)
		return errorResponse(500, "can't get groups")
	}
	userGroupRightList, err := DBM.GetUserRightsForGroups(userContext.User.ID, groupsList.GetIDs())
	if err != nil {
		logger.Error("GetUserRightsForGroups", "err", err)
		return errorResponse(500, err.Error())
	}
	groupsList = groupsList.RemoveByIDs(
		userGroupRightList.UnderLevel(dto.RightEntityLevelRead).GetEntityIDs(),
	)
	groups2level := userGroupRightList.GetEntityIDToLevel(dto.RightEntityTypeGroup)

	var allMediasIDs []int64
	if img := collection.ImageMediaID; img != nil {
		allMediasIDs = append(allMediasIDs, *img)
	}

	// users
	userList, err := c.Context.DBM.GetUsersByIDs(rightsList.GetUsersIDs())
	if err != nil {
		logger.Error("cant get users by ids", "err", err)
		return errorResponse(500, err.Error())
	}
	allMediasIDs = append(allMediasIDs, userList.GetAvatarsMediaIDs()...)

	// medias
	medias, err := c.Context.DBM.GetMediasByIDs(allMediasIDs)
	if err != nil {
		logger.Error("GetMediasByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// collections
	modelCollection := models.NewModelCollection(collection)
	modelCollection.
		WithGroupsIDs(groupsList.GetIDs()).
		WithObjectsCnt(collections2cnt[params.ID]).
		WithAccessLevel(userID2level[userContext.User.ID]).
		WithValuation(&models.CollectionValuation{
			Price:      collections2valuation[params.ID],
			CurrencyID: 1,
		})

	modelgroupsList := models.NewModelGroupList(groupsList)
	for _, modelGroup := range modelgroupsList.List {
		groupID := *modelGroup.ID
		modelGroup.
			WithAccessLevel(groups2level[groupID])
	}

	return collections.NewGetCollectionsIDOK().WithPayload(&models.ACollection{
		Collection: modelCollection,
		Groups:     modelgroupsList,
		Medias:     models.NewModelMediaList(medias),
		Users:      models.NewModelUserList(userList),
		Rights:     models.NewModelEntityAccessRightList(rightsList),
	})
}
