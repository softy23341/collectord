package dto

import (
	"fmt"
	"time"
)

//go:generate dbgen -type Message

// MessageType TBD
type MessageType int16

// Message types
const (
	MessageTypeText    MessageType = 0
	MessageTypeMedia   MessageType = 10
	MessageTypeService MessageType = 20
	MessageTypeObject  MessageType = 30
)

// MessageTypeSet TBD
var MessageTypeSet = map[MessageType]bool{
	MessageTypeText:    true,
	MessageTypeMedia:   true,
	MessageTypeService: true,
	MessageTypeObject:  true,
}

// NewMessageType TBD
func NewMessageType(t int16) (MessageType, error) {
	typo := MessageType(t)
	if _, f := MessageTypeSet[typo]; !f {
		return typo, fmt.Errorf("cat cast %d to messagetype", t)
	}
	return typo, nil
}

// ServiceMessageType TBD
type ServiceMessageType int16

// ServiceMessageType types
const (
	ServiceMessageTypeChatCreated         ServiceMessageType = 1
	ServiceMessageTypeChatUserAdded       ServiceMessageType = 2
	ServiceMessageTypeChatUserDeleted     ServiceMessageType = 3
	ServiceMessageTypeChatTitleChanged    ServiceMessageType = 4
	ServiceMessageTypeChatAvatarChanged   ServiceMessageType = 5
	ServiceMessageTypeInvite              ServiceMessageType = 6
	ServiceMessageTypeTask                ServiceMessageType = 7
	ServiceMessageTypeInviteStatusChanged ServiceMessageType = 8
)

// Message TBD
type Message struct {
	ID           int64       `db:"id"`
	UserID       int64       `db:"user_id"`
	UserUniqID   int64       `db:"user_uniq_id"`
	PeerID       int64       `db:"peer_id"`
	PeerType     PeerType    `db:"peer_type"`
	Typo         MessageType `db:"type"`
	CreationTime time.Time   `db:"creation_time"`

	MessageExtra `db:"extra,json"`
}

// PlainText TBD
func (e MessageExtra) PlainText() string {
	if e.Text == nil {
		return ""
	}
	text := e.Text.Text

	return text
}

// IsText TBD
func (m *Message) IsText() bool {
	return m.Typo == MessageTypeText
}

// IsMedia TBD
func (m *Message) IsMedia() bool {
	return m.Typo == MessageTypeMedia
}

// IsObject TBD
func (m *Message) IsObject() bool {
	return m.Typo == MessageTypeObject
}

// IsService TBD
func (m *Message) IsService() bool {
	return m.Typo == MessageTypeService
}

// Peer TBD
func (m *Message) Peer() *Peer {
	return &Peer{ID: m.PeerID, Type: m.PeerType}
}

// MediasIDs TBD
func (m *Message) MediasIDs() []int64 {
	var ids []int64

	if m.IsMedia() {
		ids = append(ids, m.Media.MediaID)
	} else if m.IsService() {
		if avatarMsg := m.Service.ChatAvatarChanged; avatarMsg != nil && avatarMsg.NewAvatarID != nil {
			ids = append(ids, *avatarMsg.NewAvatarID)
		}
	}

	return ids
}

// MessageList TBD
type MessageList []*Message

// UsersIDs TBD
func (ml MessageList) UsersIDs() []int64 {
	usersIDsSet := make(map[int64]struct{})
	for _, msg := range ml {
		usersIDsSet[msg.UserID] = struct{}{}
		if msg.PeerType.IsUser() {
			usersIDsSet[msg.PeerID] = struct{}{}
		}
		if msg.IsService() {
			if newUserMsg := msg.Service.ChatUserAdded; newUserMsg != nil {
				usersIDsSet[newUserMsg.AddedUserID] = struct{}{}
			}
			if deletedUserMsg := msg.Service.ChatUserDeleted; deletedUserMsg != nil {
				usersIDsSet[deletedUserMsg.DeletedUserID] = struct{}{}
			}
			if taskMsg := msg.Service.Task; taskMsg != nil {
				usersIDsSet[taskMsg.ActorUserID] = struct{}{}
				if taskMsg.NewAssignedUserID != nil {
					usersIDsSet[*taskMsg.NewAssignedUserID] = struct{}{}
				}
			}
		}
	}

	var usersIDs []int64
	for userID := range usersIDsSet {
		usersIDs = append(usersIDs, userID)
	}
	return usersIDs
}

// InvitesIDs TBD
func (ml MessageList) InvitesIDs() []int64 {
	invitesIDsSet := make(map[int64]struct{})
	for _, msg := range ml {
		if msg.IsService() {
			if newInviteMsg := msg.Service.Invite; newInviteMsg != nil {
				invitesIDsSet[newInviteMsg.InviteID] = struct{}{}
			}
			if newInviteChangedMsg := msg.Service.InviteStatusChanged; newInviteChangedMsg != nil {
				invitesIDsSet[newInviteChangedMsg.InviteID] = struct{}{}
			}
		}
	}

	var invitesIDs []int64
	for inviteID := range invitesIDsSet {
		invitesIDs = append(invitesIDs, inviteID)
	}
	return invitesIDs
}

