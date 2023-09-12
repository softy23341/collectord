package resource

import (
	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/users_ban_list"
	"github.com/go-openapi/runtime/middleware"
)

// UserBan TBD
type UsersBan struct {
	Context Context
}

func (u *UsersBan) GetUsersBanList(params users_ban_list.GetUsersBanListParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return users_ban_list.NewGetUsersBanListDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	userBanList, err := u.Context.DBM.GetUserBanList(userContext.User.ID)
	if err != nil {
		logger.Error("GetUserBanList", "err", err)
		return errorResponse(500, err.Error())
	}

	// get users
	users, err := u.Context.DBM.GetUsersByIDs(userBanList.UserIDs())
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// get medias
	mediasIDs := users.GetAvatarsMediaIDs()

	medias, err := u.Context.DBM.GetMediasByIDs(mediasIDs)
	if err != nil {
		logger.Error("GetMediasByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	return users_ban_list.NewGetUsersBanListOK().WithPayload(&models.AGetBanListInfo{
		Users:  models.NewModelUserList(users),
		Medias: models.NewModelMediaList(medias),
	})
}

func (u *UsersBan) UsersBanListAdd(params users_ban_list.PostUsersBanListAddParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("new ban user")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return users_ban_list.NewPostUsersBanListAddDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	isUserBanned, err := u.Context.DBM.IsUserBanned(userContext.User.ID, params.RBanAddUser.UserID)
	if err != nil {
		logger.Error("IsUserBanned", "err", err)
		return errorResponse(500, err.Error())
	}
	if isUserBanned {
		return users_ban_list.NewPostUsersBanListAddNoContent()
	}

	newUserBan := &dto.UserBan{
		CreatorUserID: userContext.User.ID,
		UserID:        params.RBanAddUser.UserID,
	}

	if err := u.Context.DBM.CreateUserBan(newUserBan); err != nil {
		logger.Error("CreateUserBan", "err", err)
		return errorResponse(500, err.Error())
	}

	return users_ban_list.NewPostUsersBanListAddNoContent()
}

func (u *UsersBan) UsersBanListRemove(params users_ban_list.PostUsersBanListRemoveParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("remove ban user")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return users_ban_list.NewPostUsersBanListAddDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	isUserBanned, err := u.Context.DBM.IsUserBanned(userContext.User.ID, params.RBanRemoveUser.UserID)
	if err != nil {
		logger.Error("IsUserBanned", "err", err)
		return errorResponse(500, err.Error())
	}
	if !isUserBanned {
		return users_ban_list.NewPostUsersBanListRemoveNotFound()
	}

	if err := u.Context.DBM.DeleteUserBan(userContext.User.ID, params.RBanRemoveUser.UserID); err != nil {
		logger.Error("CreateUserBan", "err", err)
		return errorResponse(500, err.Error())
	}

	return users_ban_list.NewPostUsersBanListRemoveNoContent()
}
