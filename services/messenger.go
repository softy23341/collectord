package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"git.softndit.com/collector/backend/dal"
	"git.softndit.com/collector/backend/delayedjob"
	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/models"
	"git.softndit.com/collector/backend/npusher"
	npusherclient "git.softndit.com/collector/backend/npusher/client"
	"git.softndit.com/collector/backend/util"
	"gopkg.in/inconshreveable/log15.v2"
)

// MessengerClient TBD
type MessengerClient struct {
	DBM        dal.TrManager
	Log        log15.Logger
	PushClient npusherclient.Client
}

// NewMessageSender TBD
func (m *MessengerClient) NewMessageSender(dbm dal.TxManager, message *dto.Message, info *MessageInfo) *MessageSender {
	return &MessageSender{
		dbm:              dbm,
		message:          message,
		info:             info,
		logger:           m.Log,
		countersByUserID: make(map[int64]MessageCounters),
		pushClient:       m.PushClient,
	}
}

// MessageSender TBD
type MessageSender struct {
	dbm              dal.TxManager
	message          *dto.Message
	info             *MessageInfo
	logger           log15.Logger
	countersByUserID map[int64]MessageCounters
	events           dto.EventList
	jobs             []delayedjob.Job
	pushClient       npusherclient.Client
}

// MessageCounters TBD
type MessageCounters struct {
	NUnreadMessages             int32
	NUnreadNotifications        int32
	NConversationUnreadMessages int32
}

// MessageInfo TBD
type MessageInfo struct {
	Media             *dto.Media
	Chat              *dto.Chat
	ConversationUsers dto.ConversationUserList
}

// MessageSenderResult TDB
type MessageSenderResult struct {
	Info   MessageInfo
	Events dto.EventList
	Jobs   []delayedjob.Job
}

// Send TBD
func (s *MessageSender) Send() (*MessageSenderResult, error) {
	if err := s.dbm.CreateMessage(s.message); err != nil {
		return nil, err
	}

	if s.message.PeerType == dto.PeerTypeChat {
		if err := s.sendToChat(); err != nil {
			return nil, err
		}
	} else if s.message.PeerType == dto.PeerTypeUser {
		if err := s.sendToDialog(); err != nil {
			return nil, err
		}
	}

	return &MessageSenderResult{
		Events: s.events,
		Jobs:   s.jobs,
	}, nil
}

