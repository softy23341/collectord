package resource

import (
	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/restapi/operations/users"
	"github.com/go-openapi/runtime/middleware"
)

// UpdateUser TBD
func (u *User) UpdateUser(params users.PutUserParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("update user")

	errorResponse := func(code int) middleware.Responder {
		return users.NewPutUserDefault(code)
	}

	targerUserID := params.REditUser.ID
	if targerUserID == nil {
		targerUserID = &userContext.User.ID
	}

	if *targerUserID != userContext.User.ID {
		logger.Error("access error")
		return users.NewPutUserForbidden()
	}

	// TODO any user editing
	params.REditUser.User.MergeToDtoUser(userContext.User)

	if err := u.Context.DBM.UpdateUser(userContext.User); err != nil {
		logger.Error("DBM.UpdateUser", "err", err)
		return errorResponse(500)
	}

	return users.NewPutUserNoContent()
}

func (u *User) UpdateUserLocale(params users.PutUserUpdateLocaleParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("update user locale")

	errorResponse := func(code int) middleware.Responder {
		return users.NewPutUserDefault(code)
	}

	newLocale := params.RUpdateUserLocale.Locale
	if newLocale == nil {
		return errorResponse(500)
	}

	if userContext.User.Locale == *newLocale {
		return users.NewPutUserNoContent()
	}

	userContext.User.Locale = *newLocale
	if err := u.Context.DBM.UpdateUser(userContext.User); err != nil {
		logger.Error("DBM.UpdateUserLocale", "err", err)
		return errorResponse(500)
	}

	return users.NewPutUserNoContent()
}
