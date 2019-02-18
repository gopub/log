package log

type Level int

const (
	AllLevel Level = iota + 1
	TraceLevel
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
)

func (l Level) String() string {
	switch l {
	case AllLevel:
		return "ALL"
	case TraceLevel:
		return "TRA"
	case DebugLevel:
		return "DEB"
	case InfoLevel:
		return "INF"
	case WarnLevel:
		return "WRN"
	case ErrorLevel:
		return "ERR"
	case FatalLevel:
		return "FAT"
	case PanicLevel:
		return "PAN"
	default:
		return ""
	}
}
