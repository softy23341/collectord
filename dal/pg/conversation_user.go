package dalpg

import (
	"fmt"

	"github.com/jackc/pgx"

	"git.softndit.com/collector/backend/dto"
)

// CreateConversationUser TBD
func (m *Manager) CreateConversationUser(user *dto.ConversationUser) error {
	if user.JoinedAt.IsZero() {
		flds := dto.ConversationUserAllFields.Del(dto.ConversationUserFieldID, dto.ConversationUserFieldCreationTime)

		sql := fmt.Sprintf(`INSERT INTO conversation_user(%s) VALUES(%s)`,
			flds.JoinedNames(), flds.Placeholders())

		_, err := m.p.Exec(sql, user.FieldsValues(flds)...)

		return err
	}

	retFlds := dto.ConversationUserFieldsList{
		dto.ConversationUserFieldJoinedAt,
		dto.ConversationUserFieldCreationTime,
	}
	insFlds := dto.ConversationUserAllFields.
		Del(retFlds...).
		Del(dto.ConversationUserFieldID)

	sql := fmt.Sprintf(`INSERT INTO conversation_user (%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		insFlds.JoinedNames(), insFlds.Placeholders(), retFlds.JoinedNames())

	return dto.QueryConversationUserRow(m.p, retFlds, sql, user.FieldsValues(insFlds)...).ScanTo(user)
}

// GetUserPeers TBD
func (m *Manager) GetUserPeers(userID int64) (dto.PeerList, error) {
	flds := dto.ConversationUserAllFields

	sql := fmt.Sprintf(`SELECT %s FROM conversation_user WHERE user_id = $1`,
		flds.JoinedNames())

	cusers, err := dto.ScanConversationUserList(m.p, flds, sql, userID)
	if err != nil {
		return nil, err
	}

	return dto.ConversationUserList(cusers).GetPeers(), nil
}

// UpdateConversationUserNUnreadMessages TBD
func (m *Manager) UpdateConversationUserNUnreadMessages(userID int64, peer *dto.Peer, delta int32) (*int32, error) {
	sql := `
          UPDATE conversation_user
          SET n_unread_messages = GREATEST(0, n_unread_messages + $1)
          WHERE user_id = $2 AND peer_id = $3 AND peer_type = $4
          RETURNING n_unread_messages`

	var n int32
	err := m.p.QueryRow(sql, delta, userID, peer.ID, int16(peer.Type)).Scan(&n)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &n, nil
}

// UpdateConversationUserLastReadMessageID TBD
func (m *Manager) UpdateConversationUserLastReadMessageID(userID int64, peer *dto.Peer, lastReadMessageID int64) (*int64, error) {
	sql := `UPDATE conversation_user
                SET last_read_message_id = $1 WHERE user_id = $2 AND peer_id = $3 AND peer_type = $4
                RETURNING last_read_message_id`

	var id int64
	err := m.p.QueryRow(sql, lastReadMessageID, userID, peer.ID, int16(peer.Type)).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &id, nil
}

// GetChatConversationUsers TBD
func (m *Manager) GetChatConversationUsers(chatID int64) (dto.ConversationUserList, error) {
	flds := dto.ConversationUserAllFields

	sql := fmt.Sprintf(
		`SELECT %s
                 FROM conversation_user
                 WHERE peer_id = $1
                 AND peer_type = $2`,
		flds.JoinedNames())

	return dto.ScanConversationUserList(m.p, flds, sql, chatID, int16(dto.PeerTypeChat))
}

// GetDialogConversationUsers TBD
func (m *Manager) GetDialogConversationUsers(firstUserID, secondUserID int64) (dto.ConversationUserList, error) {
	flds := dto.ConversationUserAllFields

	sql := fmt.Sprintf(
		`SELECT %s
                 FROM conversation_user
                 WHERE (user_id = $1 AND peer_id = $2 AND peer_type = $3)
                    OR (user_id = $2 AND peer_id = $1 AND peer_type = $3)`,
		flds.JoinedNames())

	return dto.ScanConversationUserList(m.p, flds, sql, firstUserID, secondUserID, int16(dto.PeerTypeUser))
}

// DelConversationUser TBD
func (m *Manager) DelConversationUser(userID int64, peer dto.Peer) error {
	sql := `DELETE FROM conversation_user WHERE user_id = $1 AND peer_id = $2 AND peer_type = $3`
	_, err := m.p.Exec(sql, userID, peer.ID, peer.Type)
	return err
}
