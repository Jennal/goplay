package log

var defaultLogger Logger = NewLogger("", 4)

func Log(line string) {
	defaultLogger.Log(line)
}

func Logf(format string, args ...interface{}) {
	defaultLogger.Logf(format, args...)
}

func Error(err error) {
	defaultLogger.Error(err)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}
