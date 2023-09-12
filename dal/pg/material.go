package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetMaterialsByIDs TBD
func (m *Manager) GetMaterialsByIDs(objectIDs []int64) (dto.MaterialList, error) {
	flds := dto.MaterialAllFields

	sql := fmt.Sprintf(`
          SELECT %s
          FROM material
          WHERE id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanMaterialList(m.p, flds, sql, objectIDs)
}

// GetMaterialsByObjectID TBD
func (m *Manager) GetMaterialsByObjectID(objectID int64) (dto.MaterialList, error) {
	flds := dto.MaterialAllFields

	sql := fmt.Sprintf(`
	  SELECT %s
	  FROM material AS m
	  JOIN object_material_ref AS omr ON omr.material_id = m.id
	  WHERE omr.object_id = $1
	`, flds.JoinedNamesWithAlias("m"))

	return dto.ScanMaterialList(m.p, flds, sql, objectID)
}

// GetOrCreateMaterialByNormalName TBD
func (m *Manager) GetOrCreateMaterialByNormalName(inMaterial *dto.Material) (*dto.Material, error) {
	outMaterials, err := m.GetMaterialsByNormalNames(
		*inMaterial.RootID,
		[]string{inMaterial.NormalName})

	if err != nil {
		return nil, err
	}
	if len(outMaterials) > 0 {
		return outMaterials[0], nil
	}
	if err := m.CreateMaterial(inMaterial); err != nil {
		return nil, err
	}
	return inMaterial, nil
}

// GetMaterialsByNormalNames TBD
func (m *Manager) GetMaterialsByNormalNames(rootID int64, normalNames []string) (dto.MaterialList, error) {
	flds := dto.MaterialAllFields
	sql := fmt.Sprintf(`
	  SELECT %s
	  FROM material
	  WHERE %s = $1 AND %s = any($2::varchar[])
	`, flds.JoinedNames(), dto.MaterialFieldRootID.Name(), dto.MaterialFieldNormalName.Name())

	return dto.ScanMaterialList(m.p, flds, sql, rootID, normalNames)
}

// CreateMaterial TBD
func (m *Manager) CreateMaterial(material *dto.Material) error {
	flds := dto.MaterialAllFields.Del(dto.MaterialFieldID)

	sql := fmt.Sprintf(`INSERT INTO material(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		flds.JoinedNames(),
		flds.Placeholders(),
		dto.MaterialFieldID.Name(),
	)

	return dto.QueryMaterialRow(
		m.p,
		dto.MaterialFieldsList{dto.MaterialFieldID},
		sql,
		(*dto.Material)(material).FieldsValues(flds)...,
	).ScanTo(material)
}

// CreateObjectMaterialsRefs TBD
func (m *TxManager) CreateObjectMaterialsRefs(objectID int64, materialIDs []int64) error {
	// TODO make batch insert
	for _, materialID := range materialIDs {
		ref := &dto.ObjectMaterialRef{
			ObjectID:   objectID,
			MaterialID: materialID,
		}
		if err := m.createObjectMaterialRef(ref); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) createObjectMaterialRef(ref *dto.ObjectMaterialRef) error {
	flds := dto.ObjectMaterialRefAllFields

	sql := fmt.Sprintf(`INSERT INTO object_material_ref(%[1]s) VALUES(%[2]s)`,
		flds.JoinedNames(),
		flds.Placeholders(),
	)

	_, err := m.p.Exec(sql, (*dto.ObjectMaterialRef)(ref).FieldsValues(flds)...)
	return err
}

// DeleteObjectMaterialsRefs TBD
func (m *Manager) DeleteObjectMaterialsRefs(objectID int64) error {
	sql := fmt.Sprintf(`DELETE FROM object_material_ref WHERE %s = $1`,
		dto.ObjectMaterialRefFieldObjectID.Name(),
	)

	_, err := m.p.Exec(sql, objectID)
	return err
}

// GetMaterialRefsByObjectsIDs TBD
func (m *Manager) GetMaterialRefsByObjectsIDs(objectIDs []int64) (dto.ObjectMaterialRefList, error) {
	flds := dto.ObjectMaterialRefAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM object_material_ref
                WHERE object_id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanObjectMaterialRefList(m.p, flds, sql, objectIDs)
}

// UpdateMaterial TBD
func (m *Manager) UpdateMaterial(material *dto.Material) error {
	sql := `
          UPDATE material
          SET
            name = $1,
            normal_name = $2
          WHERE id = $3
        `

	_, err := m.p.Exec(sql,
		material.Name,
		material.NormalName,
		material.ID,
	)
	return err
}

// GetMaterialsByRootID TBD
func (m *Manager) GetMaterialsByRootID(rootID int64) (dto.MaterialList, error) {
	flds := dto.MaterialAllFields

	sql := fmt.Sprintf(`
	        SELECT %s
                FROM material
                WHERE root_id = $1
                ORDER BY id DESC
        `, flds.JoinedNames())

	return dto.ScanMaterialList(m.p, flds, sql, rootID)
}

// DeleteMaterials TBD
func (m *Manager) DeleteMaterials(materialsIDs []int64) error {
	sql := `DELETE FROM material WHERE id = any($1::bigint[])`

	_, err := m.p.Exec(sql, materialsIDs)
	return err
}

// CopyMaterialsToRoot TBD
func (m *Manager) CopyMaterialsToRoot(rootID int64) error {
	sql := `
          INSERT INTO material (
            name,
            normal_name,
            root_id
          ) SELECT
            name,
            normal_name,
            $1
          FROM material WHERE root_id IS NULL
        `

	_, err := m.p.Exec(sql, rootID)
	return err
}
