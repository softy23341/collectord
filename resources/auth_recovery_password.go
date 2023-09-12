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
func (a *Auth) RecoveryPassword(params auth.PostAuthPasswordRecoveryParams) middleware.Responder {
	logger := a.Context.Log
	logger.Debug("recovery password")

	recoveryUserEmail := *params.RUserEmail.Email

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return auth.NewPostAuthPasswordRecoveryDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	var user *dto.User
	if userList, err := a.Context.DBM.GetUsersByEmail([]string{recoveryUserEmail}); err != nil {
		logger.Error("cant get user by email", "err", err)
		return errorResponse(500, err.Error())
	} else if len(userList) != 1 {
		logger.Warn(fmt.Sprintf("user (email: %s) has not been found", recoveryUserEmail))
		return auth.NewPostAuthPasswordRecoveryNoContent()
	} else {
		user = userList[0]
	}

	token, err := util.GenerateInviteToken()
	if err != nil {
		return errorResponse(500, err.Error())
	}

	if err := a.EmailTokenStorage.Set(token, user.Email); err != nil {
		return errorResponse(500, err.Error())
	}

	scheme := "https"
	resetPasswordURL := fmt.Sprintf("%v://%v/reset-password/%v",
		scheme,
		params.HTTPRequest.Host,
		token)

	a.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		subject, html := NewResetPasswordMail(a.Templates, user, resetPasswordURL)
		err := a.Context.MailClient.Send(&services.Mail{
			To:      []string{recoveryUserEmail},
			From:    services.SystemMailFrom,
			Subject: subject,
			Body:    html,
		})
		if err != nil {
			logger.Error("cant send email", "err", err)
		}
	}))

	return auth.NewPostAuthPasswordRecoveryNoContent()
}
