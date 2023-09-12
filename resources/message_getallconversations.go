package resource

import (
	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/messages"
	"github.com/go-openapi/runtime/middleware"
)

// GetAllConversations TBD
func (m *Message) GetAllConversations(params messages.PostMessageAllConversationParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("GetAllConversations")

	DBM := m.Context.DBM

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return messages.NewPostMessageAllConversationDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	peers, err := m.Context.DBM.GetUserPeers(userContext.User.ID)
	if err != nil {
		logger.Error("Get peers", "err", err)
		return errorResponse(500, err.Error())
	}

	var recentMessages dto.MessageList
	var conversationUsers dto.ConversationUserList

	rangePaginator := &dto.RangePaginator{Distance: -10}
	if n := params.RGetAllConversation.NLastMessages; n != nil {
		rangePaginator.Distance = -(*n)
	}
	for _, peer := range peers {
		var (
			err    error
			cusers dto.ConversationUserList
			msgs   dto.MessageList
		)
		switch peer.Type {
		case dto.PeerTypeUser:
			if cusers, err = DBM.GetDialogConversationUsers(userContext.User.ID, peer.ID); err == nil {
				msgs, err = DBM.GetDialogMessagesByRange(userContext.User.ID, peer.ID, rangePaginator)
			}
		case dto.PeerTypeChat:
			if cusers, err = DBM.GetChatConversationUsers(peer.ID); err == nil {
				msgs, err = DBM.GetChatMessagesByRange(peer.ID, rangePaginator)
			}
		}
		if err != nil {
			logger.Error("Get mesages", "err", err)
			return errorResponse(500, err.Error())
		}
		conversationUsers = append(conversationUsers, cusers...)
		recentMessages = append(recentMessages, msgs...)
	}
	recentMessagesUsersIDs := recentMessages.UsersIDs()

	// contacts
	relatedUserRootRefs, err := m.Context.DBM.GetUserRelatedRootRefs([]int64{userContext.User.ID})
	if err != nil {
		logger.Error("GetUserRelatedRootRefs", "err", err)
		return errorResponse(500, err.Error())
	}
	rootContactUsersIDs := relatedUserRootRefs.UsersList()

	// fetch invites
	invitesIDs := recentMessages.InvitesIDs()
	inviteList, err := m.Context.DBM.GetInvitesByIDs(invitesIDs)
	if err != nil {
		logger.Error("getinvitesbyids", "err", err)
		return errorResponse(500, err.Error())
	}
	IDToInvite := inviteList.IDToInvite()

	// fetch tasks
	tasksIDs := recentMessages.GetTaskIDs()
	taskList, err := m.Context.DBM.GetTasksByIDs(tasksIDs)
	if err != nil {
		logger.Error("gettasksbyids", "err", err)
	}

	rootIDs := inviteList.RootsIDs()
	userRootRefs, err := m.Context.DBM.GetUserRootRefs(rootIDs)
	if err != nil {
		logger.Error("get user root refs", "err", err)
		return errorResponse(500, err.Error())
	}

	// fetch chats
	chatList, err := m.Context.DBM.GetChatsByIDs(peers.ChatIDs())
	if err != nil {
		logger.Error("GetChatsByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// fetch objects
	objectIDs := recentMessages.ObjectIDs()
	objectPreviewList, err := m.Context.DBM.GetObjectsPreviewByIDs(objectIDs)
	if err != nil {
		logger.Error("GetObjectsByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// get objects rights
	userCollectionRightList, err := m.Context.DBM.GetUserRightsForCollections(userContext.User.ID, objectPreviewList.GetCollectionsIDs())
	if err != nil {
		logger.Error("GetUserRightsForCollections", "err", err)
		return errorResponse(500, err.Error())
	}
	collections2level := userCollectionRightList.GetEntityIDToLevel(dto.RightEntityTypeCollection)

	objectMeta, err := NewObjectPreviewExtractor(m.Context.DBM, logger.New("object extractor"), &ObjectPreviewExtractorOpts{
		mediaExtractorOpts: &objectMediasExtractorOpts{
			onlyPhoto: true,
		},
	}).
		SetObjectsIDs(objectIDs).
		FetchAll().
		Result()

	if err != nil {
		logger.Error("NewObjectPreviewExtractor", "err", err)
		return errorResponse(500, err.Error())
	}

	// fetch users
	usersIDs := recentMessagesUsersIDs
	usersIDs = append(usersIDs, taskList.GetUsersIDs()...)
	usersIDs = append(usersIDs, inviteList.UsersIDs()...)
	usersIDs = append(usersIDs, userRootRefs.GetRootOwnerIDs()...)
	usersIDs = append(usersIDs, rootContactUsersIDs...)

	userList, err := m.Context.DBM.GetUsersByIDs(usersIDs)
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500, err.Error())
	}
	userListMapByID := userList.IDToUser()

	// fetch medias
	mediasIDs := userList.GetAvatarsMediaIDs()
	mediasIDs = append(mediasIDs, recentMessages.MediasIDs()...)
	mediasIDs = append(mediasIDs, objectMeta.MediaExtractor.getMediasIDs()...)
	mediasIDs = append(mediasIDs, chatList.MediasIDs()...)

	var mediaList dto.MediaList
	if len(mediasIDs) > 0 {
		mediaList, err = m.Context.DBM.GetMediasByIDs(mediasIDs)
		if err != nil {
			logger.Error("GetMediasByIDs", "err", err)
			return errorResponse(500, err.Error())
		}
	}
	mediaListMapById := mediaList.IDToMedia()

	objectStatusMap := objectMeta.ObjectStatuses.ObjectToOneStatusMap()

	objectValuationMap := objectMeta.Valuations.ObjectToOneValuationMap()

	// construct object json models
	modelsObjectPreviews := models.NewModelObjectPreviewList(objectPreviewList)
	for _, modelObjectPreview := range modelsObjectPreviews.List {
		ID := *modelObjectPreview.ID
		originalLocationsIDs := objectMeta.OriginLocationExtractor.getOriginLocationsIDsByObjectID(ID)
		modelObjectPreview.
			WithMediaIDs(objectMeta.MediaExtractor.getMediasIDsByObjectID(ID)).
			WithActorIDs(objectMeta.ActorExtractor.getActorsIDsByObjectID(ID)).
			WithOriginLocationIDs(originalLocationsIDs).
			WithBadgeIDs(objectMeta.BadgeExtractor.getBadgesIDsByObjectID(ID)).
			WithStatus(objectStatusMap[ID]).
			WithValuations(models.NewModelValuationList(objectValuationMap[ID])).
			WithAccessLevel(collections2level[*modelObjectPreview.CollectionID])
	}

	// locations
	originLocations := objectMeta.OriginLocationExtractor.getOriginLocations()

	// chats
	modelsChats := models.NewModelChatList(chatList)
	for _, modelChat := range modelsChats.List {
		// yeah so slow
		modelChat.WithUsersCnt(int16(len(conversationUsers.GetUsersByPeer(modelChat.ID, dto.PeerTypeChat))))
	}

	// messages
	modelsMessages := models.NewModelMessageList(recentMessages)
	for _, modelMessage := range modelsMessages.List {
		if srvMsg := modelMessage.ServiceMessage; srvMsg != nil {
			if invMsg := srvMsg.ServiceMessageInvite; invMsg != nil {
				invite := IDToInvite[invMsg.InviteID]
				if invite != nil {
					invMsg.FillInvite(invite).
						SetRootOwner(userRootRefs.GetRootOwnerID(invite.RootID))
				}
			}
			if invMsg := srvMsg.ServiceMessageInviteStatusChanged; invMsg != nil {
				var iu *dto.User
				var ium *dto.Media
				invite := IDToInvite[invMsg.InviteID]
				if invite != nil && invite.ToUserID != nil {
					if user, found := userListMapByID[*invite.ToUserID]; found {
						iu = user
					}
				}
				if iu != nil && iu.AvatarMediaID != nil {
					if media, found := mediaListMapById[*iu.AvatarMediaID]; found {
						ium = media
					}
				}
				invMsg.FillInviteStatusChanged(iu, ium)
			}
		}
	}

	return messages.NewPostMessageAllConversationOK().WithPayload(&models.AGetAllConversation{
		TaskPreviews:              models.NewModelTaskPreviewList(taskList),
		ContactUsersIds:           append(recentMessagesUsersIDs, rootContactUsersIDs...),
		ConversationUsers:         models.NewModelConversationUserList(conversationUsers),
		Chats:                     modelsChats,
		Messages:                  modelsMessages,
		Medias:                    models.NewModelMediaList(mediaList),
		Users:                     models.NewModelUserList(userList),
		ObjectsPreview:            modelsObjectPreviews,
		Actors:                    models.NewModelActorList(objectMeta.ActorExtractor.getActors()),
		Badges:                    models.NewModelBadgeList(objectMeta.BadgeExtractor.getBadges()),
		OriginLocations:           models.NewModelOriginLocationList(originLocations),
		NTotalUnreadMessages:      userContext.User.NUnreadMessages,
		NTotalUnreadNotifications: userContext.User.NUnreadNotifications,
	})
}
