package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// CreateUserBan TBD
func (m *Manager) CreateUserBan(ban *dto.UserBan) error {
	retFlds := dto.UserBanFieldsList{dto.UserBanFieldID, dto.UserBanFieldCreationTime}
	insFlds := dto.UserBanAllFields.Del(retFlds...)

	sql := fmt.Sprintf(`INSERT INTO user_ban (%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		insFlds.JoinedNames(), insFlds.Placeholders(), retFlds.JoinedNames())

	return dto.QueryUserBanRow(m.p, retFlds, sql, ban.FieldsValues(insFlds)...).ScanTo(ban)
}

// DeleteUserBan TBD
func (m *Manager) DeleteUserBan(creatorUserID, userID int64) error {
	sql := fmt.Sprintf(`DELETE FROM user_ban WHERE creator_user_id = $1 AND user_id = $2`)
	_, err := m.p.Exec(sql, creatorUserID, userID)
	return err
}

// GetUserBanList TBD
func (m *Manager) GetUserBanList(userID int64) (dto.UserBanList, error) {
	flds := dto.UserBanAllFields
	sql := fmt.Sprintf(`SELECT %s FROM user_ban WHERE creator_user_id = $1`, flds.JoinedNames())

	return dto.ScanUserBanList(m.p, flds, sql, userID)
}

// GetUserBanList TBD
func (m *Manager) IsUserBanned(creatorUserID, userID int64) (bool, error) {
	sql := `SELECT COUNT(user_ban.*) AS cnt FROM user_ban WHERE creator_user_id = $1 AND user_id = $2`
	rows, err := m.p.Query(sql, creatorUserID, userID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var cnt int64
	for rows.Next() {
		if err := rows.Scan(&cnt); err != nil {
			return false, err
		}
	}
	if rows.Err() != nil {
		return false, err
	}
	return cnt > 0, nil
}
