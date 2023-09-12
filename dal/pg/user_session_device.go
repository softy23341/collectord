package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// DeleteSessionDevices TBD
func (m *Manager) DeleteSessionDevices(sessionTokens []string) error {
	sql := `
            DELETE FROM user_session_device WHERE session_id IN (
              SELECT id FROM user_session WHERE auth_token = any($1::character varying[])
            );
        `
	_, err := m.p.Exec(sql, sessionTokens)
	return err
}

// SetSessionDeviceByAuthToken TBD
func (m *Manager) SetSessionDeviceByAuthToken(device *dto.UserSessionDevice) error {
	flds := dto.UserSessionDeviceAllFields.Del(dto.UserSessionDeviceFieldID)

	sql := fmt.Sprintf(`INSERT INTO "user_session_device"(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		flds.JoinedNames(),
		flds.Placeholders(),
		dto.UserSessionDeviceFieldID.Name(),
	)

	return dto.QueryUserSessionDeviceRow(
		m.p,
		dto.UserSessionDeviceFieldsList{dto.UserSessionDeviceFieldID},
		sql,
		(*dto.UserSessionDevice)(device).FieldsValues(flds)...,
	).ScanTo(device)
}

// GetUserActiveDevices TBD
func (m *Manager) GetUserActiveDevices(userID int64) (dto.UserSessionDeviceList, error) {
	flds := dto.UserSessionDeviceAllFields

	sql := fmt.Sprintf(`
               SELECT %s
               FROM user_session_device AS usd
               INNER JOIN (
                       SELECT MAX(usd.id) AS id
                       FROM user_session_device AS usd
                       INNER JOIN user_session AS us
                               ON us.id = usd.session_id
                       WHERE us.user_id = $1
                       GROUP BY token
               ) AS usd_sub ON usd_sub.id = usd.id;
        `, flds.JoinedNamesWithAlias("usd"))

	return dto.ScanUserSessionDeviceList(m.p, flds, sql, userID)
}
