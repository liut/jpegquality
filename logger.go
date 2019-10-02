package jpegquality

var (
	logger Logger = &logDiscard{}
)

// Logger ...
type Logger interface {
	Print(args ...interface{})
	Printf(format string, args ...interface{})
}

// SetLogger ...
func SetLogger(l Logger) {
	logger = l
}

// GetLogger ...
func GetLogger() Logger {
	return logger
}

type logDiscard struct{}

func (l *logDiscard) Print(args ...interface{})                 {}
func (l *logDiscard) Printf(format string, args ...interface{}) {}
