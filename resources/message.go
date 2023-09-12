package resource

import (
	"errors"
	"fmt"
	"html/template"

	"git.softndit.com/collector/backend/auth"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/restapi/operations/messages"
	"git.softndit.com/collector/backend/services"
	"git.softndit.com/collector/backend/util"
	"github.com/go-openapi/runtime/middleware"
)

// Message TBD
type Message struct {
	Context   Context
	Templates *template.Template
}

// NewDirectMessageMail TBD
func NewDirectMessageMail(templates *template.Template, from, to *dto.User, message *dto.Message) (string, string) {
	T := util.GetTranslationFunc(to.Locale)

	subject := T("api.share.subject")

	var title string
	if message.IsText() && message.Text != nil {
		title = T("api.share.message",
			map[string]interface{}{"SenderName": from.FullName(), "Message": message.Text.Text})
	} else if message.IsMedia() {
		title = T("api.share.media",
			map[string]interface{}{"SenderName": from.FullName()})
	} else if message.IsObject() {
		title = T("api.share.object",
			map[string]interface{}{"SenderName": from.FullName()})
	}

	props := map[interface{}]interface{}{
		"Props": map[interface{}]interface{}{
			"Title":  title,
			"Footer": T("api.footer.info"),
		},
	}

	return subject, util.Parse(templates, "directmessage_body", props)
}

