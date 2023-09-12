package resource

import (
	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/teams"
	"github.com/go-openapi/runtime/middleware"
)

// GetTeamByRootID TBD
func (t *Team) GetTeamByRootID(params teams.GetTeamByRootIDParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("get team by root ID")

	DBM := t.Context.DBM
	targetRootID := params.RootID

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return teams.NewGetTeamByRootIDDefault(code).WithPayload(
			&models.Error{Code: &c, Message: msg},
		)
	}

	notFoundResponse := func(msg string) middleware.Responder {
		code := int32(404)
		return teams.NewGetTeamByRootIDNotFound().WithPayload(
			&models.Error{Code: &code, Message: msg},
		)
	}

	forbiddenErrorResponse := teams.NewGetTeamByRootIDForbidden

	// check rights
	ok, err := NewAccessRightsChecker(DBM, logger.New("service", "access checker")).
		IsUserRelatedToRoots(userContext.User.ID, []int64{targetRootID})
	if err != nil {
		logger.Error("get check user rights", "err", err)
		return errorResponse(500, err.Error())
	} else if !ok {
		logger.Error("user cant access to this root")
		return forbiddenErrorResponse()
	}

	userList, err := DBM.GetUsersByRootID(targetRootID)
	if err != nil {
		logger.Error("get users by root", "err", err)
		return errorResponse(500, err.Error())
	}
	if len(userList) == 0 {
		logger.Error("cant find root users")
		return notFoundResponse("cant find root users")
	}

	// invites
	inviteList, err := DBM.GetInvitesByRoot(targetRootID, dto.InviteCreated)
	if err != nil {
		logger.Error("get invites by root", "err", err)
		return errorResponse(500, err.Error())
	}
	invitedUserIDToInvite := inviteList.InvitedUserIDToInvite()

	// users
	invitedUserList, err := DBM.GetUsersByIDs(inviteList.InvitedUsersIDs())
	if err != nil {
		logger.Error("cant get users by invites")
	}

	// task count
	taskCnt, err := DBM.GetTasksCntByUserID(userContext.User.ID, userList.IDs())
	if err != nil {
		logger.Error("cant get users by invites")
	}

	allUserList := append(userList, invitedUserList...)

	mediaIDs := allUserList.GetAvatarsMediaIDs()

	medias, err := DBM.GetMediasByIDs(mediaIDs)
	if err != nil {
		logger.Error("get medias", "err", err)
		return errorResponse(500, "can't get additionalMediasIDs")
	}

	modelInvitedUserList := models.NewModeInvitedUserList(invitedUserList)
	for _, modelInvitedUser := range modelInvitedUserList.List {
		invite := invitedUserIDToInvite[*modelInvitedUser.UserID]
		if invite != nil {
			modelInvitedUser.WithInviteID(invite.ID)
		}
	}

	return teams.NewGetTeamByRootIDOK().WithPayload(&models.ATeam{
		Tusers:       models.NewModeTUserList(userList, taskCnt),
		Users:        models.NewModelUserList(allUserList),
		InvitedUsers: modelInvitedUserList,
		Medias:       models.NewModelMediaList(medias),
	})
}
