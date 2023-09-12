package dalpg

import (
	"fmt"

	"strings"

	"git.softndit.com/collector/backend/dto"
	"github.com/jackc/pgx"
)

// GetCollectionIDByUserUniqID TBD
func (m *Manager) GetCollectionIDByUserUniqID(userID, uniqID int64) (*int64, error) {
	sql := `
                SELECT id
                FROM collection
                WHERE user_id = $1
                AND user_uniq_id = $2`

	var objectID int64
	err := m.p.QueryRow(sql, userID, uniqID).Scan(&objectID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &objectID, nil
}

// CreateCollection TBD
func (m *Manager) CreateCollection(collection *dto.Collection) error {
	flds := dto.CollectionAllFields.Del(dto.CollectionFieldID)

	sql := fmt.Sprintf(`INSERT INTO collection(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		flds.JoinedNames(),
		flds.Placeholders(),
		dto.CollectionFieldID.Name(),
	)

	return dto.QueryCollectionRow(
		m.p,
		dto.CollectionFieldsList{dto.CollectionFieldID},
		sql,
		(*dto.Collection)(collection).FieldsValues(flds)...,
	).ScanTo(collection)
}

// GetCollectionsByIDs TBD
func (m *Manager) GetCollectionsByIDs(collectionsIDs []int64) (dto.CollectionList, error) {
	flds := dto.CollectionAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM collection
                WHERE id = any($1::bigint[])
                ORDER BY id DESC
        `, flds.JoinedNames())

	return dto.ScanCollectionList(m.p, flds, sql, collectionsIDs)
}

// GetCollectionsByIDsForUpdate TBD
func (m *Manager) GetCollectionsByIDsForUpdate(collectionsIDs []int64) (dto.CollectionList, error) {
	flds := dto.CollectionAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM collection
                WHERE id = any($1::bigint[])
                ORDER BY id DESC
                FOR UPDATE
        `, flds.JoinedNames())

	return dto.ScanCollectionList(m.p, flds, sql, collectionsIDs)
}

// GetTypedCollection TBD
func (m *Manager) GetTypedCollection(rootID int64, typo dto.CollectionTypo) (*dto.Collection, error) {
	flds := dto.CollectionAllFields

	sql := fmt.Sprintf(`
          SELECT %s
          FROM collection
          WHERE root_id = $1
          AND typo = $2
        `, flds.JoinedNames())

	return dto.ScanCollection(m.p, flds, sql, rootID, int16(typo))
}

// UpdateCollection TBD
func (m *Manager) UpdateCollection(collection *dto.Collection) error {
	sql := `
          UPDATE collection
          SET
            name = $1,
            description = $2,
            image_media_id = $3,
			public = $4,
			is_anonymous = $5
          WHERE id = $6
        `

	_, err := m.p.Exec(sql,
		collection.Name,
		collection.Description,
		collection.ImageMediaID,
		collection.Public,
		collection.IsAnonymous,
		collection.ID,
	)
	return err
}

// GetCollectionsByRootID TBD
func (m *Manager) GetCollectionsByRootID(rootID int64) (dto.CollectionList, error) {
	flds := dto.CollectionAllFields

	sql := fmt.Sprintf(`
	        SELECT %s
                FROM collection
                WHERE root_id = $1
                ORDER BY typo DESC, id ASC
        `, flds.JoinedNames())

	return dto.ScanCollectionList(m.p, flds, sql, rootID)
}

// GetPublicCollections TBD
func (m *Manager) GetPublicCollections(query *string, rootID *int64, paginator *dto.PagePaginator) (int64, dto.CollectionList, error) {
	flds := dto.CollectionAllFields

	var args []interface{}
	queryParts := make([]string, 0)
	totalPlaceholders := 0

	totalPlaceholders++
	queryParts = append(queryParts, fmt.Sprintf(" c.public = $%d ", totalPlaceholders))
	args = append(args, true)

	if query != nil {
		totalPlaceholders++
		queryParts = append(queryParts, fmt.Sprintf(" (c.name ILIKE $%d OR c.description ILIKE $%d) ", totalPlaceholders, totalPlaceholders))
		args = append(args, *query+"%")
	}

	if rootID != nil {
		totalPlaceholders++
		queryParts = append(queryParts, fmt.Sprintf(" c.root_id = $%d ", totalPlaceholders))
		args = append(args, *rootID)
	}

	queryPart := strings.Join(queryParts, "AND")

	var cnt int64
	err := m.p.QueryRow(fmt.Sprintf(`SELECT COUNT(c) FROM "collection" AS c WHERE %s`, queryPart), args...).Scan(&cnt)
	if err != nil {
		return 0, nil, err
	}

	sql := fmt.Sprintf(`
          SELECT %s
          FROM "collection" AS c
          WHERE %s
          ORDER BY c.id
          LIMIT $%d OFFSET $%d ;
        `, flds.JoinedNamesWithAlias("c"), queryPart, totalPlaceholders+1, totalPlaceholders+2)

	args = append(args, paginator.Cnt, paginator.Page*paginator.Cnt)

	list, err := dto.ScanCollectionList(m.p, flds, sql, args...)
	if err != nil {
		return 0, nil, err
	}

	return cnt, list, nil
}

// GetCollectionsByGroupIDs TBD
func (m *Manager) GetCollectionsByGroupIDs(groupIDs []int64) (dto.CollectionList, error) {
	flds := dto.CollectionAllFields
	sql := fmt.Sprintf(`
               SELECT %s
               FROM collection AS c
               JOIN (
                 SELECT c.id
                 FROM collection AS c
                 JOIN collection_group_ref AS cr
                   ON cr.collection_id = c.id
                 JOIN "group" AS g
                   ON cr.group_id = g.id
                 WHERE g.id = any($1::bigint[])
                 GROUP BY c.id
               ) AS uniq_c
               ON uniq_c.id = c.id
               ORDER BY c.id DESC
        `, flds.JoinedNamesWithAlias("c"))

	return dto.ScanCollectionList(m.p, flds, sql, groupIDs)
}

// GetPublicCollections TBD
func (m *Manager) SearchPublicCollections(query string, paginator *dto.PagePaginator) (dto.CollectionList, error) {
	flds := dto.CollectionAllFields
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM collection c
                WHERE 
					public IS TRUE AND
					(c.name ILIKE $1 OR c.description ILIKE $1) 
                ORDER BY typo DESC, id ASC
				LIMIT $2 OFFSET $3
        `, flds.JoinedNamesWithAlias("c"))
	sqlParams := []interface{}{query + "%", paginator.Cnt, paginator.Page * paginator.Cnt}
	return dto.ScanCollectionList(m.p, flds, sql, sqlParams...)
}

