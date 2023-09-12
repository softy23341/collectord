package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/invites"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// CreateInvite TBD
func (i *Invite) CreateInvite(params invites.PostInviteNewParams, principal interface{}) middleware.Responder {
	var (
		userContext = principal.(*auth.UserContext)
		logger      = userContext.Logger(params.HTTPRequest)

		DBM     = i.Context.DBM
		iParams = params.RCreateInvite

		errorResponse = func(code int, msg string) middleware.Responder {
			c := int32(code)
			return invites.NewPostInviteNewDefault(code).WithPayload(
				&models.Error{Code: &c, Message: msg},
			)
		}
	)

	logger.Debug("invite user")

	// get target user
	var targetUser *dto.User
	if usersList, err := DBM.GetUsersByIDs([]int64{*iParams.ToUserID}); err != nil {
		logger.Error("cant get user by id", "err", err)
		return errorResponse(500, err.Error())
	} else if len(usersList) != 1 {
		logger.Error("cant find user by id", "err", err)
		return errorResponse(500, "invalid user id")
	} else {
		targetUser = usersList[0]
	}

	if userRootRefs, err := DBM.GetUserRootRefs([]int64{*iParams.RootID}); err != nil {
		logger.Error("cant dto.NewInvite", "err", err)
		return errorResponse(500, err.Error())
	} else if len(userRootRefs) == 0 {
		err := fmt.Errorf("wrong root id: %d", *iParams.RootID)
		logger.Error("cant find root users", "err", err)
		return invites.NewPostInviteNewNotFound()
	} else {
		isInviterInRoot := false
		for _, rootRef := range userRootRefs {
			if rootRef.UserID == targetUser.ID {
				logger.Error("user already in this root", "err", err)
				return invites.NewPostInviteNewConflict()
			}
			if rootRef.UserID == userContext.User.ID {
				isInviterInRoot = true
			}
		}

		if !isInviterInRoot {
			return invites.NewPostInviteNewForbidden()
		}
	}

	invite, err := dto.NewInvite(*iParams.RootID)
	if err != nil {
		logger.Error("cant dto.NewInvite", "err", err)
		return errorResponse(500, err.Error())
	}

	invite.FromUser(userContext.User.ID).ToUser(targetUser.ID)

	oldInvite, err := DBM.GetInviteByUserRoot(userContext.User.ID, targetUser.ID, *iParams.RootID, dto.InviteCreated)
	if err != nil {
		logger.Error("cant get invite", "err", err)
		return errorResponse(500, err.Error())
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

	if oldInvite == nil {
		if err := tx.CreateInvite(invite); err != nil {
			logger.Error("cant create invite", "err", err)
			return errorResponse(500, err.Error())
		}
	} else {
		invite = oldInvite
	}

	inviteUsersIDs := invite.UsersIDs()
	userList, err := tx.GetUsersByIDs(inviteUsersIDs)
	if err != nil {
		logger.Error("GetUsersByIDs", "error", err)
		return errorResponse(500, err.Error())
	}

	mediaList, err := tx.GetMediasByIDs(userList.GetAvatarsMediaIDs())
	if err != nil {
		logger.Error("GetMediasByIDs", "error", err)
		return errorResponse(500, err.Error())
	}

	systemUser, err := tx.GetSystemUser()
	if err != nil {
		logger.Error("get system user", "err", err)
		return errorResponse(500, err.Error())
	}

	if len(targetUser.Email) != 0 {
		sendEmailJob := delayedjob.NewJob(delayedjob.Immideate, func() {
			subject, body := NewTeamInviteMail(i.Templates, userContext.User, targetUser)
			err := i.Context.MailClient.Send(&services.Mail{
				To:      []string{targetUser.Email},
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

	inviteMessage := &dto.Message{
		UserID:     systemUser.ID,
		UserUniqID: util.NextUniqID(),
		PeerID:     targetUser.ID,
		PeerType:   dto.PeerTypeUser,
		Typo:       dto.MessageTypeService,
		MessageExtra: dto.MessageExtra{
			Service: &dto.ServiceMessage{
				Type: dto.ServiceMessageTypeInvite,
				Invite: &dto.ServiceMessageInvite{
					InviteID: invite.ID,
				},
			},
		},
	}

	result, err := i.Context.MessengerClient.
		NewMessageSender(tx, inviteMessage, &services.MessageInfo{}).
		Send()

	if err != nil {
		logger.Error("cant send message", "err", err)
		return errorResponse(500, err.Error())
	}

	events = append(events, result.Events...)
	jobs = append(jobs, result.Jobs...)

	if err := tx.Commit(); err != nil {
		logger.Error("tx.commit", "err", err)
		return errorResponse(500, err.Error())
	}

	i.Context.EventSender.Send(events...)
	i.Context.JobPool.Enqueue(jobs...)

	return invites.NewPostInviteNewOK().WithPayload(&models.ACreateInvite{
		Users:  models.NewModelUserList(userList),
		Medias: models.NewModelMediaList(mediaList),
		Invite: models.NewModelInvite(invite),
	})
}
