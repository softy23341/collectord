package dalpg

import (
	"fmt"

	"github.com/jackc/pgx"

	"git.softndit.com/collector/backend/dto"
)

// GetObjectsByIDs TBD
func (m *Manager) GetObjectsByIDs(objectIDs []int64) (dto.ObjectList, error) {
	flds := dto.ObjectAllFields
	sql := fmt.Sprintf(`SELECT %s FROM object WHERE %s = any($1::bigint[])`,
		flds.JoinedNames(), dto.ObjectFieldID.Name())

	return dto.ScanObjectList(m.p, flds, sql, objectIDs)
}

// GetObjectsPreviewByIDs TBD
func (m *Manager) GetObjectsPreviewByIDs(objectIDs []int64) (dto.ObjectPreviewList, error) {
	flds := dto.ObjectPreviewAllFields
	sql := fmt.Sprintf(`
		SELECT o.id, o.collection_id, o.name, c.root_id 
		FROM object AS o
		LEFT JOIN collection AS c 
		  ON o.collection_id = c.id  
		WHERE o.%s = any($1::bigint[])`, dto.ObjectFieldID.Name())

	return dto.ScanObjectPreviewList(m.p, flds, sql, objectIDs)
}

// GetObjectIDByUserUniqID TBD
func (m *Manager) GetObjectIDByUserUniqID(userID, uniqID int64) (*int64, error) {
	sqlQuery := `
                SELECT id
                FROM object
                WHERE user_id = $1
                AND user_uniq_id = $2`

	var objectID int64
	err := m.p.QueryRow(sqlQuery, userID, uniqID).Scan(&objectID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &objectID, nil
}

// CreateObject TBD
func (m *Manager) CreateObject(object *dto.Object) error {
	flds := dto.ObjectAllFields.Del(dto.ObjectFieldID)

	sql := fmt.Sprintf(`INSERT INTO object(%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		flds.JoinedNames(),
		flds.Placeholders(),
		dto.ObjectFieldID.Name(),
	)

	return dto.QueryObjectRow(
		m.p,
		dto.ObjectFieldsList{dto.ObjectFieldID},
		sql,
		(*dto.Object)(object).FieldsValues(flds)...,
	).ScanTo(object)
}

// UpdateObject TBD
func (m *Manager) UpdateObject(object *dto.Object) error {
	sql := `
	  UPDATE object
          SET
            name = $1
          , production_date_interval_id = $2
          , production_date_interval_from = $3
          , production_date_interval_to = $4
          , description = $5
          , purchase_date = $6
          , purchase_price = $7
          , purchase_price_currency_id = $8
          , root_id_number = $9
          , collection_id = $10
          , provenance = $11
          , update_time = $12
          WHERE id = $13
        `

	_, err := m.p.Exec(sql,
		object.Name,
		object.ProductionDateIntervalID,
		object.ProductionDateIntervalFrom,
		object.ProductionDateIntervalTo,
		object.Description,
		object.PurchaseDate,
		object.PurchasePrice,
		object.PurchaseCurrencyID,
		object.RootIDNumber,
		object.CollectionID,
		object.Provenance,
		object.UpdateTime,
		object.ID,
	)
	return err
}

// GetCollectionsObjectPreviews TBD
func (m *Manager) GetCollectionsObjectPreviews(collectionsIDs []int64, orders *dto.ObjectOrders, paginator *dto.PagePaginator) (dto.ObjectPreviewList, error) {
	flds := dto.ObjectPreviewAllFields

	sql := ""
	if orders.ActorName == 0 {
		orderByClosure := "id DESC"

		// creation time
		if orders.CreationTime != 0 {
			if orders.CreationTime > 0 {
				orderByClosure = "id ASC"
			} else {
				orderByClosure = "id DESC"
			}
		}

		// object name
		if orders.Name != 0 {
			if orders.Name > 0 {
				orderByClosure = "name ASC"
			} else {
				orderByClosure = "name DESC"
			}
		}

		// update time
		if orders.UpdateTime != 0 {
			if orders.UpdateTime > 0 {
				orderByClosure = "update_time ASC"
			} else {
				orderByClosure = "update_time DESC"
			}
		}

		// sql
		sql = fmt.Sprintf(`
	            SELECT o.id, o.collection_id, o.name, c.root_id
                    FROM object AS o
					LEFT JOIN collection AS c 
					  ON o.collection_id = c.id
                    WHERE o.collection_id = any($1::bigint[])
                    ORDER BY %s, o.id
                    LIMIT $2 OFFSET $3
                `, orderByClosure)
	} else {
		orderBy := "ASC"
		if orders.ActorName < 0 {
			orderBy = "DESC"
		}
		sql = fmt.Sprintf(`
				  SELECT o.id, o.collection_id, o.name, c.root_id
                    FROM object AS o
					LEFT JOIN collection AS c 
					  ON o.collection_id = c.id
                    LEFT JOIN object_actor_ref AS oar
                      ON oar.object_id = o.id
                    LEFT JOIN actor AS a
                      ON a.id = oar.actor_id
                    WHERE (
                      oar.id IN (
                        SELECT MIN(oar_filter.id) AS id
                        FROM object_actor_ref AS oar_filter
                        GROUP BY oar_filter.object_id
                      ) OR oar.id IS NULL
                    ) AND (
                      o.collection_id = any($1::bigint[])
                    )
                    ORDER BY oar.id IS NULL, a.name %s, o.id
                    LIMIT $2 OFFSET $3
                    `, orderBy)
	}

	m.log.Debug("run sql", "sql", sql)

	sqlParams := []interface{}{collectionsIDs, paginator.Cnt, paginator.Page * paginator.Cnt}
	return dto.ScanObjectPreviewList(m.p, flds, sql, sqlParams...)
}

// GetObjects TBD
func (m *Manager) GetObjects(paginator *dto.PagePaginator) (dto.ObjectList, error) {
	flds := dto.ObjectAllFields
	sql := fmt.Sprintf(`SELECT %s
                            FROM object
                            ORDER BY id ASC
                            LIMIT $1 OFFSET $2
                           `, flds.JoinedNames())

	return dto.ScanObjectList(m.p, flds, sql, paginator.Limit(), paginator.Offset())
}

// DeleteObjectsByIDs TBD
func (m *Manager) DeleteObjectsByIDs(objectIDs []int64) error {
	sql := "DELETE FROM object WHERE id = any($1::bigint[])"

	_, err := m.p.Exec(sql, objectIDs)
	return err
}

// ChangeObjectsCollection TBD
func (m *Manager) ChangeObjectsCollection(fromID, toID int64) error {
	sql := "UPDATE object SET collection_id = $1 WHERE collection_id = $2"

	_, err := m.p.Exec(sql, toID, fromID)
	return err
}

// DeleteObjectsByCollectionsIDs TBD
func (m *Manager) DeleteObjectsByCollectionsIDs(collectionsIDs []int64) error {
	sql := "DELETE FROM object WHERE collection_id = any($1::bigint[])"

	_, err := m.p.Exec(sql, collectionsIDs)
	return err
}

// ChangeObjectsCollectionByIDs TBD
func (m *Manager) ChangeObjectsCollectionByIDs(objectsIDs []int64, collectionID int64) error {
	sql := `UPDATE object SET collection_id = $1 WHERE id = any($2::bigint[])`

	_, err := m.p.Exec(sql, collectionID, objectsIDs)
	return err
}

// GetObjectsByUserIDsWithCustomFields TBD
func (m *Manager) GetObjectsByUserIDsWithCustomFields(userIDs []int64, flds dto.ObjectFieldsList) (dto.ObjectList, error) {
	sql := fmt.Sprintf(`
	        SELECT %s
                FROM object
                WHERE user_id = any($1::bigint[])
                ORDER BY id DESC
        `, flds.JoinedNames())
	return dto.ScanObjectList(m.p, flds, sql, userIDs)
}
