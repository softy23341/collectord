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

// Chat TBD
type Chat struct {
	Context Context
}

// CreateChat TBD
func (c *Chat) CreateChat(params chat.PostChatParams, principal interface{}) middleware.Responder {
	p := params.RCreateChat

	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return chat.NewPostChatDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	createdChat, err := c.Context.DBM.GetChatByUserUniqID(userContext.User.ID, *p.ClientUniqID)
	if err != nil {
		logger.Error("GetChatByUserUniqID", "err", err)
		return errorResponse(500, err.Error())
	}
	if createdChat != nil {
		firstMessage, err := c.Context.DBM.GetMessageByUserUniqID(userContext.User.ID, *p.ClientUniqID)
		if err != nil {
			logger.Error("GetMessageByUserUniqID", "err", err)
			return errorResponse(500, err.Error())
		}
		if firstMessage == nil {
			logger.Error("chat without fist msg!", "err", err)
			return errorResponse(500, err.Error())
		}

		return chat.NewPostChatOK().WithPayload(&models.ACreateChat{
			ID:                       &createdChat.ID,
			CreateChatServiceMessage: models.NewModelMessage(firstMessage),
		})
	}

	// check users
	usersList, err := c.Context.DBM.GetUsersByIDs(p.UsersIds)
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500, err.Error())
	}
	if len(usersList) != len(p.UsersIds) {
		err := fmt.Errorf("users not found or contain duplicates %v", p.UsersIds)
		logger.Error("usersList", "err", err)
		return errorResponse(422, err.Error())
	}

	// check media
	if avatarID := p.AvatarMediaID; avatarID != nil {
		mediaList, err := c.Context.DBM.GetMediasByIDs([]int64{*avatarID})
		if err != nil {
			logger.Error("GetMediasByIDs", "err", err)
			return errorResponse(500, err.Error())
		}

		if len(mediaList) == 0 {
			err := fmt.Errorf("cant find image with id %d", *avatarID)
			logger.Error("fetch avatar by id", "err", err)
			return errorResponse(422, err.Error())
		}

		// check access rights
		ok := NewAccessRightsChecker(c.Context.DBM, logger.New("service", "access rights")).
			IsUserOwnerOfMedias(userContext.User.ID, mediaList)
		if err != nil {
			logger.Error("cant check rights", "err", err)
			return errorResponse(500, err.Error())
		} else if !ok {
			err := errors.New("user is not owner of media")
			logger.Error("user cant allow be here", "err", err)
			return chat.NewPostChatForbidden()
		}
	}

	// start transaction
	tx, err := c.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("begintx", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// create chat
	newChat := &dto.Chat{
		AdminUserID:   userContext.User.ID,
		CreatorUserID: userContext.User.ID,
		UserUniqID:    *p.ClientUniqID,
		Name:          *p.Name,
		AvatarMediaID: p.AvatarMediaID,
	}

	if err := tx.CreateChat(newChat); err != nil {
		logger.Error("CreateChat", "err", err)
		return errorResponse(500, err.Error())
	}

	// create conversation user for each user
	var conversationUsers = make(dto.ConversationUserList, 0, len(p.UsersIds)+1)

	for _, participantID := range p.UsersIds {
		conversationUser := &dto.ConversationUser{
			UserID:        participantID,
			PeerID:        newChat.ID,
			PeerType:      dto.PeerTypeChat,
			JoinedAt:      *newChat.CreationTime,
			InvitedUserID: &userContext.User.ID,
		}
		if err := tx.CreateConversationUser(conversationUser); err != nil {
			logger.Error("CreateConversationUser", "err", err)
			return errorResponse(500, err.Error())
		}

		conversationUsers = append(conversationUsers, conversationUser)
	}

	// create self conversation user and ref
	conversationUser := &dto.ConversationUser{
		UserID:          userContext.User.ID,
		PeerID:          newChat.ID,
		PeerType:        dto.PeerTypeChat,
		JoinedAt:        *newChat.CreationTime,
		NUnreadMessages: 0,
	}
	if err := tx.CreateConversationUser(conversationUser); err != nil {
		logger.Error("CreateConversationUser", "err", err)
		return errorResponse(500, err.Error())
	}

	conversationUsers = append(conversationUsers, conversationUser)

	// service message to inform other users about chat creation
	chatCreateMessage := &dto.Message{
		UserID:     userContext.User.ID,
		UserUniqID: *p.ClientUniqID,
		PeerID:     newChat.ID,
		PeerType:   dto.PeerTypeChat,
		Typo:       dto.MessageTypeService,
		MessageExtra: dto.MessageExtra{
			Service: &dto.ServiceMessage{
				Type: dto.ServiceMessageTypeChatCreated,
				ChatCreated: &dto.ServiceMessageChatCreated{
					ChatTitle: newChat.Name,
				},
			},
		},
	}

	result, err := c.Context.MessengerClient.
		NewMessageSender(tx, chatCreateMessage, &services.MessageInfo{
			ConversationUsers: conversationUsers,
			Chat:              newChat}).
		Send()

	if err != nil {
		logger.Error("send service message fail", "err", err)
		return errorResponse(500, err.Error())
	}

	// commit changes
	if err := tx.Commit(); err != nil {
		logger.Error("commit", "err", err)
		return errorResponse(500, err.Error())
	}

	c.Context.EventSender.Send(result.Events...)

	// add jobs
	c.Context.JobPool.Enqueue(result.Jobs...)

	return chat.NewPostChatOK().WithPayload(&models.ACreateChat{
		ID:                       &newChat.ID,
		CreateChatServiceMessage: models.NewModelMessage(chatCreateMessage),
	})
}

