package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// CreateUserRootRef TBD
func (m *Manager) CreateUserRootRef(ref *dto.UserRootRef) error {
	flds := dto.UserRootRefAllFields

	sql := fmt.Sprintf(`INSERT INTO "user_root_ref"(%[1]s) VALUES(%[2]s)`,
		flds.JoinedNames(),
		flds.Placeholders(),
	)

	_, err := m.p.Exec(sql, (*dto.UserRootRef)(ref).FieldsValues(flds)...)
	return err
}

// GetUserRelatedRootRefs TBD
func (m *Manager) GetUserRelatedRootRefs(userIDs []int64) (dto.UserRootRefList, error) {
	flds := dto.UserRootRefAllFields

	sql := fmt.Sprintf(`
          SELECT DISTINCT %s
          FROM user_root_ref AS user_root
          INNER JOIN user_root_ref AS users_root
            ON users_root.root_id = user_root.root_id
          WHERE user_root.user_id = any($1::bigint[])
        `, flds.JoinedNamesWithAlias("users_root"))

	return dto.ScanUserRootRefList(m.p, flds, sql, userIDs)
}

// GetUserRootRefs TBD
func (m *Manager) GetUserRootRefs(rootsIDs []int64) (dto.UserRootRefList, error) {
	flds := dto.UserRootRefAllFields

	sql := fmt.Sprintf(`
          SELECT %s
          FROM user_root_ref
          WHERE root_id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanUserRootRefList(m.p, flds, sql, rootsIDs)
}

// DeleteUserRootRef TBD
func (m *Manager) DeleteUserRootRef(userID, rootID int64) error {
	sql := "DELETE FROM user_root_ref WHERE user_id = $1 AND root_id = $2"

	_, err := m.p.Exec(sql, userID, rootID)
	return err
}
