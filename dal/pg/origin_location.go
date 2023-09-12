package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetObjectsOriginLocationRefs TBD
func (m *Manager) GetObjectsOriginLocationRefs(objectsIDs []int64) (dto.ObjectOriginLocationRefList, error) {
	flds := dto.ObjectOriginLocationRefAllFields

	sql := fmt.Sprintf(`
	        SELECT %s
                FROM object_origin_location_ref
                WHERE object_id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanObjectOriginLocationRefList(m.p, flds, sql, objectsIDs)
}

// GetOriginLocationsByIDs TBD
func (m *Manager) GetOriginLocationsByIDs(originLocationsIDs []int64) (dto.OriginLocationList, error) {
	flds := dto.OriginLocationAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM origin_location
                WHERE id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanOriginLocationList(m.p, flds, sql, originLocationsIDs)

}

// GetOrCreateOriginLocationByNormalName TBD
func (m *Manager) GetOrCreateOriginLocationByNormalName(inOriginLocation *dto.OriginLocation) (*dto.OriginLocation, error) {
	outOriginLocations, err := m.GetOriginLocationsByNormalNames(
		*inOriginLocation.RootID,
		[]string{inOriginLocation.NormalName})

	if err != nil {
		return nil, err
	}
	if len(outOriginLocations) > 0 {
		return outOriginLocations[0], nil
	}
	if err := m.CreateOriginLocation(inOriginLocation); err != nil {
		return nil, err
	}
	return inOriginLocation, nil
}

// GetOriginLocationsByNormalNames TBD
func (m *Manager) GetOriginLocationsByNormalNames(rootID int64, normalNames []string) (dto.OriginLocationList, error) {
	flds := dto.OriginLocationAllFields
	sql := fmt.Sprintf(`
	  SELECT %s
	  FROM origin_location
	  WHERE %s = $1 AND %s = any($2::varchar[])
	`, flds.JoinedNames(), dto.OriginLocationFieldRootID.Name(), dto.OriginLocationFieldNormalName.Name())

	return dto.ScanOriginLocationList(m.p, flds, sql, rootID, normalNames)
}

// CreateOriginLocation TBD
func (m *Manager) CreateOriginLocation(originLocation *dto.OriginLocation) error {
	flds := dto.OriginLocationAllFields.Del(dto.OriginLocationFieldID)

	sql := fmt.Sprintf(`INSERT INTO origin_location(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		flds.JoinedNames(),
		flds.Placeholders(),
		dto.OriginLocationFieldID.Name(),
	)

	return dto.QueryOriginLocationRow(
		m.p,
		dto.OriginLocationFieldsList{dto.OriginLocationFieldID},
		sql,
		(*dto.OriginLocation)(originLocation).FieldsValues(flds)...,
	).ScanTo(originLocation)
}

// CreateObjectOriginLocationsRefs TBD
func (m *TxManager) CreateObjectOriginLocationsRefs(objectID int64, originLocationIDs []int64) error {
	// TODO make batch insert
	for _, originLocationID := range originLocationIDs {
		ref := &dto.ObjectOriginLocationRef{
			ObjectID:         objectID,
			OriginLocationID: originLocationID,
		}
		if err := m.createObjectOriginLocationRef(ref); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) createObjectOriginLocationRef(ref *dto.ObjectOriginLocationRef) error {
	flds := dto.ObjectOriginLocationRefAllFields

	sql := fmt.Sprintf(`INSERT INTO object_origin_location_ref(%[1]s) VALUES(%[2]s)`,
		flds.JoinedNames(),
		flds.Placeholders(),
	)

	_, err := m.p.Exec(sql, (*dto.ObjectOriginLocationRef)(ref).FieldsValues(flds)...)
	return err
}

// DeleteObjectOriginLocationsRefs TBD
func (m *Manager) DeleteObjectOriginLocationsRefs(objectID int64) error {
	sql := fmt.Sprintf(`DELETE FROM object_origin_location_ref WHERE %s = $1`,
		dto.ObjectOriginLocationRefFieldObjectID.Name(),
	)

	_, err := m.p.Exec(sql, objectID)
	return err
}

// UpdateOriginLocation TBD
func (m *Manager) UpdateOriginLocation(originLocation *dto.OriginLocation) error {
	sql := `
          UPDATE origin_location
          SET
            name = $1,
            normal_name = $2
          WHERE id = $3
        `
	_, err := m.p.Exec(sql,
		originLocation.Name,
		originLocation.NormalName,
		originLocation.ID,
	)
	return err
}

// GetOriginLocationsByRootID TBD
func (m *Manager) GetOriginLocationsByRootID(rootID int64) (dto.OriginLocationList, error) {
	flds := dto.OriginLocationAllFields

	sql := fmt.Sprintf(`
	        SELECT %s
                FROM origin_location
                WHERE root_id = $1
                ORDER BY id DESC
        `, flds.JoinedNames())

	return dto.ScanOriginLocationList(m.p, flds, sql, rootID)
}

// DeleteOriginLocations TBD
func (m *Manager) DeleteOriginLocations(originLocationsIDs []int64) error {
	sql := `DELETE FROM origin_location WHERE id = any($1::bigint[])`

	_, err := m.p.Exec(sql, originLocationsIDs)
	return err
}

// CopyOriginLocationToRoot TBD
func (m *Manager) CopyOriginLocationToRoot(rootID int64) error {
	sql := `
          INSERT INTO origin_location (
            name,
            normal_name,
            root_id
          ) SELECT
            name,
            normal_name,
            $1
          FROM origin_location WHERE root_id IS NULL
        `

	_, err := m.p.Exec(sql, rootID)
	return err
}
