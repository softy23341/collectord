package loggy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/inconshreveable/log15"

	"git.softndit.com/collector/backend/gelfy"
)

// GelfFormat implements log15.Format interface, formats log records in Graylog GELF format.
// extra represents a map which will be added to all messages as GELF additional fields.
func GelfFormat(extra gelfy.ExtraFields) log15.Format {
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}

	return &gelfFormat{extra: extra, host: host}
}

type gelfFormat struct {
	extra gelfy.ExtraFields
	host  string
}

func (f *gelfFormat) Format(r *log15.Record) []byte {
	m := gelfy.Message{
		Version:      "1.1",
		Host:         f.host,
		ShortMessage: r.Msg,
		FullMessage:  f.fullMessage(r.Msg, r.Ctx),
		Timestamp:    float64(r.Time.UnixNano()) / float64(time.Second),
		Level:        f.convertLevel(r.Lvl),
		Extra:        f.extra,
	}

	b, err := json.Marshal(&m)
	if err != nil {
		m.ShortMessage = "(logger error) " + m.ShortMessage
		m.FullMessage = ""
		m.Extra = gelfy.ExtraFields{"logerr": err.Error()}
		b, _ = json.Marshal(&m)
	}

	return b
}

func (f *gelfFormat) fullMessage(msg string, ctx []interface{}) string {
	buf := &bytes.Buffer{}

	fmt.Fprintf(buf, "msg=%s", msg)

	for i := 0; i < len(ctx); i += 2 {
		fmt.Fprintf(buf, " %v=%v", ctx[i], ctx[i+1])
	}

	return buf.String()
}

func (f *gelfFormat) convertLevel(l log15.Lvl) int {
	switch l {
	case log15.LvlDebug:
		return gelfy.LvlDebug
	case log15.LvlInfo:
		return gelfy.LvlInfo
	case log15.LvlWarn:
		return gelfy.LvlWarn
	case log15.LvlError:
		return gelfy.LvlError
	case log15.LvlCrit:
		return gelfy.LvlCrit
	default:
		return gelfy.LvlNotice
	}
}
