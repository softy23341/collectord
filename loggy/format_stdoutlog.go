package loggy

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/inconshreveable/log15"
)

const (
	timeFormat  = "01-02 15:04:05.000000"
	floatFormat = 'f'
	errorKey    = "LOG15_ERROR"
)

// StdoutLogHandler TBD
var StdoutLogHandler = log15.StreamHandler(os.Stdout, StdoutLogFormat())

type stdoutLogFormat struct{}

// StdoutLogFormat TBD
func StdoutLogFormat() log15.Format {
	return &stdoutLogFormat{}
}

// DBUG[10-01 16:06:51.000001] msg="something" app=".." ...
func (f *stdoutLogFormat) Format(r *log15.Record) []byte {
	buf := &bytes.Buffer{}

	f.writeHead(r, buf)
	f.writeBody(r, buf)

	return buf.Bytes()
}

func (f *stdoutLogFormat) writeHead(r *log15.Record, buf *bytes.Buffer) {
	fmt.Fprintf(buf, "%s[%s]", f.mapLevelToStr(r.Lvl), f.fmtTime(r.Time))
	buf.WriteByte(' ')
}

func (f *stdoutLogFormat) writeBody(r *log15.Record, buf *bytes.Buffer) {
	msg := []interface{}{"msg", r.Msg}
	f.stdoutLogFmt(buf, append(msg, r.Ctx...))
}

func (f *stdoutLogFormat) stdoutLogFmt(buf *bytes.Buffer, ctx []interface{}) {
	for i := 0; i < len(ctx); i += 2 {
		if i != 0 {
			buf.WriteByte(' ')
		}

		k, ok := ctx[i].(string)
		v := f.fmtValue(ctx[i+1])
		if !ok {
			k, v = errorKey, f.fmtValue(k)
		}

		fmt.Fprintf(buf, "%s=%s", k, v)
	}

	buf.WriteByte('\n')
}

func (f *stdoutLogFormat) fmtShared(value interface{}) (result interface{}) {
	defer func() {
		if err := recover(); err != nil {
			if v := reflect.ValueOf(value); v.Kind() == reflect.Ptr && v.IsNil() {
				result = "nil"
			} else {
				panic(err)
			}
		}
	}()

	switch v := value.(type) {
	case time.Time:
		return v.Format(timeFormat)

	case error:
		return v.Error()

	case fmt.Stringer:
		return v.String()

	default:
		return v
	}
}

// formatValue formats a value for serialization
func (f *stdoutLogFormat) fmtValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	value = f.fmtShared(value)
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)
	case float32:
		return strconv.FormatFloat(float64(v), floatFormat, 3, 64)
	case float64:
		return strconv.FormatFloat(v, floatFormat, 3, 64)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", value)
	case string:
		return f.escapeString(v)
	default:
		return f.escapeString(fmt.Sprintf("%+v", value))
	}
}

func (f *stdoutLogFormat) escapeString(s string) string {
	e := bytes.Buffer{}
	e.WriteByte('"')
	for _, r := range s {
		switch r {
		case '\\', '"':
			e.WriteByte('\\')
			e.WriteByte(byte(r))
		case '\n':
			e.WriteByte('\\')
			e.WriteByte('n')
		case '\r':
			e.WriteByte('\\')
			e.WriteByte('r')
		case '\t':
			e.WriteByte('\\')
			e.WriteByte('t')
		default:
			e.WriteRune(r)
		}
	}
	e.WriteByte('"')
	return e.String()
}

// Returns the name of a Lvl
func (f *stdoutLogFormat) mapLevelToStr(l log15.Lvl) string {
	switch l {
	case log15.LvlDebug:
		return "DBUG"
	case log15.LvlInfo:
		return "INFO"
	case log15.LvlWarn:
		return "WARN"
	case log15.LvlError:
		return "EROR"
	case log15.LvlCrit:
		return "CRIT"
	default:
		panic("bad level")
	}
}

func (f *stdoutLogFormat) fmtTime(time time.Time) string {
	return time.Format(timeFormat)
}
