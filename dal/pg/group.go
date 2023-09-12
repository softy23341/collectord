package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetGroupsByIDs TBD
func (m *Manager) GetGroupsByIDs(groupsIDs []int64) (dto.GroupList, error) {
	flds := dto.GroupAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM "group"
                WHERE id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanGroupList(m.p, flds, sql, groupsIDs)
}

// GetGroupsByRootID TBD
func (m *Manager) GetGroupsByRootID(rootID int64) (dto.GroupList, error) {
	flds := dto.GroupAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM "group"
                WHERE root_id = $1
                ORDER BY id ASC
        `, flds.JoinedNames())

	return dto.ScanGroupList(m.p, flds, sql, rootID)
}

// GetCollectionsGroupRefs TBD
func (m *Manager) GetCollectionsGroupRefs(collectionsIDs []int64) (dto.CollectionGroupRefList, error) {
	flds := dto.CollectionGroupRefAllFields

	sql := fmt.Sprintf(`
	        SELECT %s
                FROM collection_group_ref
                WHERE collection_id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanCollectionGroupRefList(m.p, flds, sql, collectionsIDs)
}

// GetCollectionsGroupRefsByGroups TBD
func (m *Manager) GetCollectionsGroupRefsByGroups(groupsIDs []int64) (dto.CollectionGroupRefList, error) {
	flds := dto.CollectionGroupRefAllFields

	sql := fmt.Sprintf(`
	        SELECT %s
                FROM collection_group_ref
                WHERE group_id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanCollectionGroupRefList(m.p, flds, sql, groupsIDs)
}

// CreateGroup TBD
func (m *Manager) CreateGroup(group *dto.Group) error {
	retFlds := dto.GroupFieldsList{dto.GroupFieldID}
	insFlds := dto.GroupAllFields.
		Del(retFlds...)

	sql := fmt.Sprintf(`INSERT INTO "group" (%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		insFlds.JoinedNames(), insFlds.Placeholders(), retFlds.JoinedNames())

	return dto.QueryGroupRow(m.p, retFlds, sql, group.FieldsValues(insFlds)...).ScanTo(group)
}

// CreateCollectionGroupsRefs TBD
func (m *Manager) CreateCollectionGroupsRefs(collectionID int64, groupsIDs []int64) error {
	// TODO make batch insert
	for _, groupID := range groupsIDs {
		ref := &dto.CollectionGroupRef{
			CollectionID: collectionID,
			GroupID:      groupID,
		}
		if err := m.createCollectionGroupRef(ref); err != nil {
			return err
		}
	}
	return nil
}

// CreateGroupCollectionsRefs TBD
func (m *TxManager) CreateGroupCollectionsRefs(groupID int64, collectionsIDs []int64) error {
	// TODO make batch insert
	for _, collectionID := range collectionsIDs {
		ref := &dto.CollectionGroupRef{
			CollectionID: collectionID,
			GroupID:      groupID,
		}
		if err := m.createCollectionGroupRef(ref); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) createCollectionGroupRef(ref *dto.CollectionGroupRef) error {
	flds := dto.CollectionGroupRefAllFields

	sql := fmt.Sprintf(`INSERT INTO collection_group_ref(%[1]s) VALUES(%[2]s)`,
		flds.JoinedNames(),
		flds.Placeholders(),
	)

	_, err := m.p.Exec(sql, (*dto.CollectionGroupRef)(ref).FieldsValues(flds)...)
	return err
}

// UpdateGroup TBD
func (m *Manager) UpdateGroup(group *dto.Group) error {
	sql := `
          UPDATE "group"
          SET
            name = $1
          WHERE id = $2
        `

	_, err := m.p.Exec(sql,
		group.Name,
		group.ID,
	)
	return err
}

// DeleteCollectionsGroupRefs TBD
func (m *Manager) DeleteCollectionsGroupRefs(groupID int64) error {
	sql := fmt.Sprintf(`DELETE FROM collection_group_ref WHERE %s = $1`,
		dto.CollectionGroupRefFieldGroupID.Name(),
	)

	_, err := m.p.Exec(sql, groupID)
	return err
}

// DeleteGroups TBD
func (m *Manager) DeleteGroups(groupsIDs []int64) error {
	sql := `DELETE FROM "group" WHERE id = any($1::bigint[])`

	_, err := m.p.Exec(sql, groupsIDs)
	return err
}

// DeleteCollectionGroupsRefs TBD
func (m *Manager) DeleteCollectionGroupsRefs(collectionID int64) error {
	sql := `DELETE FROM collection_group_ref WHERE collection_id = $1`

	_, err := m.p.Exec(sql, collectionID)
	return err
}

// GetGroupByRootIDAndName TBD
func (m *Manager) GetGroupByRootIDAndName(rootID int64, name string) (*dto.Group, error) {
	flds := dto.GroupAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM "group"
                WHERE root_id = $1 AND name ILIKE $2
                ORDER BY id ASC
        `, flds.JoinedNames())

	return dto.ScanGroup(m.p, flds, sql, rootID, name)
}
