package logger

type logItem struct {
	Level   logLevel
	Content any
}

type logLevel uint8

const (
	lInfo logLevel = iota
	lWarn
	lError
	lFatal
	lDebug
)

func (level logLevel) String() string {
	switch level {
	case lInfo:
		return "INFO"
	case lWarn:
		return "WARN"
	case lError:
		return "ERROR"
	case lFatal:
		return "FATAL"
	case lDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}
