package dto

import "time"

//go:generate dbgen -type ObjectStatusRef

// ObjectStatusRef TBD
type ObjectStatusRef struct {
	ID               int64      `db:"id"`
	ObjectID         int64      `db:"object_id"`
	ObjectStatusID   int64      `db:"object_status_id"`
	StartDate        time.Time  `db:"start_date"`
	NotificationDate *time.Time `db:"notification_date"`
	CreationTime     time.Time  `db:"creation_time"`
	Description      string     `db:"description"`
}

// ObjectStatusRefList TBD
type ObjectStatusRefList []*ObjectStatusRef

// ObjectToOneStatusMap TBD
func (osr ObjectStatusRefList) ObjectToOneStatusMap() map[int64]*ObjectStatusRef {
	objectIDStatusRef := make(map[int64]*ObjectStatusRef)
	for _, ref := range osr {
		objectIDStatusRef[ref.ObjectID] = ref
	}
	return objectIDStatusRef
}
