package dalpg

// GetCurrentTransaction TBD
func (m *Manager) GetCurrentTransaction() (int64, error) {
	sql := `SELECT txid_current() AS version`

	rows, err := m.p.Query(sql)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var (
		version int64
	)
	for rows.Next() {
		if err := rows.Scan(&version); err != nil {
			return 0, err
		}
	}

	return version, nil

}
