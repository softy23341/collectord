package dto

import "time"

//go:generate dbgen -type UserSession

// UserSession TBD
type UserSession struct {
	ID           int64      `db:"id"`
	UserID       int64      `db:"user_id"`
	AuthToken    string     `db:"auth_token"`
	CreationTime *time.Time `db:"creation_time"`
}

// UserSessionList TBD
type UserSessionList []*UserSession