// AddUser TBD
func (c *Chat) AddUser(params chat.PostChatAddUserParams, principal interface{}) middleware.Responder {
	p := params.RChatAddUser

	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)

	DBM := c.Context.DBM

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return chat.NewPostChatAddUserDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	// check for repeated request
	if existMsg, err := DBM.GetMessageByUserUniqID(userContext.User.ID, *p.ClientUniqID); err != nil {
		logger.Error("GetMessageByUserUniqID", "err", err)
		return errorResponse(500, err.Error())
	} else if existMsg != nil {
		logger.Warn("return old msg", "id", existMsg.ID)
		return chat.NewPostChatAddUserOK().WithPayload(&models.AChatAddUser{
			ChatNewUserMessage: models.NewModelMessage(existMsg),
		})
	}

	//
	peer := dto.Peer{ID: *p.ChatID, Type: dto.PeerTypeChat}

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin tx", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// check users
	var chatForInvite *dto.Chat
	{
		chats, err := tx.GetChatsByIDs([]int64{peer.ID})
		if err != nil {
			logger.Error("GetChatsByIDs", "err", err)
			return errorResponse(500, err.Error())
		}
		if len(chats) < 1 {
			err = fmt.Errorf("Cant find chat with id %d", peer.ID)
			logger.Error("GetChatsByIDs", "err", err)
			return errorResponse(422, err.Error())
		}
		chatForInvite = chats[0]
	}

	conversationUsers, err := tx.GetChatConversationUsers(peer.ID)
	if err != nil {
		logger.Error("GetChatConversationUsers", "err", err)
		return errorResponse(500, err.Error())
	}

	var isCurrentUserChatMemeber bool
	for _, cu := range conversationUsers {
		switch cu.UserID {
		case userContext.User.ID:
			isCurrentUserChatMemeber = true
		case *p.UserID:
			logger.Debug("User already exists")
			return errorResponse(500, "user already present in this chat")
		}
	}
	if !isCurrentUserChatMemeber {
		err := fmt.Errorf("User violent chat rules: %d", userContext.User.ID)
		logger.Debug("user must be in chat for invite", "err", err)
		return chat.NewPostChatAddUserForbidden()
	}

	// TODO check chat and users
	newConversationUser := &dto.ConversationUser{
		UserID:        *p.UserID,
		PeerID:        peer.ID,
		PeerType:      peer.Type,
		InvitedUserID: &userContext.User.ID,
	}

	if err := tx.CreateConversationUser(newConversationUser); err != nil {
		logger.Error("CreateConversationUser", "err", err)
		return errorResponse(500, err.Error())
	}

	conversationUsers = append(conversationUsers, newConversationUser)

	//  service message to inform other users about new user
	addUserMessage := &dto.Message{
		UserID:     userContext.User.ID,
		UserUniqID: *p.ClientUniqID,
		PeerID:     chatForInvite.ID,
		PeerType:   dto.PeerTypeChat,
		Typo:       dto.MessageTypeService,
		MessageExtra: dto.MessageExtra{
			Service: &dto.ServiceMessage{
				Type: dto.ServiceMessageTypeChatUserAdded,
				ChatUserAdded: &dto.ServiceMessageChatUserAdded{
					AddedUserID: *p.UserID,
				},
			},
		},
	}

	result, err := c.Context.MessengerClient.
		NewMessageSender(tx, addUserMessage, &services.MessageInfo{
			ConversationUsers: conversationUsers, Chat: chatForInvite,
		}).Send()
	if err != nil {
		logger.Error("cant send message", "err", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.Commit(); err != nil {
		logger.Error("tx.commit", "err", err)
		return errorResponse(500, err.Error())
	}

	// add event
	c.Context.EventSender.Send(result.Events...)

	// add jobs
	c.Context.JobPool.Enqueue(result.Jobs...)

	return chat.NewPostChatAddUserOK().WithPayload(&models.AChatAddUser{
		ChatNewUserMessage: models.NewModelMessage(addUserMessage),
	})
}

// RemoveUser TBD
func (c *Chat) RemoveUser(params chat.PostChatRemoveUserParams, principal interface{}) middleware.Responder {
	p := params.RChatRemoveUser

	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	DBM := c.Context.DBM

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return chat.NewPostChatRemoveUserDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	// check for repeated request
	if existMsg, err := DBM.GetMessageByUserUniqID(userContext.User.ID, *p.ClientUniqID); err != nil {
		logger.Error("GetMessageByUserUniqID", "err", err)
		return errorResponse(500, err.Error())
	} else if existMsg != nil {
		logger.Warn("return old msg", "id", existMsg.ID)
		return chat.NewPostChatAddUserOK().WithPayload(&models.AChatAddUser{
			ChatNewUserMessage: models.NewModelMessage(existMsg),
		})
	}

	peer := dto.Peer{ID: *p.ChatID, Type: dto.PeerTypeChat}

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin tx", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// check chat
	var chatForRemove *dto.Chat
	{
		chats, err := tx.GetChatsByIDs([]int64{peer.ID})
		if err != nil {
			logger.Error("GetChatsByIDs", "err", err)
			return errorResponse(500, err.Error())
		}
		if len(chats) < 1 {
			err = fmt.Errorf("Cant find chat with id %d", peer.ID)
			logger.Error("GetChatsByIDs", "err", err)
			return errorResponse(422, err.Error())
		}
		chatForRemove = chats[0]
	}

	// check users
	conversationUsers, err := tx.GetChatConversationUsers(peer.ID)
	if err != nil {
		logger.Error("GetChatConversationUsers", "err", err)
		return errorResponse(500, err.Error())
	}

	var (
		curConversationUser *dto.ConversationUser
		delConversationUser *dto.ConversationUser
	)
	for _, cu := range conversationUsers {
		if cu.UserID == userContext.User.ID {
			curConversationUser = cu
		}
		// no "else if" here, coz current user can delete themself (c)
		if cu.UserID == *p.UserID {
			delConversationUser = cu
		}
		if curConversationUser != nil && delConversationUser != nil {
			break
		}
	}
	if curConversationUser == nil {
		err = errors.New("cant del user from chat: user must be a chat member to delete another users")
		logger.Error("remove user", "err", err)
		return chat.NewPostChatRemoveUserForbidden()
	}
	if delConversationUser == nil {
		err = errors.New("cant del user from chat: deleting user not fond")
		logger.Error("remove user", "err", err)
		return errorResponse(500, err.Error())
	}

	isAllowed := (*p.UserID == userContext.User.ID) ||
		(chatForRemove.AdminUserID == userContext.User.ID) ||
		(delConversationUser.InvitedUserID != nil && *delConversationUser.InvitedUserID == userContext.User.ID)
	if !isAllowed {
		err = errors.New("cant del user from chat: have no rights")
		logger.Error("remove user", "err", err)
		return chat.NewPostChatRemoveUserForbidden()
	}

	//  service message to inform other users about delete
	delUserMessage := &dto.Message{
		UserID:     userContext.User.ID,
		UserUniqID: *p.ClientUniqID,
		PeerID:     chatForRemove.ID,
		PeerType:   dto.PeerTypeChat,
		Typo:       dto.MessageTypeService,
		MessageExtra: dto.MessageExtra{
			Service: &dto.ServiceMessage{
				Type: dto.ServiceMessageTypeChatUserDeleted,
				ChatUserDeleted: &dto.ServiceMessageChatUserDeleted{
					DeletedUserID: *p.UserID,
				},
			},
		},
	}

	result, err := c.Context.MessengerClient.
		NewMessageSender(tx, delUserMessage, &services.MessageInfo{
			ConversationUsers: conversationUsers, Chat: chatForRemove,
		}).Send()
	if err != nil {
		logger.Error("cant send message", "err", err)
		return errorResponse(500, err.Error())
	}

	if err := tx.DelConversationUser(*p.UserID, peer); err != nil {
		logger.Error("cant DelConversationUser", "err", err)
		return errorResponse(500, err.Error())
	}

	chatForRemovePeer := &dto.Peer{
		Type: dto.PeerTypeChat,
		ID:   chatForRemove.ID,
	}

	nread, err := tx.DelAllUnreadMessagesByPeer(*p.UserID, chatForRemovePeer)
	if err != nil {
		logger.Error("DelAllUnreadMessagesByPeer", "err", err)
		return errorResponse(500, err.Error())
	}

	if n, err := tx.UpdateUserNUnreadMessages(*p.UserID, -nread); err != nil {
		logger.Error("UpdateUserNUnreadMessages", "err", err)
		return errorResponse(500, err.Error())
	} else if n == nil {
		logger.Error("UpdateUserNUnreadMessages", "err", "user not found")
		return errorResponse(422, fmt.Sprintf("cuser not dound %+v", *p.UserID))
	}

	if n, err := tx.UpdateUserNUnreadNotifications(*p.UserID, -nread); err != nil {
		logger.Error("UpdateUserNUnreadNotifications", "err", err)
		return errorResponse(500, err.Error())
	} else if n == nil {
		logger.Error("UpdateUserNUnreadNotifications", "err", "user not found")
		return errorResponse(422, fmt.Sprintf("cuser not dound %+v", *p.UserID))
	}

	if err := tx.Commit(); err != nil {
		logger.Error("tx.commit", "err", err)
		return errorResponse(500, err.Error())
	}

	// add event
	c.Context.EventSender.Send(result.Events...)

	// add jobs
	c.Context.JobPool.Enqueue(result.Jobs...)

	return chat.NewPostChatRemoveUserOK().WithPayload(&models.AChatRemoveUser{
		ChatRemoveUserMessage: models.NewModelMessage(delUserMessage),
	})
}

// ChangeName TBD
func (c *Chat) ChangeName(params chat.PostChatIDChangeNameParams, principal interface{}) middleware.Responder {
	p := params.REditChatName

	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	DBM := c.Context.DBM

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return chat.NewPostChatIDChangeNameDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	// check for repeated request
	if existMsg, err := DBM.GetMessageByUserUniqID(userContext.User.ID, *p.ClientUniqID); err != nil {
		logger.Error("GetMessageByUserUniqID", "err", err)
		return errorResponse(500, err.Error())
	} else if existMsg != nil {
		logger.Warn("return old msg", "id", existMsg.ID)
		return chat.NewPostChatIDChangeNameOK().WithPayload(&models.AEditChatName{
			NewNameServiceMessage: models.NewModelMessage(existMsg),
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
		err = errors.New("cant change chat name: user must be a chat member to change name")
		logger.Error("change name", "err", err)
		return chat.NewPostChatIDChangeNameForbidden()
	}

	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("begin tx", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	changedChat.Name = *p.Name

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
				Type: dto.ServiceMessageTypeChatTitleChanged,
				ChatTitleChanged: &dto.ServiceMessageChatTitleChanged{
					NewTitle: *p.Name,
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

	return chat.NewPostChatIDChangeNameOK().WithPayload(&models.AEditChatName{
		NewNameServiceMessage: models.NewModelMessage(serviceMsg),
	})
}
