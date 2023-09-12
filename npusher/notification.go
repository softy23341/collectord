package npusher

// Notification TBD
type Notification interface {
	Type() string
	PayloadJSON() ([]byte, error)
	DefaultMessage() string
}
