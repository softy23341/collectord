package resource

import (
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/auth"
	"github.com/go-openapi/runtime/middleware"
	"golang.org/x/crypto/bcrypt"
)

// RecoveryPassword TBD
func (a *Auth) ResetPassword(params auth.PostAuthPasswordResetTokenParams) middleware.Responder {
	logger := a.Context.Log
	logger.Debug("reset password")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return auth.NewPostAuthPasswordResetTokenDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	forbiddenResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return auth.NewPostAuthLoginForbidden().WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	email, err := a.EmailTokenStorage.Get(params.Token)
	if err != nil || len(email) == 0 {
		return forbiddenResponse(403, "password reset token is invalid or has expired")
	}

	var user *dto.User
	users, err := a.Context.DBM.GetUsersByEmail([]string{email})
	if err != nil {
		logger.Error("cant check users present GetUsersByEmail", "err", err)
		return errorResponse(500, err.Error())
	}
	if len(users) != 1 {
		logger.Error("can't find user")
		return forbiddenResponse(400, "can't find user")
	}
	user = users[0]

	// gen password
	password := *params.RPasswordReset.Password
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("cant crypt password", "err", err)
		return errorResponse(500, err.Error())
	}
	user.EncryptedPassword = string(encryptedPassword)

	// verify email
	if !user.EmailVerified {
		user.EmailVerified = true
	}

	if err := a.Context.DBM.UpdateUser(user); err != nil {
		logger.Error("cant update user", "err", err)
		return errorResponse(500, err.Error())
	}

	go a.EmailTokenStorage.Del(params.Token)

	return auth.NewPostAuthPasswordResetTokenNoContent()
}
