package npusher

import "encoding/json"

// APNSNotification TBD
type APNSNotification struct {
	Alert            string `json:"alert,omitempty"`
	Badge            *int   `json:"badge,omitempty"`
	Sound            string `json:"sound,omitempty"`
	ContentAvailable int    `json:"content-available,omitempty"`

	Custom map[string]interface{} `json:"-"`
}

// Type TBD
func (n APNSNotification) Type() string {
	return "apns"
}

// PayloadJSON TBD
func (n APNSNotification) PayloadJSON() ([]byte, error) {
	apsPayload, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}

	// add custom fields
	var m = make(map[string]interface{}, len(n.Custom)+1)
	for k, v := range n.Custom {
		m[k] = v
	}

	apsPayloadJSON := json.RawMessage(apsPayload)
	m["aps"] = &apsPayloadJSON

	return json.Marshal(m)
}

// DefaultMessage TBD
func (n APNSNotification) DefaultMessage() string {
	return n.Alert
}