// SendMessage TBD
func (m *Message) SendMessage(params messages.PostMessageParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("send message")

	errorResponse := func(code int) middleware.Responder {
		return messages.NewPostMessageDefault(code)
	}

	successResponse := func(messageID int64) middleware.Responder {
		return messages.NewPostMessageOK().WithPayload(&models.ASendMessage{
			MessageID: &messageID,
		})
	}

	DBM := m.Context.DBM
	accessChecker := NewAccessRightsChecker(m.Context.DBM, logger.New("service", "access rights"))
	sParams := params.RSendMessage

	// send message
	var (
		events dto.EventList
		jobs   []delayedjob.Job
	)

	// check by user uniq id
	msg, err := DBM.GetMessageByUserUniqID(userContext.User.ID, *sParams.ClientUniqID)
	if err != nil {
		logger.Error("GetMessageByUserUniqID", "err", err)
		return errorResponse(500)
	}
	if msg != nil {
		logger.Warn("msg already present")
		return successResponse(msg.ID)
	}

	// check users exists
	var peerUser *dto.User
	peer := models.NewDtoInputPeer(sParams.Peer)
	if peer.Type.IsUser() {
		if userList, err := DBM.GetUsersByIDs([]int64{peer.ID}); err != nil {
			logger.Debug("GetUsersByIDs", "err", err)
			return errorResponse(500)
		} else if len(userList) != 1 {
			logger.Error("cant find user", "id", peer.ID)
			return errorResponse(500)
		} else {
			peerUser = userList[0]
		}
		if isUserBanned, _ := DBM.IsUserBanned(peer.ID, userContext.User.ID); isUserBanned {
			return messages.NewPostMessageLocked()
		}
	} else if peer.Type.IsChat() {
		if chatList, err := DBM.GetChatsByIDs([]int64{peer.ID}); err != nil {
			logger.Debug("GetChatsByIDs", "err", err)
			return errorResponse(500)
		} else if len(chatList) != 1 {
			logger.Error("cant find chat", "id", peer.ID)
			return errorResponse(500)
		} else {
			ok, err := accessChecker.IsUserInChat(userContext.User.ID, peer.ID)
			if err != nil {
				logger.Error("cant check users", "err", err)
				return errorResponse(500)
			}

			if !ok {
				logger.Error("user is not in chat")
				return messages.NewPostMessageForbidden()
			}
		}
	}

	messageType, err := dto.NewMessageType(*sParams.Typo)

	// create message
	message := &dto.Message{
		UserID:     userContext.User.ID,
		UserUniqID: *sParams.ClientUniqID,
		PeerType:   peer.Type,
		PeerID:     peer.ID,
		Typo:       messageType,
	}

	if message.IsText() {
		message.Text = &dto.TextMessage{
			Text: *sParams.TextMessage.Text,
		}
	} else if message.IsMedia() {
		mediaID := *sParams.MediaMessage.MediaID
		if mediaList, err := DBM.GetMediasByIDs([]int64{mediaID}); err != nil {
			logger.Error("GetMediasByIDs", "err", err)
			return errorResponse(500)
		} else if len(mediaList) != 1 {
			err := fmt.Errorf("cant find media: %d", mediaID)
			logger.Error("GetMediasByIDs", "err", err)
			return errorResponse(500)
		} else {
			// check access rights
			ok := accessChecker.IsUserOwnerOfMedias(userContext.User.ID, mediaList)
			if err != nil {
				logger.Error("cant check rights", "err", err)
				return errorResponse(500)
			} else if !ok {
				err := errors.New("user is not owner of media")
				logger.Error("user cant allow be here", "err", err)
				return messages.NewPostMessageForbidden()
			}
		}

		message.Media = &dto.MediaMessage{
			MediaID: mediaID,
		}
	} else if message.IsObject() {
		objectID := *sParams.CobjectMessage.ObjectID
		// get object collection
		var collection *dto.Collection
		if collections, err := DBM.GetCollectionsByObjectIDs([]int64{objectID}); err != nil {
			err := fmt.Errorf("cant find collection for: %d", objectID)
			logger.Error("GetCollectionsByObjectIDs", "err", err)
			return errorResponse(500)
		} else if len(collections) != 1 {
			err := fmt.Errorf("cant find collection for: %d", objectID)
			logger.Error("GetCollectionsByObjectIDs", "err", err)
			return errorResponse(500)
		} else {
			collection = collections[0]
		}
		// check is collection public
		if !collection.Public {
			if objectList, err := DBM.GetObjectsByIDs([]int64{objectID}); err != nil {
				logger.Error("GetObjectsByIDs", "err", err)
				return errorResponse(500)
			} else if len(objectList) != 1 {
				err := fmt.Errorf("cant find object: %d", objectID)
				logger.Error("GetObjectsByIDs", "err", err)
				return errorResponse(500)
			} else {
				ok, err := accessChecker.HasUserRightsForObjects(userContext.
					User.ID,
					dto.RightEntityLevelRead,
					objectList.GetIDs())
				if err != nil {
					logger.Error("cant check rights", "err", err)
					return errorResponse(500)
				} else if !ok {
					err := errors.New("user is not allowed to touch object")
					logger.Error("user cant allow be here", "err", err)
					return messages.NewPostMessageForbidden()
				}
			}
		}
		message.Object = &dto.ObjectMessage{
			ObjectID: objectID,
		}
	}

	// send message
	tx, err := DBM.BeginTx()
	if err != nil {
		logger.Error("dbm.begintx fail", "err", err)
		return errorResponse(500)
	}
	defer tx.Rollback()

	result, err := m.Context.MessengerClient.
		NewMessageSender(tx, message, &services.MessageInfo{}).
		Send()

	if err != nil {
		logger.Error("cant send message", "err", err)
		return errorResponse(500)
	}

	// try to send message by email
	if peerUser != nil && len(peerUser.Email) != 0 {
		sendEmailJob := delayedjob.NewJob(delayedjob.Immideate, func() {
			subject, body := NewDirectMessageMail(m.Templates, userContext.User, peerUser, message)
			err := m.Context.MailClient.Send(&services.Mail{
				To:      []string{peerUser.Email},
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

	events = append(events, result.Events...)
	jobs = append(jobs, result.Jobs...)

	if err := tx.Commit(); err != nil {
		logger.Error("tx.commit", "err", err)
		return errorResponse(500)
	}

	// add event
	m.Context.EventSender.Send(events...)

	// add jobs
	m.Context.JobPool.Enqueue(jobs...)

	return successResponse(message.ID)
}

// GetMessages TBD
func (m *Message) GetMessages(params messages.PostMessageByidsParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("get messages")

	var (
		errorResponse = func(code int) middleware.Responder {
			return messages.NewPostMessageByidsDefault(code)
		}
		accessChecker = NewAccessRightsChecker(m.Context.DBM, logger.New("service", "access rights"))
		gParams       = params.RGetMessages
	)

	messageList, err := m.Context.DBM.GetMessagesByIDs(gParams.Ids)
	if err != nil {
		logger.Error("GetMessagesByIDs", "err", err)
		return errorResponse(500)
	}

	chatsToCheck := make(map[int64]struct{}, len(messageList))
	for _, message := range messageList {
		peer := message.Peer()
		if peer.IsUser() {
			if message.UserID == userContext.User.ID || peer.ID == userContext.User.ID {
				continue
			}
			err := fmt.Errorf("cant access to msg: %d", message.ID)
			logger.Error("user cant access to this message", "err", err)
			return messages.NewPostMessageByidsForbidden()
		}

		if peer.IsChat() {
			if message.Service != nil && message.Service.ChatUserDeleted != nil {
				continue
			}
			chatsToCheck[peer.ID] = struct{}{}
		}
	}

	for chatID := range chatsToCheck {
		ok, err := accessChecker.IsUserInChat(userContext.User.ID, chatID)
		if err != nil {
			logger.Error("cant check user in chat", "err", err)
			return errorResponse(500)
		} else if !ok {
			err := fmt.Errorf("cant access to chat: %d", chatID)
			logger.Error("user cant access to this chat", "err", err)
			return messages.NewPostMessageByidsForbidden()
		}

	}

	// fetch invites
	invitesIDs := messageList.InvitesIDs()
	inviteList, err := m.Context.DBM.GetInvitesByIDs(invitesIDs)
	if err != nil {
		logger.Error("getinvitesbyids", "err", err)
	}
	IDToInvite := inviteList.IDToInvite()

	rootIDs := inviteList.RootsIDs()
	userRootRefs, err := m.Context.DBM.GetUserRootRefs(rootIDs)
	if err != nil {
		logger.Error("get user root refs", "err", err)
		return errorResponse(500)
	}

	// fetch tasks
	tasksIDs := messageList.GetTaskIDs()
	taskList, err := m.Context.DBM.GetTasksByIDs(tasksIDs)
	if err != nil {
		logger.Error("gettasksbyids", "err", err)
	}

	// users
	usersIDs := messageList.UsersIDs()
	usersIDs = append(usersIDs, taskList.GetUsersIDs()...)
	usersIDs = append(usersIDs, inviteList.UsersIDs()...)
	usersIDs = append(usersIDs, userRootRefs.GetRootOwnerIDs()...)
	userList, err := m.Context.DBM.GetUsersByIDs(usersIDs)
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500)
	}
	userListMapByID := userList.IDToUser()

	// fetch chats
	chatList, err := m.Context.DBM.GetChatsByIDs(messageList.ChatsIDs())
	if err != nil {
		logger.Error("GetChatsByIDs", "err", err)
		return errorResponse(500)
	}

	// fetch chat cnt
	chats2cnt, err := m.Context.DBM.GetChatMembersCntByChats(chatList.IDs())
	if err != nil {
		logger.Error("GetChatMembersCntByChats", "err", err)
		return errorResponse(500)
	}

	// fetch objects
	objectIDs := messageList.ObjectIDs()
	objectPreviewList, err := m.Context.DBM.GetObjectsPreviewByIDs(objectIDs)
	if err != nil {
		logger.Error("GetObjectsByIDs", "err", err)
		return errorResponse(500)
	}

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
		return errorResponse(500)
	}

	// get objects rights
	userCollectionRightList, err := m.Context.DBM.GetUserRightsForCollections(userContext.User.ID, objectPreviewList.GetCollectionsIDs())
	if err != nil {
		logger.Error("GetUserRightsForCollections", "err", err)
		return errorResponse(500)
	}
	collections2level := userCollectionRightList.GetEntityIDToLevel(dto.RightEntityTypeCollection)

	mediasIDs := userList.GetAvatarsMediaIDs()
	mediasIDs = append(mediasIDs, messageList.MediasIDs()...)
	mediasIDs = append(mediasIDs, objectMeta.MediaExtractor.getMediasIDs()...)
	mediasIDs = append(mediasIDs, chatList.MediasIDs()...)

	var mediaList dto.MediaList
	if len(mediasIDs) > 0 {
		mediaList, err = m.Context.DBM.GetMediasByIDs(mediasIDs)
		if err != nil {
			logger.Error("GetMediasByIDs", "err", err)
			return errorResponse(500)
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

	// origin locations
	originLocations := objectMeta.OriginLocationExtractor.getOriginLocations()

	// messages
	modelsMessages := models.NewModelMessageList(messageList)
	for _, modelMessage := range modelsMessages.List {
		if srvMsg := modelMessage.ServiceMessage; srvMsg != nil {
			if invMsg := srvMsg.ServiceMessageInvite; invMsg != nil {
				invite := IDToInvite[invMsg.InviteID]
				invMsg.
					FillInvite(invite).
					SetRootOwner(userRootRefs.GetRootOwnerID(invite.RootID))
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

	// fill chats with cnt
	chatModels := models.NewModelChatList(chatList)
	for _, chatModel := range chatModels.List {
		chatModel.WithUsersCnt(int16(chats2cnt[chatModel.ID]))
	}

	return messages.NewPostMessageByidsOK().WithPayload(&models.AGetMessages{
		TaskPreviews:    models.NewModelTaskPreviewList(taskList),
		Messages:        modelsMessages,
		Medias:          models.NewModelMediaList(mediaList),
		Users:           models.NewModelUserList(userList),
		ObjectsPreview:  modelsObjectPreviews,
		Actors:          models.NewModelActorList(objectMeta.ActorExtractor.getActors()),
		Badges:          models.NewModelBadgeList(objectMeta.BadgeExtractor.getBadges()),
		OriginLocations: models.NewModelOriginLocationList(originLocations),
		Chats:           chatModels,
	})
}

// GetMessagesRange TBD
func (m *Message) GetMessagesRange(params messages.PostMessageRangeParams, principal interface{}) middleware.Responder {
	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("get messages range")

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return messages.NewPostMessageRangeDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	rangePaginator := models.NewDtoRangePaginator(params.RGetMessagesRange.Range)
	peer := models.NewDtoInputPeer(params.RGetMessagesRange.Peer)
	accessChecker := NewAccessRightsChecker(m.Context.DBM, logger.New("service", "access rights"))

	var (
		messageList dto.MessageList
		err         error
		usersIDs    []int64
		chatsIDs    []int64
	)
	if peer.Type.IsUser() {
		messageList, err = m.Context.DBM.GetDialogMessagesByRange(userContext.User.ID, peer.ID, rangePaginator)
		usersIDs = append(usersIDs, peer.ID)
	} else if peer.Type.IsChat() {
		ok, err := accessChecker.IsUserInChat(userContext.User.ID, peer.ID)
		if err != nil {
			logger.Error("cant check users", "err", err)
			return errorResponse(500, err.Error())
		}

		if !ok {
			logger.Error("user is not in chat")
			return messages.NewPostMessageRangeForbidden()
		}

		messageList, err = m.Context.DBM.GetChatMessagesByRange(peer.ID, rangePaginator)
		chatsIDs = append(chatsIDs, peer.ID)
	} else {
		logger.Error("cant recognize type")
		return errorResponse(500, "cant recognize type")
	}

	if err != nil {
		logger.Error("GetMessagesByRange", "err", err)
		return errorResponse(500, err.Error())
	}

	// fetch invites
	invitesIDs := messageList.InvitesIDs()
	inviteList, err := m.Context.DBM.GetInvitesByIDs(invitesIDs)
	if err != nil {
		logger.Error("getinvitesbyids", "err", err)
	}
	IDToInvite := inviteList.IDToInvite()

	// fetch tasks
	tasksIDs := messageList.GetTaskIDs()
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

	usersIDs = append(usersIDs, messageList.UsersIDs()...)
	usersIDs = append(usersIDs, inviteList.UsersIDs()...)
	usersIDs = append(usersIDs, userRootRefs.GetRootOwnerIDs()...)
	userList, err := m.Context.DBM.GetUsersByIDs(usersIDs)
	if err != nil {
		logger.Error("GetUsersByIDs", "err", err)
		return errorResponse(500, err.Error())
	}
	userListMapByID := userList.IDToUser()

	// fetch chats
	chatsIDs = append(chatsIDs, messageList.ChatsIDs()...)
	chatList, err := m.Context.DBM.GetChatsByIDs(chatsIDs)
	if err != nil {
		logger.Error("GetChatsByIDs", "err", err)
		return errorResponse(500, err.Error())
	}

	// fetch chat cnt
	chats2cnt, err := m.Context.DBM.GetChatMembersCntByChats(chatList.IDs())
	if err != nil {
		logger.Error("GetChatMembersCntByChats", "err", err)
		return errorResponse(500, err.Error())
	}

	// fetch objects
	objectIDs := messageList.ObjectIDs()
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

	mediasIDs := userList.GetAvatarsMediaIDs()
	mediasIDs = append(mediasIDs, taskList.GetUsersIDs()...)
	mediasIDs = append(mediasIDs, messageList.MediasIDs()...)
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

	// origin locations
	originLocations := objectMeta.OriginLocationExtractor.getOriginLocations()

	// messages
	modelsMessages := models.NewModelMessageList(messageList)
	for _, modelMessage := range modelsMessages.List {
		if srvMsg := modelMessage.ServiceMessage; srvMsg != nil {
			if invMsg := srvMsg.ServiceMessageInvite; invMsg != nil {
				if invite := IDToInvite[invMsg.InviteID]; invite != nil {
					invMsg.
						FillInvite(invite).
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

	// fill chats with cnt
	chatModels := models.NewModelChatList(chatList)
	for _, chatModel := range chatModels.List {
		chatModel.WithUsersCnt(int16(chats2cnt[chatModel.ID]))
	}

	return messages.NewPostMessageRangeOK().WithPayload(&models.AGetMessagesRange{
		TaskPreviews:    models.NewModelTaskPreviewList(taskList),
		Messages:        modelsMessages,
		Medias:          models.NewModelMediaList(mediaList),
		Users:           models.NewModelUserList(userList),
		ObjectsPreview:  modelsObjectPreviews,
		Actors:          models.NewModelActorList(objectMeta.ActorExtractor.getActors()),
		Badges:          models.NewModelBadgeList(objectMeta.BadgeExtractor.getBadges()),
		OriginLocations: models.NewModelOriginLocationList(originLocations),
		Chats:           chatModels,
	})
}

// ReadHistory TBD
func (m *Message) ReadHistory(params messages.PostMessageReadHistoryParams, principal interface{}) middleware.Responder {
	p := params.RReadHistory

	userContext := principal.(*auth.UserContext)
	logger := userContext.Logger(params.HTTPRequest)
	logger.Debug("read history",
		"last msg id", *p.LastReadMessageID,
		"peer", fmt.Sprintf("%+v", *p.Peer),
	)

	errorResponse := func(code int, msg string) middleware.Responder {
		c := int32(code)
		return messages.NewPostMessageReadHistoryDefault(code).
			WithPayload(&models.Error{Code: &c, Message: msg})
	}

	peer := models.NewDtoInputPeer(p.Peer)

	tx, err := m.Context.DBM.BeginTx()
	if err != nil {
		logger.Error("dbm.begintx fail", "err", err)
		return errorResponse(500, err.Error())
	}
	defer tx.Rollback()

	// update counters
	var peerUnreadMessages int32
	var totalUnreadMessages int32
	var totalUnreadNotifications int32

	// remove messages
	nread, lastID, err := tx.DelUnreadMessages(userContext.User.ID, peer, *p.LastReadMessageID)
	if err != nil {
		logger.Error("DelUnreadMessages", "err", err)
		return errorResponse(500, err.Error())
	}

	// update conversation user counters
	if n, err := tx.UpdateConversationUserNUnreadMessages(userContext.User.ID, peer, -nread); err != nil {
		logger.Error("UpdateConversationUserNUnreadMessages", "err", err)
		return errorResponse(500, err.Error())
	} else if n == nil {
		peerUnreadMessages = 0
	}

	if n, err := tx.UpdateUserNUnreadMessages(userContext.User.ID, -nread); err != nil {
		logger.Error("UpdateUserNUnreadMessages", "err", err)
		return errorResponse(500, err.Error())
	} else if n == nil {
		logger.Error("UpdateUserNUnreadMessages", "err", "user not found")
		return errorResponse(422, fmt.Sprintf("cuser not dound %+v", *peer))
	} else {
		totalUnreadMessages = *n
	}

	if n, err := tx.UpdateUserNUnreadNotifications(userContext.User.ID, -nread); err != nil {
		logger.Error("UpdateUserNUnreadNotifications", "err", err)
		return errorResponse(500, err.Error())
	} else if n == nil {
		logger.Error("UpdateUserNUnreadNotifications", "err", "user not found")
		return errorResponse(422, fmt.Sprintf("cuser not dound %+v", *peer))
	} else {
		totalUnreadNotifications = *n
	}

	// send event
	var eventList dto.EventList
	if nread > 0 && lastID > 0 {
		if n, err := tx.UpdateConversationUserLastReadMessageID(userContext.User.ID, peer, lastID); err != nil {
			logger.Error("UpdateConversationUserLastReadMessageID", "err", err)
			return errorResponse(500, err.Error())
		} else if n == nil {
		}

		var affectedUsers []int64

		// > implying only dialogs
		if peer.Type == dto.PeerTypeUser {
			affectedUsers = []int64{peer.ID, userContext.User.ID}
		} else if peer.Type == dto.PeerTypeChat {
			chatLastReadMessageID, err := tx.GetChatLastReadMessageID(peer.ID)
			if err != nil {
				logger.Error("GetChatLastReadMessageID", "err", err)
				return errorResponse(500, err.Error())
			}
			if chatLastReadMessageID == nil {
				err := errors.New("chat not found")
				logger.Error("chatLastReadMessageID", "err", err)
				return errorResponse(500, err.Error())
			}
			if lastID > *chatLastReadMessageID {
				if _, err := tx.SetChatLastReadMessageIDIfGreater(peer.ID, lastID); err != nil {
					logger.Error("SetChatLastReadMessageIDIfGreater", "err", err)
					return errorResponse(500, err.Error())
				}
				ids, err := tx.GetChatMessagesAuthors(peer.ID, *chatLastReadMessageID+1, lastID)
				if err != nil {
					logger.Error("GetChatMessagesAuthors", "err", err)
					return errorResponse(500, err.Error())
				}
				set := util.Int64SliceToSet(ids)
				affectedUsers = set.ToSlice()
			}
		}

		selfPeer := dto.Peer{ID: userContext.User.ID, Type: dto.PeerTypeUser}
		for _, affectedUserID := range affectedUsers {
			var eventPeer dto.Peer
			if peer.Type == dto.PeerTypeUser {
				eventPeer = selfPeer
			} else if peer.Type == dto.PeerTypeChat {
				eventPeer = *peer
			}

			event := &dto.Event{
				UserID: affectedUserID,
				Type:   dto.EventTypeHistoryRead,
				EventUnion: dto.EventUnion{
					HistoryRead: &dto.EventHistoryRead{
						Peer:              eventPeer,
						LastReadMessageID: lastID,
					},
				},
			}

			if _, err := services.EmplaceEvent(tx, event); err != nil {
				logger.Error("EmplaceEvent", "err", err)
				return errorResponse(500, err.Error())
			}

			eventList = append(eventList, event)
		}
	}

	if err := tx.Commit(); err != nil {
		logger.Error("tx.commit", "err", err)
		return errorResponse(500, err.Error())
	}

	// add event
	m.Context.EventSender.Send(eventList...)

	return messages.NewPostMessageReadHistoryOK().WithPayload(&models.AReadHistory{
		NPeerUnreadMessages:       peerUnreadMessages,
		NTotalUnreadMessages:      totalUnreadMessages,
		NTotalUnreadNotifications: totalUnreadNotifications,
	})
}
