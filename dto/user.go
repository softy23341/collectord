package dto

import (
	"golang.org/x/crypto/bcrypt"
)

//go:generate dbgen -type User

// User dto
type (
	User struct {
		ID                int64  `db:"id"`
		FirstName         string `db:"first_name"`
		LastName          string `db:"last_name"`
		Email             string `db:"email"`
		EmailVerified     bool   `db:"email_verified"`
		EncryptedPassword string `db:"encrypted_password"`
		AvatarMediaID     *int64 `db:"avatar_media_id"`
		Description       string `db:"description"`
		SystemUser        bool   `db:"system_user"`
		Locale            string `db:"locale"`
		IsAnonymous       bool   `db:"is_anonymous"`

		LastEventSeqNo       int64 `db:"last_event_seq_no"`
		NUnreadMessages      int32 `db:"n_unread_messages"`
		NUnreadNotifications int32 `db:"n_unread_notifications"`

		Tags       []string `db:"tags"`
		Speciality *string  `db:"speciality"`
	}

	// UserList TBD
	UserList []*User
)

// IsEqualEPass TBD
func (u *User) IsEqualEPass(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

// FullName TBD
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IDs TBD
func (uu UserList) IDs() []int64 {
	o := make([]int64, len(uu))
	for i := range uu {
		o[i] = uu[i].ID
	}
	return o
}

// GetAvatarsMediaIDs TBD
func (uu UserList) GetAvatarsMediaIDs() []int64 {
	ids := make([]int64, 0, len(uu))
	for _, user := range uu {
		if user.AvatarMediaID != nil {
			ids = append(ids, *user.AvatarMediaID)
		}
	}
	return ids
}

// IDToUser TBD
func (uu UserList) IDToUser() map[int64]*User {
	id2user := make(map[int64]*User, 0)
	for _, user := range uu {
		id2user[user.ID] = user
	}
	return id2user
}
