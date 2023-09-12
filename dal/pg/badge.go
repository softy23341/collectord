package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetObjectsBadgeRefs TBD
func (m *Manager) GetObjectsBadgeRefs(objectsIDs []int64) (dto.ObjectBadgeRefList, error) {
	flds := dto.ObjectBadgeRefAllFields

	sql := fmt.Sprintf(`
	        SELECT %s
                FROM object_badge_ref
                WHERE object_id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanObjectBadgeRefList(m.p, flds, sql, objectsIDs)
}

// GetBadgesByIDs TBD
func (m *Manager) GetBadgesByIDs(badgesIDs []int64) (dto.BadgeList, error) {
	flds := dto.BadgeAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM badge
                WHERE id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanBadgeList(m.p, flds, sql, badgesIDs)

}

// GetOrCreateBadgeByNormalNameAndColor TBD
func (m *Manager) GetOrCreateBadgeByNormalNameAndColor(inBadge *dto.Badge) (*dto.Badge, error) {
	outBadges, err := m.GetBadgesByNormalNamesOrColors(
		*inBadge.RootID,
		[]string{inBadge.NormalName},
		[]string{inBadge.Color},
	)

	if err != nil {
		return nil, err
	}
	if len(outBadges) > 0 {
		return outBadges[0], nil
	}
	if err := m.CreateBadge(inBadge); err != nil {
		return nil, err
	}
	return inBadge, nil
}

// GetBadgesByNormalNamesOrColors TBD
func (m *Manager) GetBadgesByNormalNamesOrColors(rootID int64, normalNames, colors []string) (dto.BadgeList, error) {
	flds := dto.BadgeAllFields
	sql := fmt.Sprintf(`
	  SELECT %s
	  FROM badge
	  WHERE %s = $1 AND (%s = any($2::varchar[]) OR %s = any($3::varchar[]))
	`, flds.JoinedNames(),
		dto.BadgeFieldRootID.Name(),
		dto.BadgeFieldNormalName.Name(),
		dto.BadgeFieldColor.Name())

	return dto.ScanBadgeList(m.p, flds, sql, rootID, normalNames, colors)
}

// CreateBadge TBD
func (m *Manager) CreateBadge(badge *dto.Badge) error {
	flds := dto.BadgeAllFields.Del(dto.BadgeFieldID)

	sql := fmt.Sprintf(`INSERT INTO badge(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		flds.JoinedNames(),
		flds.Placeholders(),
		dto.BadgeFieldID.Name(),
	)

	return dto.QueryBadgeRow(
		m.p,
		dto.BadgeFieldsList{dto.BadgeFieldID},
		sql,
		(*dto.Badge)(badge).FieldsValues(flds)...,
	).ScanTo(badge)
}

// CreateObjectBadgesRefs TBD
func (m *TxManager) CreateObjectBadgesRefs(objectID int64, badgeIDs []int64) error {
	// TODO make batch insert
	for _, badgeID := range badgeIDs {
		ref := &dto.ObjectBadgeRef{
			ObjectID: objectID,
			BadgeID:  badgeID,
		}
		if err := m.createObjectBadgeRef(ref); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) createObjectBadgeRef(ref *dto.ObjectBadgeRef) error {
	flds := dto.ObjectBadgeRefAllFields

	sql := fmt.Sprintf(`INSERT INTO object_badge_ref(%[1]s) VALUES(%[2]s)`,
		flds.JoinedNames(),
		flds.Placeholders(),
	)

	_, err := m.p.Exec(sql, (*dto.ObjectBadgeRef)(ref).FieldsValues(flds)...)
	return err
}

// DeleteObjectBadgesRefs TBD
func (m *Manager) DeleteObjectBadgesRefs(objectID int64) error {
	sql := fmt.Sprintf(`DELETE FROM object_badge_ref WHERE %s = $1`,
		dto.ObjectBadgeRefFieldObjectID.Name(),
	)

	_, err := m.p.Exec(sql, objectID)
	return err
}

// GetBadgesByRootID TBD
func (m *Manager) GetBadgesByRootID(rootID int64) (dto.BadgeList, error) {
	flds := dto.BadgeAllFields

	sql := fmt.Sprintf(`
	        SELECT %s
                FROM badge
                WHERE root_id = $1
                ORDER BY id DESC
        `, flds.JoinedNames())

	return dto.ScanBadgeList(m.p, flds, sql, rootID)
}

// UpdateBadge TBD
func (m *Manager) UpdateBadge(badge *dto.Badge) error {
	sql := `
          UPDATE badge
          SET
            name = $1,
            normal_name = $2,
            color = $3
          WHERE id = $4
        `

	_, err := m.p.Exec(sql,
		badge.Name,
		badge.NormalName,
		badge.Color,
		badge.ID,
	)
	return err
}

// DeleteBadges TBD
func (m *Manager) DeleteBadges(badgesIDs []int64) error {
	sql := `DELETE FROM badge WHERE id = any($1::bigint[])`

	_, err := m.p.Exec(sql, badgesIDs)
	return err
}
