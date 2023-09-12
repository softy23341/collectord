package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
	"git.softndit.com/collector/backend/util"
)

// CreateDialog TBD
func (m *Manager) CreateDialog(dialog *dto.Dialog) error {
	retFlds := dto.DialogFieldsList{dto.DialogFieldID, dto.DialogFieldCreationTime}
	insFlds := dto.DialogAllFields.Del(retFlds...)

	sql := fmt.Sprintf(`INSERT INTO dialog (%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		insFlds.JoinedNames(), insFlds.Placeholders(), retFlds.JoinedNames())

	return dto.QueryDialogRow(m.p, retFlds, sql, dialog.FieldsValues(insFlds)...).ScanTo(dialog)
}

// GetDialog TBD
func (m *Manager) GetDialog(firstUserID, secondUserID int64) (*dto.Dialog, error) {
	flds := dto.DialogAllFields

	sql := fmt.Sprintf(
		`SELECT %s FROM dialog WHERE least_user_id = $1 AND greatest_user_id = $2 LIMIT 1`,
		flds.JoinedNames())

	least, greatest := util.Int64MinMax(firstUserID, secondUserID)

	return dto.ScanDialog(m.p, flds, sql, least, greatest)
}