func (s *MessageSender) sendToChat() error {
	var conversationUsers = s.info.ConversationUsers
	if len(s.info.ConversationUsers) == 0 {
		cusers, err := s.dbm.GetChatConversationUsers(s.message.PeerID)
		if err != nil {
			return err
		}
		if len(cusers) == 0 {
			return nil
		}
		conversationUsers = cusers
	}

	for _, cuser := range conversationUsers {
		if cuser.UserID == s.message.UserID {
			continue
		}

		unreadMessage := &dto.UnreadMessage{
			MessageID: s.message.ID,
			UserID:    cuser.UserID,
			PeerID:    s.message.PeerID,
			PeerType:  s.message.PeerType,
		}
		if err := s.dbm.SaveUnreadMessage(unreadMessage); err != nil {
			return err
		}

		if err := s.incrUserCounters(cuser.UserID, s.message.Peer()); err != nil {
			return err
		}
	}

	// save events about new chat creation
	if err := s.createEvents(conversationUsers.GetUsersIDs()); err != nil {
		return err
	}

	// push notifications
	cusers := make(dto.ConversationUserList, 0, len(conversationUsers))
	for _, cu := range conversationUsers {
		cusers = append(cusers, cu)
	}
	if len(cusers) == 0 {
		return nil
	}

	var accountedUsersIDs = util.Int64Slice(cusers.GetUsersIDs()).Delete(s.message.UserID)

	userID2Device := make(map[int64]dto.UserSessionDeviceList)
	for _, userID := range accountedUsersIDs {
		if devices, err := s.dbm.GetUserActiveDevices(userID); err != nil {
			return err
		} else if len(devices) != 0 {
			userID2Device[userID] = devices
		}
	}

	var chat *dto.Chat
	if chats, err := s.dbm.GetChatsByIDs([]int64{s.message.PeerID}); err != nil {
		return err
	} else if len(chats) < 1 {
		err := fmt.Errorf("cant find chat with id %d", s.message.PeerID)
		return err
	} else {
		chat = chats[0]
	}
	var chatTitle = chat.Name

	var media = s.info.Media
	if (s.message.Typo == dto.MessageTypeMedia) && (media == nil || media.ID != s.message.Media.MediaID) {
		medias, err := s.dbm.GetMediasByIDs([]int64{s.message.Media.MediaID})
		if err != nil || len(medias) != 1 {
			return err
		}
		media = medias[0]
	}

	var actorID *int64
	if s.message.Typo == dto.MessageTypeService {
		switch s.message.Service.Type {
		case dto.ServiceMessageTypeChatUserAdded:
			actorID = &s.message.Service.ChatUserAdded.AddedUserID
		case dto.ServiceMessageTypeChatUserDeleted:
			actorID = &s.message.Service.ChatUserDeleted.DeletedUserID
		}
	}

	var (
		authorName string
		actorName  string
		usersIDs   = make([]int64, 0, 2)
	)

	userID2Locale := make(map[int64]string, 0)
	usersIDs = append(usersIDs, s.message.UserID)

	if actorID != nil {
		usersIDs = append(usersIDs, *actorID)
	}
	if len(usersIDs) > 0 {
		users, err := s.dbm.GetUsersByIDs(usersIDs)
		if err != nil {
			return err
		}
		m := users.IDToUser()

		if authorName == "" {
			if author := m[s.message.UserID]; author != nil {
				authorName = author.FullName()
			} else {
				return nil
			}
		}

		if actorID != nil && actorName == "" {
			if actor := m[*actorID]; actor != nil {
				actorName = actor.FullName()
			} else {
				return nil
			}
		}

		for _, user := range users {
			userID2Locale[user.ID] = user.Locale
		}
	}

	for userID, devices := range userID2Device {
		if userID == s.message.UserID {
			continue
		}

		var alert string
		T := util.GetTranslationFunc(userID2Locale[userID])

		switch s.message.Typo {
		case dto.MessageTypeText:
			var text string
			if s.message.Text != nil {
				text = s.normalizeText(s.message.PlainText())
			}
			if text == "" {
				alert = T("api.messenger.text",
					map[string]interface{}{"SenderName": authorName, "ChatName": chatTitle})
			} else {
				alert = T("api.messenger.direct_text",
					map[string]interface{}{"SenderName": authorName, "ChatName": chatTitle, "Message": text})
			}
		case dto.MessageTypeMedia:
			switch media.Type {
			case dto.MediaTypePhoto:
				alert = T("%api.messenger.media.photo",
					map[string]interface{}{"SenderName": authorName, "ChatName": chatTitle})
			case dto.MediaTypeVideo:
				alert = T("api.messenger.media.video",
					map[string]interface{}{"SenderName": authorName, "ChatName": chatTitle})
			case dto.MediaTypeDocument:
				alert = T("api.messenger.media.document",
					map[string]interface{}{"SenderName": authorName, "ChatName": chatTitle})
			}
		case dto.MessageTypeService:
			switch s.message.Service.Type {
			case dto.ServiceMessageTypeChatCreated:
				alert = T("api.messenger.chat.created",
					map[string]interface{}{"SenderName": authorName, "ChatName": chatTitle})
			case dto.ServiceMessageTypeChatTitleChanged:
				newTitle := s.message.MessageExtra.Service.ChatTitleChanged.NewTitle
				alert = T("api.messenger.chat.title_changed",
					map[string]interface{}{"SenderName": authorName, "ChatName": newTitle})
			case dto.ServiceMessageTypeChatAvatarChanged:
				alert = T("api.messenger.chat.avatar_changed",
					map[string]interface{}{"SenderName": authorName, "ChatName": chatTitle})
			case dto.ServiceMessageTypeChatUserAdded:
				addedUserID := s.message.Service.ChatUserAdded.AddedUserID
				if addedUserID == userID {
					alert = T("api.messenger.chat.you_invited",
						map[string]interface{}{"SenderName": authorName, "ChatName": chatTitle})
				} else {
					alert = T("api.messenger.chat.user_invited",
						map[string]interface{}{"SenderName": authorName, "ChatName": chatTitle, "UserName": actorName})
				}
			case dto.ServiceMessageTypeChatUserDeleted:
				deletedUserID := s.message.Service.ChatUserDeleted.DeletedUserID
				if deletedUserID == userID {
					alert = T("api.messenger.chat.you_kicked",
						map[string]interface{}{"SenderName": authorName, "ChatName": chatTitle})
				} else {
					if s.message.UserID == deletedUserID {
						alert = T("api.messenger.chat.user_left",
							map[string]interface{}{"SenderName": authorName, "ChatName": chatTitle})
					} else {
						alert = T("api.messenger.chat.user_kicked",
							map[string]interface{}{"SenderName": authorName, "ChatName": chatTitle, "UserName": actorName})
					}
				}
			}
		}

		if alert != "" {
			var badge *int
			if counters, ok := s.countersByUserID[userID]; ok {
				value := int(counters.NUnreadNotifications)
				badge = &value
			}

			typo := dto.PushNotificationMessengerAction
			for _, device := range devices {
				job := NewPushNotificationJob(
					s.pushClient,
					device.Token,
					alert,
					badge,
					&models.PushNotification{
						Typo: (*int16)(&typo),
						MessengerAction: &models.PushNotificationMessengerAction{
							Peer: models.NewModelPeer(s.message.Peer()),
						},
					},
					device.Sandbox,
				)

				s.jobs = append(s.jobs, job)
			}
		}
	}

	return nil
}

