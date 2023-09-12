package npusher

import "encoding/json"

// NotificationTask TBD
type NotificationTask struct {
	Token          string          `json:"token"`
	Type           string          `json:"type"`
	Sandbox        bool            `json:"sandbox"`
	Payload        json.RawMessage `json:"payload"`
	DefaultMessage string          `json:"default_message,omitempty"`
}
