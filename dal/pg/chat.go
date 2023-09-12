package dalpg

import (
	"fmt"

	"github.com/jackc/pgx"

	"git.softndit.com/collector/backend/dto"
)

// GetChatByUserUniqID TBD
func (m *Manager) GetChatByUserUniqID(creatorUserID int64, uniqID int64) (*dto.Chat, error) {
	flds := dto.ChatAllFields

	sql := fmt.Sprintf(`
          SELECT %s
          FROM chat
          WHERE creator_user_id = $1
          AND user_uniq_id = $2
        `, flds.JoinedNames())

	return dto.ScanChat(m.p, flds, sql, creatorUserID, uniqID)
}

// CreateChat TBD
func (m *Manager) CreateChat(chat *dto.Chat) error {
	retFlds := dto.ChatFieldsList{dto.ChatFieldID, dto.ChatFieldCreationTime}
	insFlds := dto.ChatAllFields.Del(retFlds...)

	sql := fmt.Sprintf(`INSERT INTO chat (%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		insFlds.JoinedNames(), insFlds.Placeholders(), retFlds.JoinedNames())

	return dto.QueryChatRow(m.p, retFlds, sql, chat.FieldsValues(insFlds)...).ScanTo(chat)
}

// GetChatLastReadMessageID TBD
func (m *Manager) GetChatLastReadMessageID(chatID int64) (*int64, error) {
	sql := `SELECT last_read_message_id FROM chat WHERE id = $1 LIMIT 1`
	var id int64
	if err := m.p.QueryRow(sql, chatID).Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &id, nil
}

// SetChatLastReadMessageIDIfGreater TBD
func (m *Manager) SetChatLastReadMessageIDIfGreater(chatID, lastReadMessageID int64) (realLastReadMessageID *int64, err error) {
	sql := `UPDATE chat SET last_read_message_id = GREATEST(last_read_message_id, $1) WHERE id = $2 RETURNING last_read_message_id`
	var realID int64
	if err := m.p.QueryRow(sql, lastReadMessageID, chatID).Scan(&realID); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &realID, nil
}

// GetChatMessagesAuthors TBD
func (m *Manager) GetChatMessagesAuthors(chatID, startMessageID, endMessageID int64) (usersIDs []int64, err error) {
	empty := (startMessageID < 0) ||
		(endMessageID < 0) ||
		(startMessageID == 0 && endMessageID == 0) ||
		(endMessageID < startMessageID)
	if empty {
		return nil, nil
	}

	var sql = `SELECT DISTINCT(user_id) FROM message WHERE peer_id = $1 AND id BETWEEN $2 AND $3`
	rows, err := m.p.Query(sql, chatID, startMessageID, endMessageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int64
	for rows.Next() {
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		usersIDs = append(usersIDs, id)
	}

	return
}

// GetChatMessagesByRange TBD
func (m *Manager) GetChatMessagesByRange(chatID int64, p *dto.RangePaginator) (dto.MessageList, error) {
	flds := dto.MessageAllFields
	sql := ""
	var args []interface{}

	// abs distance
	distance := p.Distance
	if distance < 0 {
		distance = -distance
	}

	if p.ID == nil {
		sql = fmt.Sprintf(
			`SELECT %s
                         FROM message
                         WHERE peer_type = $1
                         AND peer_id = $2 ORDER BY id DESC
                         LIMIT $3`,
			flds.JoinedNames())

		args = []interface{}{int16(dto.PeerTypeChat), chatID, distance}
	} else {
		direction := " <"
		if p.Distance > 0 {
			direction = " >"
		}
		if p.Include {
			direction += "="
		}
		direction += " "

		sql = fmt.Sprintf(
			`SELECT %s
                         FROM message
                         WHERE peer_type = $1
                         AND peer_id = $2
                         AND (
                           id %s $3
                         )
                         ORDER BY id DESC
                         LIMIT $4`,
			flds.JoinedNames(), direction)
		args = []interface{}{int16(dto.PeerTypeChat), chatID, p.ID, distance}
	}

	messages, err := dto.ScanMessageList(m.p, flds, sql, args...)
	if err != nil {
		return nil, err
	}

	// from ID in direction; by default sort DESC
	if p.Distance > 0 {
		// reverse slice to desc order
		for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
			messages[i], messages[j] = messages[j], messages[i]
		}
	}

	return messages, nil
}

// GetChatsByIDs TBD
func (m *Manager) GetChatsByIDs(IDs []int64) (dto.ChatList, error) {
	flds := dto.ChatAllFields
	sql := fmt.Sprintf(`SELECT %s FROM chat WHERE %s = any($1::bigint[])`,
		flds.JoinedNames(), dto.ChatFieldID.Name())

	return dto.ScanChatList(m.p, flds, sql, IDs)
}

// UpdateChat TBD
func (m *Manager) UpdateChat(chat *dto.Chat) error {
	sql := `
	  UPDATE "chat"
          SET
            creator_user_id = $1,
            admin_user_id = $2,
            name = $3,
            avatar_media_id = $4
          WHERE id = $5
        `

	_, err := m.p.Exec(sql,
		chat.CreatorUserID,
		chat.AdminUserID,
		chat.Name,
		chat.AvatarMediaID,
		chat.ID,
	)
	return err
}

// GetChatMembersCntByChat TBD
func (m *Manager) GetChatMembersCntByChats(chatsIDs []int64) (map[int64]int64, error) {
	sql := `
	  SELECT c.id, COUNT(cu.user_id) AS cnt
	  FROM chat AS c
	  INNER JOIN conversation_user AS cu
	    ON cu.peer_id = c.id AND cu.peer_type = 2
	  WHERE c.id = any($1::bigint[])
	  GROUP BY c.id
	`

	rows, err := m.p.Query(sql, chatsIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		chatID int64
		cnt    int64
	)
	chat2cnt := make(map[int64]int64)
	for rows.Next() {
		if err := rows.Scan(&chatID, &cnt); err != nil {
			return nil, err
		}
		chat2cnt[chatID] = cnt
	}
	return chat2cnt, nil
}