func (s *MessageSender) sendToDialog() error {
	// implying only dialog
	fromUserID, toUserID := s.message.UserID, s.message.PeerID
	dialog, err := s.dbm.GetDialog(fromUserID, toUserID)
	if err != nil {
		s.logger.Error("GetDialog", "err", err)
		return err
	}

	if dialog == nil {
		least, greatest := util.Int64MinMax(fromUserID, toUserID)
		dialog := &dto.Dialog{
			LeastUserID:    least,
			GreatestUserID: greatest,
		}
		if err := s.dbm.CreateDialog(dialog); err != nil {
			s.logger.Error("CreateDialog fail", "err", err)
			return err
		}

		currentConversationUser := &dto.ConversationUser{
			UserID:   fromUserID,
			PeerID:   toUserID,
			PeerType: s.message.PeerType,
			JoinedAt: dialog.CreationTime,
		}
		if err := s.dbm.CreateConversationUser(currentConversationUser); err != nil {
			s.logger.Error("CreateConversationUser fail", "err", err)
			return err

		}

		peerConversationUser := &dto.ConversationUser{
			UserID:   toUserID,
			PeerID:   fromUserID,
			PeerType: s.message.PeerType,
			JoinedAt: dialog.CreationTime,
		}
		if err := s.dbm.CreateConversationUser(peerConversationUser); err != nil {
			s.logger.Error("CreateConversationUser fail", "err", err)
			return err
		}
	}

	unreadMessage := &dto.UnreadMessage{
		MessageID: s.message.ID,
		UserID:    toUserID,
		PeerID:    fromUserID,
		PeerType:  dto.PeerTypeUser,
	}
	if err := s.dbm.SaveUnreadMessage(unreadMessage); err != nil {
		s.logger.Error("CreateConversationUser fail", "err", err)
		return err
	}

	// update counters
	var selfPeer = &dto.Peer{ID: fromUserID, Type: dto.PeerTypeUser}
	if err := s.incrUserCounters(toUserID, selfPeer); err != nil {
		s.logger.Error("IncrUserCounters", "err", err)
		return err
	}

	// save events
	if err := s.createEvents([]int64{s.message.Peer().ID, s.message.UserID}); err != nil {
		return err
	}

	if (s.message.Typo == dto.MessageTypeMedia) && (s.info.Media == nil) {
		medias, err := s.dbm.GetMediasByIDs([]int64{s.message.Media.MediaID})
		if err != nil {
			s.logger.Error("Msg: get media", "err", err)
			return err
		}
		if len(medias) != 1 {
			s.logger.Error("Msg: media not found", "id", s.message.Media.MediaID)
			return errors.New("media not found")
		}
		s.info.Media = medias[0]
	}

	devices, err := s.dbm.GetUserActiveDevices(toUserID)
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		return nil
	}

	var author *dto.User
	if users, err := s.dbm.GetUsersByIDs([]int64{fromUserID}); err != nil {
		return err
	} else if len(users) != 1 {
		return errors.New("user not found")
	} else {
		author = users[0]
	}

	var targetUser *dto.User
	if users, err := s.dbm.GetUsersByIDs([]int64{toUserID}); err != nil {
		return err
	} else if len(users) != 1 {
		return errors.New("user not found")
	} else {
		targetUser = users[0]
	}

	// i18n
	T := util.GetTranslationFunc(targetUser.Locale)

	authorName := author.FullName()

	var alert string
	switch s.message.Typo {
	case dto.MessageTypeText:
		var text string
		if s.message.Text != nil {
			text = s.normalizeText(s.message.PlainText())
		}
		if text == "" {
			alert = T("api.messenger.dialog.empty_message",
				map[string]interface{}{"SenderName": authorName})
		} else {
			alert = fmt.Sprintf("%s: %s", authorName, text)
		}
	case dto.MessageTypeObject:
		alert = T("api.messenger.dialog.object_message",
			map[string]interface{}{"SenderName": authorName})
	case dto.MessageTypeMedia:
		if s.info.Media.Type.IsPhoto() {
			alert = T("api.messenger.dialog.photo_message",
				map[string]interface{}{"SenderName": authorName})
		} else if s.info.Media.Type.IsDocument() {
			alert = T("api.messenger.dialog.document_message",
				map[string]interface{}{"SenderName": authorName})
		} else if s.info.Media.Type.IsVideo() {
			alert = T("api.messenger.dialog.video_message",
				map[string]interface{}{"SenderName": authorName})
		}
	case dto.MessageTypeService:
		switch s.message.Service.Type {
		case dto.ServiceMessageTypeChatCreated:
			alert = T("New chat created")
		case dto.ServiceMessageTypeInvite:
			alert = T("api.messenger.dialog.new_invite",
				map[string]interface{}{"SenderName": authorName})
		case dto.ServiceMessageTypeInviteStatusChanged:
			if s.message.Service.InviteStatusChanged != nil {
				switch s.message.Service.InviteStatusChanged.InviteStatus {
				case dto.InviteAccepted:
					alert = T("Your invite has been accepted")
				case dto.InviteRejected:
					alert = T("Your invite has been rejected")
				case dto.InviteCanceled:
					alert = T("Your invite has been cancelled")
				}
			}
		case dto.ServiceMessageTypeTask:
			if s.message.Service.Task != nil {
				if s.message.Service.Task.NewTask {
					alert = T("api.messenger.dialog.task.created_assigned_to_you",
						map[string]interface{}{"TaskName": s.message.Service.Task.TaskTitle})
					break
				}
				if s.message.Service.Task.ChangedStatus {
					alert = T("api.messenger.dialog.task.status_changed",
						map[string]interface{}{"TaskName": s.message.Service.Task.TaskTitle, "TaskStatus": s.message.Service.Task.NewStatus.Text()})
					break
				}
				if s.message.Service.Task.ChangedAssignedUser {
					if s.message.Service.Task.NewAssignedUserID != nil {
						if *s.message.Service.Task.NewAssignedUserID == author.ID {
							alert = T("api.messenger.dialog.task.assigned_to_you",
								map[string]interface{}{"TaskName": s.message.Service.Task.TaskTitle})
						} else {
							alert = T("api.messenger.dialog.task.assigned_to_user",
								map[string]interface{}{"TaskName": s.message.Service.Task.TaskTitle})
						}
						break
					}
				}
				if s.message.Service.Task.ChangedArchive {
					if s.message.Service.Task.NewArchiveStatus == true {
						alert = T("api.messenger.dialog.task.archived",
							map[string]interface{}{"TaskName": s.message.Service.Task.TaskTitle})
					} else {
						alert = T("api.messenger.dialog.task.unarchived",
							map[string]interface{}{"TaskName": s.message.Service.Task.TaskTitle})
					}
					break
				}
				if s.message.Service.Task.WasChanged {
					alert = T("api.messenger.dialog.task.changed",
						map[string]interface{}{"TaskName": s.message.Service.Task.TaskTitle})
				}
				if s.message.Service.Task.WasDeleted {
					alert = T("api.messenger.dialog.task.removed",
						map[string]interface{}{"TaskName": s.message.Service.Task.TaskTitle})
					break
				}
			}
		}
	}

	var badge *int
	if counters, ok := s.countersByUserID[toUserID]; ok {
		value := int(counters.NUnreadNotifications)
		badge = &value
	}

	if alert != "" {
		typo := dto.PushNotificationMessengerAction

		for _, device := range devices {
			job := NewPushNotificationJob(
				s.pushClient,
				device.Token,
				alert,
				badge,
				&models.PushNotification{
					Typo: (*int16)(&typo),
					MessengerAction: &models.PushNotificationMessengerAction{
						Peer: models.NewModelPeer(selfPeer),
					},
				},
				device.Sandbox,
			)

			s.jobs = append(s.jobs, job)
		}
	}

	return nil
}

