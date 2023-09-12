package resource

import (
	"fmt"
	"time"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/invites"
	"git.softndit.com/collector/backend/services"
	"github.com/go-openapi/runtime/middleware"
)

// Send TBD
func (i *InviteByEmail) Send(params invites.PostInviteByEmailParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)

		DBM     = i.Context.DBM
		iParams = params.RInviteByEmail

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return invites.NewPostInviteNewDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
	)

	logger.Debug("invite by email user")

	// check user
	if usersList, _ := DBM.GetUsersByEmail([]string{*iParams.ToUserEmail}); len(usersList) != 0 {
		return invites.NewPostInviteByEmailConflict()
	}

	invite, err := dto.NewInvite(*iParams.RootID)
	if err != nil {
		logger.Error("cant dto.NewInvite", "err", err)
		return errorResponse(500, err.Error())
	}

	invite.FromUser(userContext.User.ID)
	invite.ToUserEmail = iParams.ToUserEmail

	if err != nil {
		logger.Error("generate token", "err", err.Error())
		return errorResponse(500, err.Error())
	}

	var oldInvite *dto.Invite
	invitesList, err := DBM.GetInvitesByRoot(*iParams.RootID, dto.InviteCreated)
	if err != nil {
		logger.Error("cant get invite", "err", err)
		return errorResponse(500, err.Error())
	}
	for i, in := range invitesList {
		if in.ToUserEmail != nil && *in.ToUserEmail == *iParams.ToUserEmail {
			oldInvite = invitesList[i]
		}
	}

	if oldInvite == nil {
		if err := DBM.CreateInvite(invite); err != nil {
			logger.Error("cant create invite", "err", err)
			return errorResponse(500, err.Error())
		}
	} else {
		invite = oldInvite
		// do not let users spam invites
		if time.Now().Sub(oldInvite.CreationTime) < 3*time.Hour {
			return invites.NewPostInviteByEmailNoContent()
		}
	}

	//send email
	scheme := "https"
	confirmURL := fmt.Sprintf("%v://%v/register?email=%s&invite=%s",
		scheme,
		params.HTTPRequest.Host,
		*iParams.ToUserEmail,
		invite.Token)

	i.Context.JobPool.Enqueue(delayedjob.NewJob(delayedjob.Immideate, func() {
		subject, html := NewInviteByEmailMail(i.Templates, userContext.User, &dto.User{Email: *iParams.ToUserEmail}, confirmURL)
		err := i.Context.MailClient.Send(&services.Mail{
			To:      []string{*iParams.ToUserEmail},
			From:    services.SystemMailFrom,
			Subject: subject,
			Body:    html,
		})
		if err != nil {
			logger.Error("cant send email", "err", err)
		}
	}))

	return invites.NewPostInviteByEmailNoContent()
}
