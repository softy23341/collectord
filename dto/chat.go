package dto

import "time"

//go:generate dbgen -type Chat

// Chat TBD
type Chat struct {
	ID                int64      `db:"id"`
	CreatorUserID     int64      `db:"creator_user_id"`
	UserUniqID        int64      `db:"user_uniq_id"`
	AdminUserID       int64      `db:"admin_user_id"`
	Name              string     `db:"name"`
	AvatarMediaID     *int64     `db:"avatar_media_id"`
	LastReadMessageID int64      `db:"last_read_message_id"`
	CreationTime      *time.Time `db:"creation_time"`
}

func (c Chat) hasMedia() bool {
	return c.AvatarMediaID != nil
}

func (c Chat) mediaID() *int64 {
	return c.AvatarMediaID
}

// ChatList TBD
type ChatList []*Chat

// MediasIDs TBD
func (chatList ChatList) MediasIDs() []int64 {
	medias := make([]int64, 0, len(chatList))
	for _, c := range chatList {
		if c.hasMedia() {
			medias = append(medias, *c.mediaID())
		}
	}
	return medias
}

// IDs TBD
func (chatList ChatList) IDs() []int64 {
	ids := make([]int64, 0, len(chatList))
	for _, c := range chatList {
		ids = append(ids, c.ID)
	}
	return ids
}
