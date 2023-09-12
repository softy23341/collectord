package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetObjectValuations TBD
func (m *Manager) GetObjectValuations(objectsIDs []int64) (dto.ValuationList, error) {
	flds := dto.ValuationAllFields

	sql := fmt.Sprintf(`
SELECT %s
FROM "object_valuations"
WHERE object_id = any($1::bigint[])`,
		flds.JoinedNames(),
	)

	return dto.ScanValuationList(m.p, flds, sql, objectsIDs)
}

// CreateValuation TBD
func (m *Manager) CreateValuation(valuation *dto.Valuation) (*dto.Valuation, error) {
	flds := dto.ValuationAllFields.Del(dto.ValuationFieldID)

	sql := fmt.Sprintf(`
INSERT INTO object_valuations(%[1]s)
VALUES(%[2]s)
RETURNING %[3]s`,
		flds.JoinedNames(),
		flds.Placeholders(),
		dto.ValuationFieldID.Name(),
	)

	err := dto.QueryValuationRow(
		m.p,
		dto.ValuationFieldsList{dto.ValuationFieldID},
		sql,
		(*dto.Valuation)(valuation).FieldsValues(flds)...,
	).ScanTo(valuation)

	if err != nil {
		return nil, err
	}

	return valuation, nil
}

// DeleteValuationsByObjectID TBD
func (m *Manager) DeleteValuationsByObjectID(objectID int64) error {
	sql := "DELETE FROM object_valuations WHERE object_id = $1::bigint"

	_, err := m.p.Exec(sql, objectID)
	return err
}

// GetValuationsByCollectionIDs TBD
func (m *Manager) GetValuationsByCollectionIDs(collectionIDs []int64) (map[int64]int64, error) {
	sql := `
		(SELECT o.collection_id AS id, sum_int8range(array_agg(v.price)) AS sum FROM "object" AS o
			JOIN LATERAL (
				SELECT * FROM object_valuations WHERE object_valuations.object_id = o.id ORDER BY object_valuations.date DESC NULLS LAST LIMIT 1
			) v on true WHERE v.price_currency_id = 1 AND o.collection_id = any($1::bigint[]) GROUP BY o.collection_id)	
		UNION ALL 
		(SELECT o.collection_id AS id, sum(o.purchase_price) AS sum FROM "object" as o 
			LEFT JOIN object_valuations ON o.id = object_valuations.object_id 
		WHERE object_valuations.id IS NULL AND o.purchase_price_currency_id = 1 AND o.collection_id = any($1::bigint[]) GROUP BY o.collection_id)
`
	rows, err := m.p.Query(sql, collectionIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		collectionID int64
		valuation    int64
	)
	collection2valuation := make(map[int64]int64)
	for rows.Next() {
		if err := rows.Scan(&collectionID, &valuation); err != nil {
			return nil, err
		}
		if _, ok := collection2valuation[collectionID]; !ok {
			collection2valuation[collectionID] = 0
		}
		collection2valuation[collectionID] += valuation
	}

	if rows.Err() != nil {
		return nil, err
	}

	return collection2valuation, nil
}
