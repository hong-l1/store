package logger

type Field struct {
	Key   string
	Value any
}

type Loggerv1 interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
}
