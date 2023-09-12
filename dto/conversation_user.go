package dto

import "time"

//go:generate dbgen -type ConversationUser

// ConversationUser TBD
type ConversationUser struct {
	ID                int64     `db:"id"`
	UserID            int64     `db:"user_id"`
	PeerID            int64     `db:"peer_id"`
	PeerType          PeerType  `db:"peer_type"`
	JoinedAt          time.Time `db:"joined_at"`
	LastReadMessageID *int64    `db:"last_read_message_id"`
	NUnreadMessages   int32     `db:"n_unread_messages"`
	InvitedUserID     *int64    `db:"inviter_user_id"`
	CreationTime      time.Time `db:"creation_time"`
}

// Peer TBD
func (c *ConversationUser) Peer() Peer {
	return Peer{
		ID:   c.PeerID,
		Type: c.PeerType,
	}
}

// ConversationUserList TBD
type ConversationUserList []*ConversationUser

// GetPeers TBD
func (cl ConversationUserList) GetPeers() PeerList {
	pl := make(PeerList, len(cl))
	for idx, cu := range cl {
		peer := cu.Peer()
		pl[idx] = &peer
	}
	return pl
}

// GetUsersIDs TBD
func (cl ConversationUserList) GetUsersIDs() []int64 {
	var sz int
	for _, cu := range cl {
		if cu.PeerType == PeerTypeUser {
			sz++
		}
	}
	var ids = make([]int64, len(cl)+sz)
	for _, cu := range cl {
		ids = append(ids, cu.UserID)
		if cu.PeerType == PeerTypeUser {
			ids = append(ids, cu.PeerID)
		}
	}
	return ids
}

// GetUsersByPeer TBD
func (cl ConversationUserList) GetUsersByPeer(peerID int64, peerType PeerType) ConversationUserList {
	users := make(ConversationUserList, 0, len(cl))
	for _, user := range cl {
		if user.PeerID == peerID && user.PeerType == peerType {
			users = append(users, user)
		}
	}
	return users
}
