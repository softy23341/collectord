package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// CreateInvite TBD
func (m *Manager) CreateInvite(invite *dto.Invite) error {
	flds := dto.InviteAllFields.Del(dto.InviteFieldID)

	sql := fmt.Sprintf(`INSERT INTO "invite"(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		flds.JoinedNames(),
		flds.Placeholders(),
		dto.InviteFieldID.Name(),
	)

	return dto.QueryInviteRow(
		m.p,
		dto.InviteFieldsList{dto.InviteFieldID},
		sql,
		(*dto.Invite)(invite).FieldsValues(flds)...,
	).ScanTo(invite)
}

// GetInvitesByIDs TBD
func (m *Manager) GetInvitesByIDs(inviteIDs []int64) (dto.InviteList, error) {
	if len(inviteIDs) == 0 {
		return nil, nil
	}

	flds := dto.InviteAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM invite
                WHERE id = any($1::bigint[])
				AND to_user_id IS NOT NULL 
                ORDER BY id DESC
        `, flds.JoinedNames())

	return dto.ScanInviteList(m.p, flds, sql, inviteIDs)
}

// ChangeInviteStatus TBD
func (m *Manager) ChangeInviteStatus(inviteID int64, status dto.InviteStatus) error {
	sql := `
          UPDATE invite
          SET status = $1
          WHERE id = $2
        `

	_, err := m.p.Exec(sql, int16(status), inviteID)
	return err
}

// ChangeInviteToUserID TBD
func (m *Manager) ChangeInviteToUserID(inviteID int64, toUserID int64) error {
	sql := `
          UPDATE invite
          SET to_user_id = $1
          WHERE id = $2
        `

	_, err := m.p.Exec(sql, toUserID, inviteID)
	return err
}

// GetInviteByUserRoot TBD
func (m *Manager) GetInviteByUserRoot(fromUserID, toUserID, rootID int64, status dto.InviteStatus) (*dto.Invite, error) {
	flds := dto.InviteAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM "invite"
                WHERE creator_user_id = $1
                AND to_user_id = $2
                AND root_id = $3
                AND status = $4
        `, flds.JoinedNames())

	return dto.ScanInvite(m.p, flds, sql, fromUserID, toUserID, rootID, int16(status))
}

// GetInviteByToken TBD
func (m *Manager) GetInviteByToken(token string, status dto.InviteStatus) (*dto.Invite, error) {
	flds := dto.InviteAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM "invite"
                WHERE token = $1
                AND status = $2
        `, flds.JoinedNames())

	return dto.ScanInvite(m.p, flds, sql, token, int16(status))
}

// GetInvitesByRoot TBD
func (m *Manager) GetInvitesByRoot(rootID int64, status dto.InviteStatus) (dto.InviteList, error) {
	flds := dto.InviteAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM "invite"
                WHERE
                root_id = $1
                AND status = $2
                ORDER BY id DESC
        `, flds.JoinedNames())

	return dto.ScanInviteList(m.p, flds, sql, rootID, int16(status))
}
