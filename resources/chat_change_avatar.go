package resource

import (
	"errors"
	"fmt"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/chat"
	"git.softndit.com/collector/backend/services"
	"github.com/go-openapi/runtime/middleware"
)

// ChangeAvatar TBD
func (c *Chat) ChangeAvatar(params chat.PostChatIDChangeAvatarParams, principal interface{}) middleware.Responder {
	p := params.REditChatAvatar

	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)

	DBM := c.Context.DBM

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return chat.NewPostChatIDChangeAvatarDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	// get media
	var avatarMedia *dto.Media
	if p.AvatarID != nil {
		if mediaList, err := DBM.GetMediasByIDs([]int64{*p.AvatarID}); err != nil {
			logger.Error("GetMediasByIDs", "err", err)
			return errorResponse(500, err.Error())
		} else if len(mediaList) != 1 {
			err := fmt.Errorf("Avatar not found: %d", *p.AvatarID)
			logger.Error("GetMediasByIDs", "err", err)
			return errorResponse(422, err.Error())
		} else {
			// check access rights
			ok := NewAccessRightsChecker(c.Context.DBM, logger.New("service", "access rights")).
				IsUserOwnerOfMedias(userContext.User.ID, mediaList)
			if err != nil {
				logger.Error("cant check rights", "err", err)
				return errorResponse(500, err.Error())
			} else if !ok {
				err := errors.New("user is not owner of media")
				logger.Error("user cant allow be here", "err", err)
				return chat.NewPostChatIDChangeAvatarForbidden()
			}

			avatarMedia = mediaList[0]
		}
	}

	// check for repeated request
	if existMsg, err := DBM.GetMessageByUserUniqID(userContext.User.ID, *p.ClientUniqID); err != nil {
		logger.Error("GetMessageByUserUniqID", "err", err)
		return errorResponse(500, err.Error())
	} else if existMsg != nil {
		logger.Warn("return old msg", "id", existMsg.ID)
		return chat.NewPostChatIDChangeAvatarOK().WithPayload(&models.AEditChatAvatar{
			NewAvatarServiceMessage: models.NewModelMessage(existMsg),
			Avatar:                  models.NewModelMedia(avatarMedia),
		})
	}

	// get chat
	peer := dto.Peer{ID: params.ID, Type: dto.PeerTypeChat}
	var changedChat *dto.Chat
	if chats, err := DBM.GetChatsByIDs([]int64{peer.ID}); err != nil {
		logger.Error("GetChatsByIDs", "err", err)
		return errorResponse(500, err.Error())
	} else if len(chats) < 1 {
		err := fmt.Errorf("cant find chat with id %d", peer.ID)
		logger.Error("GetChatsByIDs", "err", err)
		return errorResponse(422, err.Error())
	} else {
		changedChat = chats[0]
	}

	// get users in chat
	conversationUsers, err := DBM.GetChatConversationUsers(peer.ID)
	if err != nil {
		logger.Error("GetChatConversationUsers", "err", err)
		return errorResponse(500, err.Error())
	}

	// check user present
	var curConversationUser *dto.ConversationUser
	for _, cu := range conversationUsers {
		if cu.UserID == userContext.User.ID {
			curConversationUser = cu
			break
		}
	}

	if curConversationUser == nil {
		err = errors.New("cant change chat avatar: user must be a chat member to change avatar")
		logger.Error("change avatar", "err", err)
		return chat.NewPostChatIDChangeAvatarForbidden()
	}

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin tx", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	changedChat.AvatarMediaID = p.AvatarID

	if err := tx.UpdateChat(changedChat); err != nil {
		logger.Error("UpdateChat", "err", err)
		return errorResponse(500, err.Error())
	}

	//  service message to inform other users about chat title change
	serviceMsg := &dto.Message{
		UserID:     userContext.User.ID,
		UserUniqID: *p.ClientUniqID,
		PeerID:     peer.ID,
		PeerType:   dto.PeerTypeChat,
		Typo:       dto.MessageTypeService,
		MessageExtra: dto.MessageExtra{
			Service: &dto.ServiceMessage{
				Type: dto.ServiceMessageTypeChatAvatarChanged,
				ChatAvatarChanged: &dto.ServiceMessageChatAvatarChanged{
					NewAvatarID: p.AvatarID,
				},
			},
		},
	}

	result, err := c.Context.MessengerClient.
		NewMessageSender(tx, serviceMsg, &services.MessageInfo{
			ConversationUsers: conversationUsers, Chat: changedChat,
		}).Send()
	if err != nil {
		logger.Error("send msg", "err", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("tx.commit", "err", err)
		return errorResponse(500, err.Error())
	}

	c.Context.EventSender.Send(result.Events...)

	// add jobs
	c.Context.JobPool.Enqueue(result.Jobs...)

	return chat.NewPostChatIDChangeAvatarOK().WithPayload(&models.AEditChatAvatar{
		NewAvatarServiceMessage: models.NewModelMessage(serviceMsg),
		Avatar:                  models.NewModelMedia(avatarMedia),
	})
}
