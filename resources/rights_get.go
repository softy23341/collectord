package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/rights"
	"github.com/go-openapi/runtime/middleware"
)

// GetRights TBD
func (uer *UserEntityRight) GetRights(params rights.GetRightParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)
		DBM         = uer.Context.DBM

		defaultErrorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return rights.NewGetRightDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}

		forbiddenResponse = rights.NewGetRightForbidden

		accessChecker = NewAccessRightsChecker(DBM, logger.New("service", "access rights"))

		targetRootID = params.RootID
		targetUserID = params.UserID
	)

	if ok, err := accessChecker.IsUserRelatedToRoots(userContext.User.ID, []int64{targetRootID}); err != nil {
		logger.Error("cant check user relation to root", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else if !ok {
		err := fmt.Errorf("user has not rights for this root: %d", targetRootID)
		logger.Error("access error", "err", err)
		return forbiddenResponse()
	}

	var userMainRoot *dto.Root
	if rootsList, err := DBM.GetMainUserRoot(userContext.User.ID); err != nil {
		logger.Error("GetMainUserRoot", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else if len(rootsList) != 1 {
		err := fmt.Errorf("cant find user main root")
		logger.Error("GetMainUserRoot", "err", err)
		return defaultErrorResponse(500, err.Error())
	} else {
		userMainRoot = rootsList[0]
	}

	userIsOwner := userMainRoot.ID == targetRootID
	selfAwareUser := targetUserID == userContext.User.ID

	if !userIsOwner && !selfAwareUser {
		logger.Error("user is not self aware and is not owner")
		return forbiddenResponse()
	}

	rightsList, err := DBM.GetUserRightsForCollectionsInRoot(targetUserID, targetRootID)
	if err != nil {
		logger.Error("cant get rights for user in root", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	if selfAwareUser && !userIsOwner {
		rightsList = rightsList.UnderLevel(dto.RightEntityLevelRead)
	}

	level2collection := rightsList.GetEntityIDToLevel(dto.RightEntityTypeCollection)
	collectionsIDs := rightsList.GetEntityIDs()
	collectionsList, err := DBM.GetCollectionsByIDs(collectionsIDs)
	if err != nil {
		logger.Error("cant getcollectionsbyids", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	mediaList, err := DBM.GetMediasByIDs(collectionsList.GetMediasIDs())
	if err != nil {
		logger.Error("cant getmediasids", "err", err)
		return defaultErrorResponse(500, err.Error())
	}

	collectionsModel := models.NewModelCollectionList(collectionsList)
	for _, collectionModel := range collectionsModel.List {
		collectionModel.WithAccessLevel(level2collection[*collectionModel.ID])
	}

	return rights.NewGetRightOK().WithPayload(&models.AGetRight{
		Collections: collectionsModel,
		Medias:      models.NewModelMediaList(mediaList),
		Rights:      models.NewModelShortUserEntityRightList(rightsList),
	})
}
