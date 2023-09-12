package dto

import "time"

//go:generate dbgen -type Dialog

// Dialog TBD
type Dialog struct {
	ID             int64     `db:"id"`
	LeastUserID    int64     `db:"least_user_id"`
	GreatestUserID int64     `db:"greatest_user_id"`
	CreationTime   time.Time `db:"creation_time"`
}
