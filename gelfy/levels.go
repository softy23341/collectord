package gelfy

// Graylog using syslog like log levels.
const (
	// LvlEmerg usualy reserved for kernel logging.
	LvlEmerg = iota

	// LvlAlert usualy reserved for kernel logging.
	LvlAlert

	// LvlCrit for fatal situation.
	LvlCrit

	// LvlError common errors.
	LvlError

	// LvlWarn for situations wich can became a error in some cases.
	LvlWarn

	// LvlNotice valueable information.
	LvlNotice

	// LvlInfo any information which can be usefull to end-user.
	LvlInfo

	// LvlDebug information for developers.
	LvlDebug
)
