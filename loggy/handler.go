package loggy

import (
	"git.softndit.com/collector/backend/gelfy"
	"github.com/inconshreveable/log15"
)

// GraylogHandler implements log15.Handler interface.
// GraylogHandler writes log records formated with fmtr formatter to graylog server via UDP.
// Server ip and port must be specified in addr.
func GraylogHandler(addr string, fmtr log15.Format) (log15.Handler, error) {
	tr, err := gelfy.NewTransport(addr, gelfy.DefaultOptions)
	if err != nil {
		return nil, err
	}
	return &graylogHandler{fmtr: fmtr, tr: tr}, nil
}

type graylogHandler struct {
	fmtr log15.Format
	tr   *gelfy.Transport
}

func (h *graylogHandler) Log(r *log15.Record) error {
	return h.tr.SendRaw(h.fmtr.Format(r))
}

// LvlFilterHandler returns a Handler that only writes
// records which are less than the given verbosity
func LvlFilterHandler(maxLvl log15.Lvl, h log15.Handler) log15.Handler {
	if maxLvl == log15.LvlDebug {
		return h
	}
	return log15.LvlFilterHandler(maxLvl, h)
}

// MultiHandler dispatches any write to each of its handlers.
func MultiHandler(hs ...log15.Handler) log15.Handler {
	switch len(hs) {
	case 0:
		return log15.DiscardHandler()
	case 1:
		return hs[0]
	default:
		return log15.MultiHandler(hs...)
	}
}
