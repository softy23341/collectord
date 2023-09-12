package dto

//go:generate dbgen -type UserSessionDevice

// UserSessionDeviceType TBD
type UserSessionDeviceType int16

// user devices types
const (
	UserSessionDeviceTypeAndroid = 10
	UserSessionDeviceTypeIOS     = 20
)

// UserSessionDevice TBD
type UserSessionDevice struct {
	ID        int64                 `db:"id"`
	SessionID int64                 `db:"session_id"`
	Typo      UserSessionDeviceType `db:"typo"`
	Token     string                `db:"token"`
	Sandbox   bool                  `db:"push_notification_sandbox"`
}

// UserSessionDeviceList TBD
type UserSessionDeviceList []*UserSessionDevice
