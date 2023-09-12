package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetRootsByUserID TBD
func (m *Manager) GetRootsByUserID(userID int64) (dto.RootList, error) {
	flds := dto.RootAllFields
	sql := fmt.Sprintf(`
SELECT %s FROM "root" AS r
INNER JOIN user_root_ref AS ur
  ON ur.root_id = r.id
WHERE ur.user_id = $1
`, flds.JoinedNamesWithAlias("r"))

	return dto.ScanRootList(m.p, flds, sql, userID)
}

// GetMainUserRoot TBD
func (m *Manager) GetMainUserRoot(userID int64) (dto.RootList, error) {
	flds := dto.RootAllFields
	sql := fmt.Sprintf(`
SELECT %s FROM "root" AS r
INNER JOIN user_root_ref AS ur
  ON ur.root_id = r.id
WHERE ur.user_id = $1
AND ur.typo = $2
`, flds.JoinedNamesWithAlias("r"))

	return dto.ScanRootList(m.p, flds, sql, userID, dto.UserRootTypeOwner)
}

// GetRootsByIDs TBD
func (m *Manager) GetRootsByIDs(rootIDs []int64) (dto.RootList, error) {
	flds := dto.RootAllFields
	sql := fmt.Sprintf(`
           SELECT %s
           FROM "root"
           WHERE id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanRootList(m.p, flds, sql, rootIDs)
}

// CreateRoot TBD
func (m *Manager) CreateRoot(root *dto.Root) error {
	flds := dto.RootAllFields.Del(dto.RootFieldID)

	// sql := fmt.Sprintf(`INSERT INTO "root" (%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
	// 	flds.JoinedNames(),
	// 	flds.Placeholders(),
	// 	dto.RootFieldID.Name(),
	// )

	sql := fmt.Sprintf(`INSERT INTO "root" DEFAULT VALUES RETURNING %[1]s`,
		dto.RootFieldID.Name(),
	)

	return dto.QueryRootRow(
		m.p,
		dto.RootFieldsList{dto.RootFieldID},
		sql,
		(*dto.Root)(root).FieldsValues(flds)...,
	).ScanTo(root)
}
