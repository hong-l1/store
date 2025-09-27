//go:build wireinject

package main

import (
	"awesomeProject1/internal/api"
	logger2 "awesomeProject1/internal/pkg/zapx"
	"awesomeProject1/internal/repository"
	service2 "awesomeProject1/internal/service"
	"awesomeProject1/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitDateServer() *gin.Engine {
	wire.Build(
		InitGinEngine,
		ioc.InitZapLogger,
		logger2.NewZapLogger,
		repository.NewStorageRepository,
		service2.NewNodeService,
		api.NewObjectHandler,
		ioc.Initkafka,
		ioc.InitSyncProducer,
		ioc.InitSyncConsumer,
		repository.NewLocateRepository,
		service2.NewObjectService,
	)
	return new(gin.Engine)
}
func InitGinEngine(ObjectHandler *api.DataServerHandler) *gin.Engine {
	server := gin.Default()
	ObjectHandler.RegisterObjectRoute(server)
	return server
}
