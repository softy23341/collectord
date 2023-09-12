package dto

import "time"

//go:generate dbgen -type UserBan

// Chat TBD
type UserBan struct {
	ID            int64      `db:"id"`
	CreatorUserID int64      `db:"creator_user_id"`
	UserID        int64      `db:"user_id"`
	CreationTime  *time.Time `db:"creation_time"`
}

// UserBanList TBD
type UserBanList []*UserBan

// IDs TBD
func (list UserBanList) IDs() []int64 {
	ids := make([]int64, 0, len(list))
	for _, c := range list {
		ids = append(ids, c.ID)
	}
	return ids
}

// IDs TBD
func (list UserBanList) UserIDs() []int64 {
	ids := make([]int64, 0, len(list))
	for _, c := range list {
		ids = append(ids, c.UserID)
	}
	return ids
}