// GetCollectionsByGroupIDs TBD
func (m *Manager) GetCollectionsByObjectIDs(objectIDs []int64) (dto.CollectionList, error) {
	flds := dto.CollectionAllFields
	sql := fmt.Sprintf(`
		SELECT %s
		FROM collection AS c
        JOIN object AS obj
		ON obj.collection_id = c.id
        WHERE obj.id = any($1::bigint[])
        GROUP BY c.id ORDER BY c.id DESC
	`, flds.JoinedNamesWithAlias("c"))
	return dto.ScanCollectionList(m.p, flds, sql, objectIDs)
}

// GetObjectsCnt TBD
func (m *Manager) GetObjectsCnt(collectionsIDs []int64) (int64, error) {
	c2c, err := m.GetObjectsCntByCollections(collectionsIDs)
	if err != nil {
		return 0, err
	}

	acc := int64(0)
	for _, id := range collectionsIDs {
		acc += c2c[id]
	}
	return acc, nil
}

// GetObjectsCntByCollections TBD
func (m *Manager) GetObjectsCntByCollections(collectionsIDs []int64) (map[int64]int64, error) {
	sql := `
          SELECT collection_id, COUNT(*) AS cnt
          FROM object
          WHERE collection_id = any($1::bigint[])
          GROUP BY collection_id
        `

	rows, err := m.p.Query(sql, collectionsIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		collectionID int64
		cnt          int64
	)
	collection2cnt := make(map[int64]int64)
	for rows.Next() {
		if err := rows.Scan(&collectionID, &cnt); err != nil {
			return nil, err
		}
		collection2cnt[collectionID] = cnt
	}

	if rows.Err() != nil {
		return nil, err
	}

	return collection2cnt, nil
}

// DeleteCollections TBD
func (m *Manager) DeleteCollections(collectionsIDs []int64) error {
	sql := "DELETE FROM collection WHERE id = any($1::bigint[])"

	_, err := m.p.Exec(sql, collectionsIDs)
	return err
}

// DeleteCollectionsGroupRefsByGroupAndCollections TBD
func (m *Manager) DeleteCollectionsGroupRefsByGroupAndCollections(collectionsIDs []int64, groupID int64) error {
	sql := "DELETE FROM collection_group_ref WHERE collection_id = any($1::bigint[]) AND group_id = $2"

	_, err := m.p.Exec(sql, collectionsIDs, groupID)
	return err
}

// GetCollectionsByUserIDsWithCustomFields TBD
func (m *Manager) GetCollectionsByUserIDsWithCustomFields(userIDs []int64, flds dto.CollectionFieldsList) (dto.CollectionList, error) {
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM collection
                WHERE user_id = any($1::bigint[])
                ORDER BY id DESC
        `, flds.JoinedNames())
	return dto.ScanCollectionList(m.p, flds, sql, userIDs)
}

// GetCollectionsByIDsWithCustomFields TBD
func (m *Manager) GetCollectionsByIDsWithCustomFields(collectionsIDs []int64, flds dto.CollectionFieldsList) (dto.CollectionList, error) {
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM collection
                WHERE id = any($1::bigint[])
                ORDER BY id DESC
        `, flds.JoinedNames())
	return dto.ScanCollectionList(m.p, flds, sql, collectionsIDs)
}
