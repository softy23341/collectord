package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// SaveUnreadMessage TBD
func (m *Manager) SaveUnreadMessage(unreadMessage *dto.UnreadMessage) error {
	flds := dto.UnreadMessageAllFields.Del(dto.UnreadMessageFieldID)

	sql := fmt.Sprintf(`INSERT INTO unread_message(%s) VALUES(%s)`,
		flds.JoinedNames(), flds.Placeholders())

	_, err := m.p.Exec(sql, unreadMessage.FieldsValues(flds)...)
	return err
}

// DelUnreadMessages TBD
func (m *Manager) DelUnreadMessages(userID int64, peer *dto.Peer, maxMessageID int64) (int32, int64, error) {
	sql := `
          WITH deleted AS (
            DELETE FROM unread_message
            WHERE user_id = $1 AND peer_id = $2 AND peer_type = $3 AND message_id <= $4
            RETURNING message_id
          )
          SELECT COUNT(message_id), COALESCE(MAX(message_id), 0)
          FROM deleted`
	var n, lastID int64
	err := m.p.QueryRow(sql, userID, peer.ID, peer.Type, maxMessageID).Scan(&n, &lastID)
	return int32(n), lastID, err
}

// DelAllUnreadMessagesByPeer TBD
func (m *Manager) DelAllUnreadMessagesByPeer(userID int64, peer *dto.Peer) (int32, error) {
	sql := `
          WITH deleted AS (
            DELETE FROM unread_message
            WHERE user_id = $1 AND peer_id = $2 AND peer_type = $3
            RETURNING message_id
          )
          SELECT COUNT(message_id)
          FROM deleted`
	var n int64
	err := m.p.QueryRow(sql, userID, peer.ID, peer.Type).Scan(&n)

	return int32(n), err
}
