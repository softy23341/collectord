package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/roots"
	"github.com/go-openapi/runtime/middleware"
)

// Root TBD
type Root struct {
	Context Context
}

// RootInfo TBD
func (r *Root) RootInfo(params roots.GetRootIDParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("Root info")

	var (
		errorResponse = func(code int) middleware.Responder {
			return roots.NewGetRootIDDefault(code)
		}
		accessChecker = NewAccessRightsChecker(
			r.Context.DBM,
			logger.New("service", "access rights"),
		)
	)

	var mediasIDs []int64

	if ok, err := accessChecker.IsUserRelatedToRoots(userContext.User.ID, []int64{params.ID}); err != nil {
		logger.Error("cant access to root", "err", err)
		return errorResponse(500)
	} else if !ok {
		err := fmt.Errorf("user is not related to root: %d", params.ID)
		logger.Error("user not allowed be here", "err", err)
		return roots.NewGetRootIDForbidden()
	}

	rootsList, err := r.Context.DBM.GetRootsByIDs([]int64{params.ID})
	if err != nil {
		logger.Error("GetRootsByIDs", "err", err)
		return errorResponse(500)
	}
	if len(rootsList) != 1 {
		logger.Error("GetRootsByIDs", "err", "404")
		return errorResponse(404)
	}
	root := rootsList[0]

	users, err := r.Context.DBM.GetUsersByRootID(root.ID)
	if err != nil {
		logger.Error("GetUsersByRootID", "err", err)
		return errorResponse(500)
	}
	mediasIDs = append(mediasIDs, users.GetAvatarsMediaIDs()...)

	var medias dto.MediaList

	if len(mediasIDs) > 0 {
		medias, err = r.Context.DBM.GetMediasByIDs(mediasIDs)
		if err != nil {
			logger.Error("GetMediasByIDs", "err", err)
			return errorResponse(500)
		}
	}

	return roots.NewGetRootIDOK().WithPayload(&models.ARoot{
		Medias: models.NewModelMediaList(medias),
		Root:   models.NewModelRoot(root),
		Users:  models.NewModelUserList(users),
	})
}

// UserRoots TBD
func (r *Root) UserRoots(params roots.GetRootByUserParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("user roots")

	errorResponse := func(code int) middleware.Responder {
		return roots.NewGetRootByUserDefault(code)
	}

	// depreciated
	return roots.NewGetRootByUserForbidden()

	var mediaIDs []int64

	targetUserID := userContext.User.ID
	if params.UserID != nil {
		targetUserID = *params.UserID
	}
	rootsList, err := r.Context.DBM.GetRootsByUserID(targetUserID)
	if err != nil {
		logger.Error("GetUserRoots", "err", err)
		return errorResponse(500)
	}

	userRootRefs, err := r.Context.DBM.GetUserRootRefs(rootsList.IDs())
	if err != nil {
		logger.Error("GetUserRoots", "err", err)
		return errorResponse(500)
	}

	userListIDs := userRootRefs.UsersList()
	rootToUsers := userRootRefs.RootToUsers()

	rootsUsersList, err := r.Context.DBM.GetUsersByIDs(userListIDs)
	if err != nil {
		logger.Error("rootsUsersList", "err", err)
		return errorResponse(500)
	}

	mediaIDs = append(mediaIDs, rootsUsersList.GetAvatarsMediaIDs()...)

	var medias dto.MediaList
	if len(mediaIDs) > 0 {
		var err error
		medias, err = r.Context.DBM.GetMediasByIDs(mediaIDs)
		if err != nil {
			logger.Error("GetMediasByIDs(mediaIDs)", "err", err)
			return errorResponse(500)
		}
	}

	rootModelList := models.NewModelRootList(rootsList)
	for _, rootModel := range rootModelList.List {
		rootUsers := rootToUsers[*rootModel.ID]
		if rootUsers != nil {
			rootModel.
				WithOwnerID(rootUsers.OwnerID).
				WithUsersIDs(rootUsers.UsersIDs)
		}
	}

	return roots.NewGetRootByUserOK().WithPayload(&models.ARoots{
		Users:  models.NewModelUserList(rootsUsersList),
		Roots:  rootModelList,
		Medias: models.NewModelMediaList(medias),
	})
}

