package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetNamedDateIntervalsByIDs TBD
func (m *Manager) GetNamedDateIntervalsByIDs(ids []int64) (dto.NamedDateIntervalList, error) {
	flds := dto.NamedDateIntervalAllFields

	sql := fmt.Sprintf(`
          SELECT %s
          FROM named_date_interval
          WHERE id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanNamedDateIntervalList(m.p, flds, sql, ids)
}

// GetNamedDateIntervalsForRoots TBD
func (m *Manager) GetNamedDateIntervalsForRoots(rootIDs []int64) (dto.NamedDateIntervalList, error) {
	flds := dto.NamedDateIntervalAllFields

	sql := fmt.Sprintf(`
          SELECT %s
          FROM named_date_interval
          WHERE root_id = any($1::bigint[])
          ORDER BY id DESC
        `, flds.JoinedNames())

	return dto.ScanNamedDateIntervalList(m.p, flds, sql, rootIDs)
}

// CreateNamedDateInterval TBD
func (m *Manager) CreateNamedDateInterval(interval *dto.NamedDateInterval) error {
	flds := dto.NamedDateIntervalAllFields.Del(dto.NamedDateIntervalFieldID)

	sql := fmt.Sprintf(`INSERT INTO named_date_interval(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		flds.JoinedNames(),
		flds.Placeholders(),
		dto.NamedDateIntervalFieldID.Name(),
	)

	return dto.QueryNamedDateIntervalRow(
		m.p,
		dto.NamedDateIntervalFieldsList{dto.NamedDateIntervalFieldID},
		sql,
		(*dto.NamedDateInterval)(interval).FieldsValues(flds)...,
	).ScanTo(interval)
}

// UpdateNamedDateInterval TBD
func (m *Manager) UpdateNamedDateInterval(interval *dto.NamedDateInterval) error {
	sql := `
          UPDATE named_date_interval
          SET
            name = $1,
            normal_name = $2,
            production_date_interval_from = $3,
            production_date_interval_to = $4
          WHERE id = $5
        `

	_, err := m.p.Exec(sql,
		interval.Name,
		interval.NormalName,
		interval.ProductionDateIntervalFrom,
		interval.ProductionDateIntervalTo,
		interval.ID,
	)
	return err
}

// GetNamedDayeIntervalsByNormalNames TBD
func (m *Manager) GetNamedDayeIntervalsByNormalNames(rootID int64, names []string) (dto.NamedDateIntervalList, error) {
	flds := dto.NamedDateIntervalAllFields
	sql := fmt.Sprintf(`
	  SELECT %s
	  FROM named_date_interval
	  WHERE %s = $1 AND %s = any($2::varchar[])
	`, flds.JoinedNames(),
		dto.NamedDateIntervalFieldRootID.Name(),
		dto.NamedDateIntervalFieldNormalName.Name())

	return dto.ScanNamedDateIntervalList(m.p, flds, sql, rootID, names)
}

// DeleteNamedDateIntervalsByIDs TBD
func (m *Manager) DeleteNamedDateIntervalsByIDs(ids []int64) error {
	sql := `DELETE FROM named_date_interval WHERE id = any($1::bigint[])`

	_, err := m.p.Exec(sql, ids)
	return err
}

// CopyNamedDateIntervalsToRoot TBD
func (m *Manager) CopyNamedDateIntervalsToRoot(rootID int64) error {
	sql := `
          INSERT INTO named_date_interval (
            production_date_interval_from,
            production_date_interval_to,
            name,
            normal_name,
            root_id
          ) SELECT
            production_date_interval_from,
            production_date_interval_to,
            name,
            normal_name,
            $1
          FROM named_date_interval WHERE root_id IS NULL
        `

	_, err := m.p.Exec(sql, rootID)
	return err
}
