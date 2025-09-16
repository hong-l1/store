package logger

import "go.uber.org/zap"

func NewZapLogger(log *zap.Logger) Loggerv1 {
	return &ZapLogger{
		log: log,
	}
}

type ZapLogger struct {
	log *zap.Logger
}

func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.log.Debug(msg, z.ToZaoField(args)...)
}
func (z *ZapLogger) Info(msg string, args ...Field) {
	z.log.Info(msg, z.ToZaoField(args)...)
}
func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.log.Warn(msg, z.ToZaoField(args)...)
}
func (z *ZapLogger) Error(msg string, args ...Field) {
	z.log.Error(msg, z.ToZaoField(args)...)
}
func (z ZapLogger) ToZaoField(args []Field) []zap.Field {
	ans := make([]zap.Field, 0, len(args))
	for k := range args {
		ans = append(ans, zap.Any(args[k].Key, args[k].Value))
	}
	return ans
}
func String(key, val string) Field {
	return Field{Key: key, Value: val}
}
func Error(err error) Field {
	return Field{Key: "error", Value: err}
}
func Int64(key string, val int64) Field {
	return Field{Key: key, Value: val}
}
func Int32(key string, val int32) Field {
	return Field{Key: key, Value: val}
}
