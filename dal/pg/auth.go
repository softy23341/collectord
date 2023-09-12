package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// CreateUserSession TBD
func (m *Manager) CreateUserSession(userSession *dto.UserSession) error {
	retFlds := dto.UserSessionFieldsList{dto.UserSessionFieldID, dto.UserSessionFieldCreationTime}
	insFlds := dto.UserSessionAllFields.Del(retFlds...)

	sql := fmt.Sprintf(`INSERT INTO "user_session"(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		insFlds.JoinedNames(),
		insFlds.Placeholders(),
		retFlds.JoinedNames(),
	)

	return dto.QueryUserSessionRow(
		m.p,
		retFlds,
		sql,
		(*dto.UserSession)(userSession).FieldsValues(insFlds)...,
	).ScanTo(userSession)
}

// GetUsersBySessionTokens TBD
func (m *Manager) GetUsersBySessionTokens(tokens []string) (dto.UserList, error) {
	flds := dto.UserAllFields
	sql := fmt.Sprintf(`
                SELECT %s
                FROM "user" AS u
                INNER JOIN "user_session" AS us
                  ON u.id = us.user_id
                WHERE us.auth_token = any($1::character varying[])
        `, flds.JoinedNamesWithAlias("u"))

	return dto.ScanUserList(m.p, flds, sql, tokens)
}

// DestroyUserSessionByTokens TBD
func (m *Manager) DestroyUserSessionByTokens(tokens []string) error {
	sql := `
          DELETE FROM "user_session"
          WHERE auth_token = any($1::character varying[])
        `
	_, err := m.p.Exec(sql, tokens)
	return err
}

// GetSessionsBySessionTokens TBD
func (m *Manager) GetSessionsBySessionTokens(tokens []string) (dto.UserSessionList, error) {
	flds := dto.UserSessionAllFields
	sql := fmt.Sprintf(`
          SELECT %s
          FROM "user_session"
          WHERE auth_token = any($1::character varying[])
	`, flds.JoinedNames())

	return dto.ScanUserSessionList(m.p, flds, sql, tokens)
}
