package dalpg

import (
	"fmt"

	"git.softndit.com/collector/backend/dto"
)

// GetObjectStatuses TBD
func (m *Manager) GetObjectStatuses() (dto.ObjectStatusList, error) {
	flds := dto.ObjectStatusAllFields

	sql := fmt.Sprintf(`
                SELECT %s
                FROM object_status
                ORDER BY name ASC
        `, flds.JoinedNames())

	return dto.ScanObjectStatusList(m.p, flds, sql)
}

// GetCurrentObjectsStatusesRefs TBD
func (m *Manager) GetCurrentObjectsStatusesRefs(objectsIDs []int64) (dto.ObjectStatusRefList, error) {
	flds := dto.ObjectStatusRefAllFields

	sql := fmt.Sprintf(`
                SELECT %s
                FROM object_status_ref AS osr
                JOIN (
                  SELECT osr.object_id, MAX(osr.start_date) AS start_date
                  FROM object_status_ref AS osr
                  WHERE osr.start_date < current_timestamp
                  GROUP BY osr.object_id
                ) AS closest_ref ON closest_ref.object_id = osr.object_id
                                 AND closest_ref.start_date = osr.start_date
                WHERE osr.object_id = any($1::bigint[]);
        `, flds.JoinedNamesWithAlias("osr"))

	return dto.ScanObjectStatusRefList(m.p, flds, sql, objectsIDs)
}

// GetObjectStatusByIDs TBD
func (m *Manager) GetObjectStatusByIDs(IDs []int64) (dto.ObjectStatusList, error) {
	flds := dto.ObjectStatusAllFields

	sql := fmt.Sprintf(`
                SELECT %s
                FROM object_status
                WHERE id = any($1::bigint[])
        `, flds.JoinedNames())

	return dto.ScanObjectStatusList(m.p, flds, sql, IDs)
}

// CreateObjectStatusRef TBD
func (m *Manager) CreateObjectStatusRef(ref *dto.ObjectStatusRef) error {
	retFlds := dto.ObjectStatusRefFieldsList{
		dto.ObjectStatusRefFieldID,
		dto.ObjectStatusRefFieldStartDate,
		dto.ObjectStatusRefFieldCreationTime,
	}
	insFlds := dto.ObjectStatusRefAllFields.Del(retFlds...)

	sql := fmt.Sprintf(`INSERT INTO object_status_ref (%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		insFlds.JoinedNames(), insFlds.Placeholders(), retFlds.JoinedNames())

	return dto.QueryObjectStatusRefRow(m.p, retFlds, sql, ref.FieldsValues(insFlds)...).ScanTo(ref)
}

// CreateObjectStatus TBD
func (m *Manager) CreateObjectStatus(objectStatus *dto.ObjectStatus) error {
	retFlds := dto.ObjectStatusFieldsList{
		dto.ObjectStatusFieldID,
	}
	insFlds := dto.ObjectStatusAllFields.Del(retFlds...)

	sql := fmt.Sprintf(`INSERT INTO object_status (%[1]s) VALUES(%[2]s) RETURNING %[3]s`,
		insFlds.JoinedNames(), insFlds.Placeholders(), retFlds.JoinedNames())

	return dto.QueryObjectStatusRow(m.p, retFlds, sql, objectStatus.FieldsValues(insFlds)...).ScanTo(objectStatus)
}
