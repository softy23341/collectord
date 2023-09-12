package npusher

import "encoding/json"

// GCMNotification TBD
type GCMNotification struct {
	Notification *NotificationPayload   `json:"notification,omitempty"`
	Data         map[string]interface{} `json:"data,omitempty"`

	CollapseKey           *string `json:"collapse_key,omitempty"`
	Priority              *string `json:"priority,omitempty"`
	TimeToLive            *int64  `json:"time_to_live,omitempty"`
	RestrictedPackageName *string `json:"restricted_package_name,omitempty"`
	DryRun                *bool   `json:"dry_run,omitempty"`
}

// NotificationPayload TBD
type NotificationPayload struct {
	Title        string   `json:"title"`
	Body         *string  `json:"body,omitempty"`
	Icon         *string  `json:"icon,omitempty"`
	Sound        *string  `json:"sound,omitempty"`
	Tag          *string  `json:"tag,omitempty"`
	Color        *string  `json:"color,omitempty"`
	ClickAction  *string  `json:"click_action,omitempty"`
	BodyLocKey   *string  `json:"body_loc_key,omitempty"`
	BodyLocArgs  []string `json:"body_loc_args,omitempty"`
	TitleLocKey  *string  `json:"title_loc_key,omitempty"`
	TitleLocArgs []string `json:"title_loc_args,omitempty"`
}

// Type TBD
func (n GCMNotification) Type() string {
	return "gcm"
}

// PayloadJSON TBD
func (n GCMNotification) PayloadJSON() ([]byte, error) {
	return json.Marshal(n)
}

// DefaultMessage TBD
func (n GCMNotification) DefaultMessage() string {
	if notification := n.Notification; notification != nil {
		return notification.Title
	}
	return ""
}
