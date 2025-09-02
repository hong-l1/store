package ioc

import (
	logger "awesomeProject1/store/pkg/zapx"
	"go.uber.org/zap"
)

func InitLogger() logger.Loggerv1 {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
