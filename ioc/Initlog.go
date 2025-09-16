package ioc

import "go.uber.org/zap"

func InitZapLogger() *zap.Logger {
	l, _ := zap.NewDevelopment()
	return l
}
