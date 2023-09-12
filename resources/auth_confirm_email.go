package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/auth"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// RecoveryPassword TBD
func (a *Auth) ConfirmEmail(params auth.GetAuthRegConfirmEmailParams) middleware.Responder {
	logger := a.Context.Log
	logger.Debug("confirm email")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return auth.NewGetAuthRegConfirmEmailDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	email, err := a.EmailTokenStorage.Get(params.Token)
	if err != nil || len(email) == 0 {
		return errorResponse(400, "token is invalid or has expired")
	}

	var user *dto.User
	usersList, err := a.Context.DBM.GetUsersByEmail([]string{email})
	if err != nil {
		logger.Error("cant check users present GetUsersByEmail", "err", err)
		return errorResponse(500, err.Error())
	}
	if len(usersList) != 1 {
		logger.Error("can't find user")
		return errorResponse(400, "can't find user")
	}
	user = usersList[0]

	if user.EmailVerified == true {
		return errorResponse(400, "email has already been verified")
	}

	user.EmailVerified = true
	if err := a.Context.DBM.UpdateUser(user); err != nil {
		logger.Error("cant update user", "err", err)
		return errorResponse(500, err.Error())
	}

	go a.EmailTokenStorage.Del(params.Token)
	return auth.NewGetAuthRegConfirmEmailNoContent()
}

func (a *Auth) SendConfirmEmailToken(params auth.PostAuthRegConfirmEmailParams) middleware.Responder {
	logger := a.Context.Log
	logger.Debug("send confirm email")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return auth.NewGetAuthRegConfirmEmailDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	// check email
	regParams := params.RUserEmail
	users, err := a.Context.DBM.GetUsersByEmail([]string{*regParams.Email})
	if err != nil {
		logger.Error("cant check users present GetUsersByEmail", "err", err)
		return errorResponse(500, err.Error())
	}
	if len(users) != 1 {
		logger.Error("can't find user")
		return errorResponse(403, "user not found")
	}
	user := users[0]

	if user.EmailVerified == true {
		return errorResponse(400, "email has already been confirmed")
	}

	//send confirm_email
	confirmEmailToken, err := util.GenerateInviteToken()
	if err != nil {
		logger.Error("generate token", "err", err.Error())
		return errorResponse(500, err.Error())
	}
	if err := a.EmailTokenStorage.Set(confirmEmailToken, user.Email); err != nil {
		logger.Error("redis: set key", "err", err.Error())
		return errorResponse(500, err.Error())
	}

	scheme := "https"
	confirmEmailURL := fmt.Sprintf("%v://%v/confirm-email/%v",
		scheme,
		params.HTTPRequest.Host,
		confirmEmailToken)

	a.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		subject, html := NewConfirmEmailMail(a.Templates, user, confirmEmailURL)
		err := a.Context.MailClient.Send(&services.Mail{
			To:      []string{user.Email},
			From:    services.SystemMailFrom,
			Subject: subject,
			Body:    html,
		})
		if err != nil {
			logger.Error("cant send email", "err", err)
		}
	}))

	return auth.NewPostAuthRegConfirmEmailNoContent()
}
