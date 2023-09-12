package resource

import (
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/chat"
	"github.com/go-openapi/runtime/middleware"
)

// GetChat TBD
func (c *Chat) GetChat(params chat.GetChatIDParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return chat.NewGetChatIDDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	notFoundResponse := func() middleware.Responder {
		return chat.NewGetChatIDNotFound()
	}

	// get chat
	var nChat *dto.Chat
	if chats, err := c.Context.DBM.GetChatsByIDs([]int64{params.ID}); err != nil {
		logger.Error("GetChatsByIDs", "err", err)
		return errorResponse(500, err.Error())
	} else if len(chats) < 1 {
		err := fmt.Errorf("cant find chat with id %d", params.ID)
		logger.Error("GetChatsByIDs", "err", err)
		return notFoundResponse()
	} else {
		nChat = chats[0]
	}

	// get users
	cUsers, err := c.Context.DBM.GetChatConversationUsers(nChat.ID)
	if err != nil {
		logger.Error("GetChatConversationUsers", "err", err)
		return errorResponse(500, err.Error())
	}

	foundCurrentUser := false
	for _, conversationUser := range cUsers {
		if conversationUser.UserID == userContext.User.ID {
			foundCurrentUser = true
			break
		}
	}
	if !foundCurrentUser {
		logger.Error("you cant access to this chat")
		return chat.NewGetChatIDForbidden()
	}

	users, err := c.Context.DBM.GetUsersByIDs(cUsers.GetUsersIDs())
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// get medias
	mediasIDs := users.GetAvatarsMediaIDs()
	if nChat.AvatarMediaID != nil {
		mediasIDs = append(mediasIDs, *nChat.AvatarMediaID)
	}

	medias, err := c.Context.DBM.GetMediasByIDs(mediasIDs)
	if err != nil {
		logger.Error("GetMediasByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	return chat.NewGetChatIDOK().WithPayload(&models.AGetChatInfo{
		Chat:   models.NewModelChat(nChat).WithUsersCnt(int16(len(cUsers))),
		Medias: models.NewModelMediaList(medias),
		Users:  models.NewModelUserList(users),
	})

}