// AddUserToRoot TBD
func (r *Root) AddUserToRoot(params roots.PostRootAddUserParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("add user to root")

	errorResponse := func(code int) middleware.Responder {
		return roots.NewPostRootAddUserDefault(code)
	}

	aParams := params.RAddUserToRoot
	targetRootID, targetUserID := *aParams.RootID, *aParams.UserID

	// depreciated method
	return roots.NewPostRootAddUserForbidden()

	_, err := r.Context.DBM.GetUsersByIDs([]int64{targetUserID})
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500)
	}

	_, err = r.Context.DBM.GetRootsByIDs([]int64{targetRootID})
	if err != nil {
		logger.Error("GetRootsByIDs", "err", err)
		return errorResponse(500)
	}

	rootRef := dto.NewRegularUserRootRef(targetUserID, targetRootID)
	if err := r.Context.DBM.CreateUserRootRef(rootRef); err != nil {
		logger.Error("CreateUserRootRef", "err", err)
		return errorResponse(500)
	}

	return roots.NewPostRootAddUserNoContent()
}

// RemoveUserFromRoot TBD
func (r *Root) RemoveUserFromRoot(params roots.PostRootRemoveUserParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("remove user from root")

	var (
		errorResponse = func(code int) middleware.Responder {
			return roots.NewPostRootRemoveUserDefault(code)
		}

		notFoundResponse = func(code int) middleware.Responder {
			return roots.NewPostRootRemoveUserNotFound()
		}

		rmParams                   = params.RRemoveUserFromRoot
		targetRootID, targetUserID = *rmParams.RootID, *rmParams.UserID
		accessChecker              = NewAccessRightsChecker(
			r.Context.DBM,
			logger.New("service", "access rights"),
		)
	)

	// check root owner
	if ok, err := accessChecker.IsUserRootOwner(userContext.User.ID, targetRootID); err != nil {
		logger.Error("cant check user owners", "err", err)
		return errorResponse(500)
	} else if !ok {
		err := fmt.Errorf("user is not the owner of the root (%d)", targetRootID)
		logger.Error("access forbidden", "err", err)
		return roots.NewPostRootRemoveUserForbidden()
	}

	usersList, err := r.Context.DBM.GetUsersByIDs([]int64{targetUserID})
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500)
	}
	if len(usersList) != 1 {
		logger.Error("cant find user", "userID", targetUserID)
		return notFoundResponse(404)
	}

	// get target root
	rootsList, err := r.Context.DBM.GetRootsByIDs([]int64{targetRootID})
	if err != nil {
		logger.Error("GetRootsByIDs", "err", err)
		return errorResponse(500)
	}
	if len(rootsList) != 1 {
		logger.Error("cant find root", "rootID", targetRootID)
		return notFoundResponse(404)
	}

	// get roots refs
	rootsRefList, err := r.Context.DBM.GetUserRelatedRootRefs([]int64{targetUserID})
	if err != nil {
		logger.Error("GetUserRelatedRootRefs", "err", err)
		return errorResponse(500)
	}

	rootFound := false
	for _, rootRef := range rootsRefList {
		if rootRef.RootID == targetRootID && rootRef.UserID == targetUserID {
			if rootRef.IsOwner() {
				logger.Error("cant delete user from own root")
				return roots.NewPostRootRemoveUserForbidden()
			}
			rootFound = true
			break // ok. it's just ref
		}
	}
	if !rootFound {
		logger.Error("cant find user root", "rootID", targetRootID, "userID", targetUserID)
		return errorResponse(500)
	}

	err = r.Context.DBM.DeleteUserRootRef(targetUserID, targetRootID)
	if err != nil {
		logger.Error("DeleteUserRootRef", "err", err)
		return errorResponse(500)
	}

	return roots.NewPostRootRemoveUserNoContent()
}