func (s *MessageSender) incrUserCounters(userID int64, peer *dto.Peer) error {
	var (
		counters MessageCounters
		n        *int32
		err      error
	)

	if n, err = s.dbm.UpdateUserNUnreadMessages(userID, +1); err != nil || n == nil {
		return err
	}
	counters.NUnreadMessages = *n

	if n, err = s.dbm.UpdateUserNUnreadNotifications(userID, +1); err != nil || n == nil {
		return err
	}
	counters.NUnreadNotifications = *n

	if n, err = s.dbm.UpdateConversationUserNUnreadMessages(userID, peer, +1); err != nil || n == nil {
		return err
	}
	counters.NConversationUnreadMessages = *n

	s.countersByUserID[userID] = counters

	return nil
}

func (s *MessageSender) createEvents(usersIDs []int64) error {
	var events dto.EventList
	for _, userID := range usersIDs {
		event := &dto.Event{
			UserID:       userID,
			Type:         dto.EventTypeNewMessage,
			CreationTime: time.Now(),
			EventUnion: dto.EventUnion{
				NewMessage: &dto.EventNewMessage{
					MessageID: s.message.ID,
				},
			},
		}
		if _, err := EmplaceEvent(s.dbm, event); err != nil {
			return err
		}

		events = append(events, event)
	}

	s.events = events

	return nil
}

