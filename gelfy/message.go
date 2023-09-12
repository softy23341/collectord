package gelfy

import (
	"encoding/json"
	"strings"
)

// ExtraFields represents GELF additional fields
type ExtraFields map[string]interface{}

// Message describes GELF message.
type Message struct {
	Version      string
	Host         string
	Level        int
	ShortMessage string
	FullMessage  string
	Timestamp    float64
	Extra        ExtraFields
}

// MarshalJSON implements json.Marshaller interface for Message.
func (m *Message) MarshalJSON() ([]byte, error) {
	d := map[string]interface{}{
		"host":          m.Host,
		"short_message": m.ShortMessage,
	}

	if m.Version == "" {
		d["version"] = "1.1"
	}

	if m.FullMessage != "" {
		d["full_message"] = m.FullMessage
	}

	if m.Timestamp != 0.0 {
		d["timestamp"] = m.Timestamp
	}

	if m.Level != 0 {
		d["level"] = m.Level
	}

	for k, v := range m.Extra {
		k := strings.Trim(k, " ")
		if k == "" {
			continue
		}

		if !strings.HasPrefix(k, "_") {
			k = "_" + k
		}

		if k == "_id" {
			k = "__id"
		}

		d[k] = v
	}

	return json.Marshal(d)
}