// ChatsIDs TBD
func (ml MessageList) ChatsIDs() []int64 {
	chatsIDsSet := make(map[int64]struct{})
	for _, msg := range ml {
		if msg.PeerType.IsChat() {
			chatsIDsSet[msg.PeerID] = struct{}{}
		}
	}

	var chatsIDs []int64
	for chatID := range chatsIDsSet {
		chatsIDs = append(chatsIDs, chatID)
	}
	return chatsIDs
}

// ObjectIDs TBD
func (ml MessageList) ObjectIDs() []int64 {
	objectsIDsSet := make(map[int64]struct{})
	for _, msg := range ml {
		if msg.IsObject() {
			objectsIDsSet[msg.Object.ObjectID] = struct{}{}
		}
	}

	var objectsIDs []int64
	for objectID := range objectsIDsSet {
		objectsIDs = append(objectsIDs, objectID)
	}
	return objectsIDs
}

// GetTaskIDs TBD
func (ml MessageList) GetTaskIDs() []int64 {
	tasksIDsSet := make(map[int64]struct{})
	for _, msg := range ml {
		if msg.IsService() {
			if taskMsg := msg.Service.Task; taskMsg != nil && taskMsg.TaskID != nil {
				tasksIDsSet[*taskMsg.TaskID] = struct{}{}
			}
		}
	}

	var tasksIDs []int64
	for taskID := range tasksIDsSet {
		tasksIDs = append(tasksIDs, taskID)
	}
	return tasksIDs
}

// MediasIDs TBD
func (ml MessageList) MediasIDs() []int64 {
	var mediasIDs []int64
	for _, media := range ml {
		if ids := media.MediasIDs(); len(ids) > 0 {
			mediasIDs = append(mediasIDs, ids...)
		}
	}
	return mediasIDs
}

// MessageExtra TBD
type MessageExtra struct {
	Forward *MessageForward `json:"forward,omitempty"`
	Text    *TextMessage    `json:"text,omitempty"`
	Media   *MediaMessage   `json:"media,omitempty"`
	Service *ServiceMessage `json:"service,omitempty"`
	Object  *ObjectMessage  `json:"object,omitempty"`
}

// ObjectMessage TBD
type ObjectMessage struct {
	ObjectID int64 `db:"object_id"`
}

// ServiceMessage TBD
type ServiceMessage struct {
	Type ServiceMessageType `json:"type"`

	ChatCreated         *ServiceMessageChatCreated         `json:"chat_created,omitempty"`
	ChatUserAdded       *ServiceMessageChatUserAdded       `json:"chat_user_added,omitempty"`
	ChatUserDeleted     *ServiceMessageChatUserDeleted     `json:"chat_user_deleted,omitempty"`
	ChatTitleChanged    *ServiceMessageChatTitleChanged    `json:"chat_title_changed,omitempty"`
	ChatAvatarChanged   *ServiceMessageChatAvatarChanged   `json:"chat_avatar_changed,omitempty"`
	Invite              *ServiceMessageInvite              `json:"invite,omitempty"`
	Task                *ServiceMessageTask                `json:"task,omitempty"`
	InviteStatusChanged *ServiceMessageInviteStatusChanged `json:"invite_status_changed,omitempty"`
}

// ServiceMessageTask TBD
type ServiceMessageTask struct {
	TaskID    *int64 `json:"task_id"`
	TaskTitle string `json:"task_title"`

	ActorUserID   int64  `json:"actor_user_id"`
	ActorUserName string `json:"actor_user_name"`

	NewTask bool `json:"new_task"`

	NewAssignedUserID   *int64 `json:"new_assigned_user_id"`
	ChangedAssignedUser bool   `json:"changed_assigned_user"`

	NewStatus     TaskStatus `json:"new_status"`
	ChangedStatus bool       `json:"changed_status"`

	NewArchiveStatus bool `json:"new_archive_status"`
	ChangedArchive   bool `json:"changed_archive"`

	WasChanged bool `json:"was_changed"`
	WasDeleted bool `json:"was_deleted"`
}

// ServiceMessageInvite TBD
type ServiceMessageInvite struct {
	InviteID int64
}

// ServiceMessageChatCreated TBD
type ServiceMessageChatCreated struct {
	ChatTitle string `json:"chat_title"`
}

// ServiceMessageChatUserAdded TBD
type ServiceMessageChatUserAdded struct {
	AddedUserID int64 `json:"added_user_id"`
}

// ServiceMessageChatUserDeleted TBD
type ServiceMessageChatUserDeleted struct {
	DeletedUserID int64 `json:"deleted_user_id"`
}

// ServiceMessageChatTitleChanged TBD
type ServiceMessageChatTitleChanged struct {
	NewTitle string `json:"new_title"`
}

// ServiceMessageChatAvatarChanged TBD
type ServiceMessageChatAvatarChanged struct {
	NewAvatarID *int64 `json:"new_avatar_id,omitempty"`
}

// MessageForward TBD
type MessageForward struct {
	OriginalID     int64 `json:"original_id"`
	OriginalUserID int64 `json:"original_user_id"`
}

// TextMessage TBD
type TextMessage struct {
	Text string `json:"text,omitempty"`
}

// MediaMessage TBD
type MediaMessage struct {
	MediaID int64 `json:"media_id"`
}

// InviteStatusChanged TBD
type ServiceMessageInviteStatusChanged struct {
	InviteID     int64        `json:"invite_id"`
	InviteStatus InviteStatus `json:"invite_status"`
}