const maxTextByteSize = 700

func (s *MessageSender) normalizeText(text string) string {
	if len(text) <= maxTextByteSize {
		if utf8.ValidString(text) {
			return text
		}
		return ""
	}

	var buf = bytes.NewBuffer(make([]byte, 0, maxTextByteSize))
	for _, r := range text {
		l := utf8.RuneLen(r)
		if l == -1 {
			return ""
		}
		if buf.Len()+l > maxTextByteSize {
			break
		}
		buf.WriteRune(r)
	}

	str := buf.String()
	if !utf8.ValidString(str) {
		return ""
	}
	return str
}

// EmplaceEvent TBD
func EmplaceEvent(dbm dal.Manager, event *dto.Event) (userFound bool, err error) {
	n, err := dbm.UpdateUserLastEventSeqNo(event.UserID, +1)
	if err != nil {
		return false, err
	}
	if n == nil {
		return false, nil
	}

	event.SeqNo = *n
	if err := dbm.SaveEvent(event); err != nil {
		return false, err
	}

	return true, nil
}

// NewPushNotificationJob TBD
func NewPushNotificationJob(pc npusherclient.Client, token, text string, badge *int, data *models.PushNotification, sandbox bool) delayedjob.Job {
	return delayedjob.NewJob(delayedjob.Immideate, func() {
		var collectorPayloadString string
		if data != nil {
			payload, err := json.Marshal(data)
			if err != nil {
				return
			}
			collectorPayloadString = base64.StdEncoding.EncodeToString(payload)
		}

		n := npusher.APNSNotification{
			Alert: text,
			Sound: "default",
			Badge: badge,
		}
		if collectorPayloadString != "" {
			n.Custom = map[string]interface{}{"collector_payload": collectorPayloadString}
		}

		if err := pc.SendPush(token, sandbox, n); err != nil {
			panic(err)
		}
	})
}
