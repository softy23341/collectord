package loggy

import (
	"fmt"
	"runtime"

	"github.com/inconshreveable/log15"

	"git.softndit.com/collector/backend/gelfy"
)

// Flagger TBD
type Flagger interface {
	String(key string) string
}

// StdoutAndGraylogHandler TBD
func StdoutAndGraylogHandler(f Flagger, appName string) (log15.Handler, error) {
	hs := []log15.Handler{}

	if llStr := f.String("stdoutloglvl"); llStr != "" {
		ll, err := log15.LvlFromString(llStr)
		if err != nil {
			return nil, fmt.Errorf("invalid stdoutloglvl: %v", err)
		}
		hs = append(hs, LvlFilterHandler(ll, StdoutLogHandler))
	}

	if srv, llStr := f.String("graylogsrv"), f.String("grayloglvl"); srv != "" && llStr != "" {
		ll, err := log15.LvlFromString(llStr)
		if err != nil {
			return nil, fmt.Errorf("invalid grayloglvl: %v", err)
		}

		extra := make(gelfy.ExtraFields)
		if appName != "" {
			extra["app"] = appName
		}
		h, err := GraylogHandler(srv, GelfFormat(extra))
		if err != nil {
			return nil, fmt.Errorf("invalid graylogsrv: %v", err)
		}

		hs = append(hs, LvlFilterHandler(ll, h))
	}

	return MultiHandler(hs...), nil
}

// LogPanic -- useful only with deffer
func LogPanic(logger log15.Logger) {
	if err := recover(); err != nil {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]
		errInfo := fmt.Sprintf("%v", err)
		stackInfo := fmt.Sprintf("%s", buf)
		logger.Crit("panic", "err", errInfo, "stack", stackInfo)
	}
}
