package rest

type Logger interface {
	Printf(format string, v ...interface{})
}

type noOpLogger struct {
}

func (l noOpLogger) Printf(format string, v ...interface{}) {
}

type logcat struct {
	info Logger
	warn Logger
}
