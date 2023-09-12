package dalpg

import (
	"fmt"
	"time"

	"git.softndit.com/collector/backend/dto"
)

// SaveEvent TBD
func (m *Manager) SaveEvent(event *dto.Event) error {
	flds := dto.EventAllFields.Del(dto.EventFieldID)

	sql := fmt.Sprintf(`INSERT INTO event(%s) VALUES(%s)`,
		flds.JoinedNames(), flds.Placeholders())

	_, err := m.p.Exec(sql, event.FieldsValues(flds)...)

	return err
}

// GetEvents TBD
func (m *Manager) GetEvents(userID, beginSeqNo int64, endSeqNo *int64, statuses []dto.EventStatus) (dto.EventList, error) {
	flds := dto.EventAllFields

	statusClosure := func(placeHolder string) string {
		if len(statuses) != 0 {
			return fmt.Sprintf("AND status = any(%s::smallint[])", placeHolder)
		}
		return ""
	}

	intStatuses := make([]int16, len(statuses))
	for i, s := range statuses {
		intStatuses[i] = int16(s)
	}

	var sql string
	var args []interface{}
	if endSeqNo == nil {
		sql = fmt.Sprintf(
			`
                         SELECT %s
                         FROM event
                         WHERE user_id = $1 AND seq_no >= $2 %s
                         ORDER BY seq_no`,
			flds.JoinedNames(), statusClosure("$3"))
		args = []interface{}{userID, beginSeqNo}
		if len(statuses) != 0 {
			args = append(args, intStatuses)
		}
	} else {
		sql = fmt.Sprintf(
			`
                        SELECT %s
                        FROM event
                        WHERE user_id = $1 AND (seq_no BETWEEN $2 AND $3) %s
                        ORDER BY seq_no`,
			flds.JoinedNames(), statusClosure("$4"))
		args = []interface{}{userID, beginSeqNo, *endSeqNo}
		if len(statuses) != 0 {
			args = append(args, intStatuses)
		}
	}

	return dto.ScanEventList(m.p, flds, sql, args...)
}

// ConfirmEvents TBD
func (m *Manager) ConfirmEvents(userID int64, endSeqNo int64) error {
	sql := `UPDATE event SET status = $1 WHERE user_id = $2 AND seq_no <= $3`

	_, err := m.p.Exec(sql, int16(dto.EventStatusConfirmed), userID, endSeqNo)

	return err
}

// DeleteOldEvents TBD
func (m *Manager) DeleteOldEvents(ago time.Duration) error {
	sql := `DELETE FROM event WHERE status = $1 AND creation_time < $2`

	_, err := m.p.Exec(sql, int16(dto.EventStatusConfirmed), time.Now().Add(-ago))
	return err
}

// GetOldEventsCnt TBD
func (m *Manager) GetOldEventsCnt(ago time.Duration) (int32, error) {
	sql := `SELECT COUNT(*) AS cnt FROM event WHERE status = $1 AND creation_time < $2`

	var n int32
	err := m.p.QueryRow(sql, int16(dto.EventStatusConfirmed), time.Now().Add(-ago)).Scan(&n)
	return n, err
}
