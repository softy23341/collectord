package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetMessageByUserUniqID TBD
func (m *Manager) GetMessageByUserUniqID(userID, uniqID int64) (*dto.Message, error) {
	flds := dto.MessageAllFields

	sqlQuery := fmt.Sprintf(`
                SELECT %s
                FROM message
                WHERE user_id = $1
                AND user_uniq_id = $2
       `, flds.JoinedNames())

	return dto.ScanMessage(m.p, flds, sqlQuery, userID, uniqID)
}

// CreateMessage TBD
func (m *Manager) CreateMessage(msg *dto.Message) error {
	retFlds := dto.MessageFieldsList{dto.MessageFieldID, dto.MessageFieldCreationTime}
	insFlds := dto.MessageAllFields.Del(retFlds...)

	sql := fmt.Sprintf(`INSERT INTO message (%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		insFlds.JoinedNames(), insFlds.Placeholders(), retFlds.JoinedNames())

	return dto.QueryMessageRow(m.p, retFlds, sql, msg.FieldsValues(insFlds)...).ScanTo(msg)
}

// GetMessagesByIDs TBD
func (m *Manager) GetMessagesByIDs(ids []int64) (dto.MessageList, error) {
	flds := dto.MessageAllFields
	sql := fmt.Sprintf(`SELECT %s FROM message WHERE %s = any($1::bigint[])`,
		flds.JoinedNames(), dto.MessageFieldID.Name())

	return dto.ScanMessageList(m.p, flds, sql, ids)
}

// GetDialogMessagesByRange TBD
func (m *Manager) GetDialogMessagesByRange(firstUserID, secondUserID int64, p *dto.RangePaginator) (dto.MessageList, error) {
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
                         AND (
                           (user_id = $2 AND peer_id = $3)
                           OR
                           (user_id = $3 AND peer_id = $2)
                         ) ORDER BY id DESC
                         LIMIT $4`,
			flds.JoinedNames())

		args = []interface{}{int16(dto.PeerTypeUser), firstUserID, secondUserID, distance}
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
                         AND (
                           (user_id = $2 AND peer_id = $3)
                           OR
                           (user_id = $3 AND peer_id = $2)
                         )
                         AND (
                           id %s $4
                         )
                         ORDER BY id DESC
                         LIMIT $5`,
			flds.JoinedNames(), direction)
		args = []interface{}{int16(dto.PeerTypeUser), firstUserID, secondUserID, p.ID, distance}
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
