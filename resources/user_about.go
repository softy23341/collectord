package resource

import (
	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/users"
	"github.com/go-openapi/runtime/middleware"
)

// UserInfo TBD
func (u *User) UserInfo(params users.GetUserAboutParams, principal interface{}) middleware.Responder {

	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)
		DBM         = u.Context.DBM
		currentUser = userContext.User

		errorResponse = func(code int) middleware.Responder {
			return users.NewGetUserAboutDefault(code)
		}
		forbiddenResponse = users.NewGetUserAboutForbidden

		targetUserID = userContext.User.ID
		targetRootID = params.RootID

		targetUser *dto.User
		mediaIDs   []int64

		selfAware = false
		rootAdmin = false
	)

	logger.Debug("user into")

	// check root context
	var currentUserRoot *dto.Root
	if rootList, err := DBM.GetMainUserRoot(userContext.User.ID); err != nil {
		logger.Error("cant get user root", "err", err)
		return errorResponse(500)
	} else if len(rootList) != 1 {
		logger.Error("cant find user root")
		return errorResponse(404)
	} else {
		currentUserRoot = rootList[0]
	}

	if targetRootID != nil {
		if *targetRootID != currentUserRoot.ID {
			logger.Error("cant access to root")
			return forbiddenResponse()
		}
		rootAdmin = true
	}

	// check params
	if userID := params.ID; userID != nil && *userID != currentUser.ID {
		if !rootAdmin {
			logger.Error("you cant get info about user without being rootAdmin")
			return forbiddenResponse()
		}

		targetUserID = *userID
	} else if userID == nil || (userID != nil && *userID == currentUser.ID) {
		selfAware = true
	}

	if !rootAdmin && !selfAware {
		logger.Error("rootadmin and selfAware")
		return forbiddenResponse()
	}

	// get target user
	if usersList, err := DBM.GetUsersByIDs([]int64{targetUserID}); err != nil {
		logger.Error("DBM.GetUsersByIDs", "err", err)
		return errorResponse(500)
	} else if len(usersList) != 1 {
		logger.Error("DBM.GetUsersByIDs cant find user")
		return errorResponse(404)
	} else {
		targetUser = usersList[0]
	}

	var (
		userRootRefs dto.UserRootRefList
		err          error
		usersIDs     = []int64{targetUserID}
		rootToUsers  map[int64]*dto.RootUsers
	)

	if selfAware {
		// all relations
		userRootRefs, err = DBM.GetUserRelatedRootRefs([]int64{currentUser.ID})
		if err != nil {
			logger.Error("GetUserRelatedRootRefs", "err", err)
			return errorResponse(500)
		}
	} else if rootAdmin && targetRootID != nil {
		// only target
		userRootRefs, err = DBM.GetUserRootRefs([]int64{*targetRootID})
		if err != nil {
			logger.Error("GetUserRelatedRootRefs", "err", err)
			return errorResponse(500)
		}
	}

	usersIDs = append(usersIDs, userRootRefs.UsersList()...)
	rootToUsers = userRootRefs.RootToUsers()

	// get users
	rootsUsersList, err := DBM.GetUsersByIDs(usersIDs)
	if err != nil {
		logger.Error("rootsUsersList", "err", err)
		return errorResponse(500)
	}

	// medias
	mediaIDs = append(mediaIDs, rootsUsersList.GetAvatarsMediaIDs()...)
	if targetUser.AvatarMediaID != nil {
		mediaIDs = append(mediaIDs, *targetUser.AvatarMediaID)
	}

	var (
		rightsList     dto.ShortUserEntityRightList
		collectionsIDs []int64
	)

	if rootAdmin && targetRootID != nil {
		rightsList, err = DBM.GetUserRightsForCollectionsInRoot(targetUserID, *targetRootID)
		if err != nil {
			logger.Error("cant get rights for user in root", "err", err)
			return errorResponse(500)
		}

		collectionsIDs = rightsList.GetEntityIDs()
	}
	level2collection := rightsList.GetEntityIDToLevel(dto.RightEntityTypeCollection)

	collectionsList, err := DBM.GetCollectionsByIDs(collectionsIDs)
	if err != nil {
		logger.Error("cant getcollectionsbyids", "err", err)
		return errorResponse(500)
	}
	mediaIDs = append(mediaIDs, collectionsList.GetMediasIDs()...)

	// tasks
	var tasksList dto.TaskList
	if selfAware {
		tasksList, err = DBM.GetUserRelatedTasks(currentUser.ID)
	} else if rootAdmin {
		tasksList, err = DBM.GetAssignedToUserTasksFrom(currentUser.ID, targetUserID)
	}

	// fetch medias
	var medias dto.MediaList
	if len(mediaIDs) > 0 {
		var err error
		medias, err = DBM.GetMediasByIDs(mediaIDs)
		if err != nil {
			logger.Error("GetMediasByIDs(mediaIDs)", "err", err)
			return errorResponse(500)
		}
	}

	// get roots
	rootList, err := DBM.GetRootsByIDs(userRootRefs.GetRootIDs())
	if err != nil {
		logger.Error("final get roots", "err", err)
		return errorResponse(500)
	}

	// roots models
	rootModelList := models.NewModelRootList(rootList)
	for _, rootModel := range rootModelList.List {
		rootUsers := rootToUsers[*rootModel.ID]
		if rootUsers != nil {
			rootModel.
				WithOwnerID(rootUsers.OwnerID).
				WithUsersIDs(rootUsers.UsersIDs)
		}
	}

	// collection model
	collectionsModel := models.NewModelCollectionList(collectionsList)
	for _, collectionModel := range collectionsModel.List {
		collectionModel.WithAccessLevel(level2collection[*collectionModel.ID])
	}

	return users.NewGetUserAboutOK().WithPayload(&models.AUser{
		UserID: &targetUserID,
		Medias: models.NewModelMediaList(medias),
		Users:  models.NewModelUserList(rootsUsersList),
		Roots:  rootModelList,

		Collections:  collectionsModel,
		Rights:       models.NewModelShortUserEntityRightList(rightsList),
		TaskPreviews: models.NewModelTaskPreviewList(tasksList),
	})
}
