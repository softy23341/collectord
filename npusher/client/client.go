package client

import "git.softndit.com/collector/backend/npusher"

// Client TBD
type Client interface {
	SendPush(token string, sandbox bool, n npusher.Notification) error
}
