package resource

import (
	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/invites"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// AcceptInvite TBD
func (i *Invite) AcceptInvite(params invites.PostInviteIDAcceptParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("accept invite")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return invites.NewPostInviteIDAcceptDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	DBM := i.Context.DBM

	var invite *dto.Invite
	if inviteList, err := DBM.GetInvitesByIDs([]int64{params.ID}); err != nil {
		logger.Error("cant get invites by ids", "err", err)
		return errorResponse(500, err.Error())
	} else if len(inviteList) != 1 {
		logger.Error("cant get invites by ids", "err", "cant find invite")
		return invites.NewPostInviteIDAcceptNotFound()
	} else {
		invite = inviteList[0]
	}

	if invite.IsAccepted() {
		logger.Warn("invite already was accepted")
		return invites.NewPostInviteIDAcceptNoContent()
	}

	if !invite.IsCreated() {
		return invites.NewPostInviteIDAcceptConflict()
	}

	canAccept := NewAccessRightsChecker(DBM, logger.New("service", "access checker")).
		WasUserInvited(userContext.User.ID, invite)
	if !canAccept {
		logger.Error("cant accept invite")
		return invites.NewPostInviteIDAcceptForbidden()
	}

	// send message
	var (
		events dto.EventList
		jobs   []delayedjob.Job
	)

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("dbm.begintx fail", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	if err := tx.ChangeInviteStatus(params.ID, dto.InviteAccepted); err != nil {
		logger.Error("cant change invite status", "err", err)
		return errorResponse(500, err.Error())
	}

	systemUser, err := tx.GetSystemUser()
	if err != nil {
		logger.Error("get system user", "err", err)
		return errorResponse(500, err.Error())
	}

	// get creator user
	var creatorUser *dto.User
	if usersList, err := DBM.GetUsersByIDs([]int64{invite.CreatorUserID}); err != nil {
		logger.Error("cant get user by id", "err", err)
		return errorResponse(500, err.Error())
	} else if len(usersList) != 1 {
		logger.Error("cant find user by id", "err", err)
		return errorResponse(500, "invalid user id")
	} else {
		creatorUser = usersList[0]
	}

	if len(creatorUser.Email) != 0 {
		sendEmailJob := delayedjob.NewJob(delayedjob.Immideate, func() {
			subject, body := InviteWasAcceptedMail(i.Templates, userContext.User, creatorUser)
			err := i.Context.MailClient.Send(&services.Mail{
				To:      []string{creatorUser.Email},
				From:    services.SystemMailFrom,
				Subject: subject,
				Body:    body,
			})
			if err != nil {
				logger.Error("cant send email", "err", err)
			}
		})
		jobs = append(jobs, sendEmailJob)
	}

	inviteAcceptedMessage := &dto.Message{
		UserID:     systemUser.ID,
		UserUniqID: util.NextUniqID(),
		PeerID:     invite.CreatorUserID,
		PeerType:   dto.PeerTypeUser,
		Typo:       dto.MessageTypeService,
		MessageExtra: dto.MessageExtra{
			Service: &dto.ServiceMessage{
				Type: dto.ServiceMessageTypeInviteStatusChanged,
				InviteStatusChanged: &dto.ServiceMessageInviteStatusChanged{
					InviteID:     invite.ID,
					InviteStatus: dto.InviteAccepted,
				},
			},
		},
	}

	rootRef := dto.NewRegularUserRootRef(*invite.ToUserID, invite.RootID)
	if err := tx.CreateUserRootRef(rootRef); err != nil {
		logger.Error("CreateUserRootRef", "err", err)
		return errorResponse(500, err.Error())
	}

	result, err := i.Context.MessengerClient.
		NewMessageSender(tx, inviteAcceptedMessage, &services.MessageInfo{}).
		Send()

	if err != nil {
		logger.Error("cant send message", "err", err)
		return errorResponse(500, err.Error())
	}

	events = append(events, result.Events...)
	jobs = append(jobs, result.Jobs...)

	if err := tx.Commit(); err != nil {
		logger.Error("commit", "error", err)
		return errorResponse(500, err.Error())
	}

	i.Context.EventSender.Send(events...)
	i.Context.JobPool.Enqueue(jobs...)

	return invites.NewPostInviteIDAcceptNoContent()
}
