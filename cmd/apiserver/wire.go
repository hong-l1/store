//go:build wireinject

package main

import (
	"awesomeProject1/internal/api"
	logger2 "awesomeProject1/internal/pkg/zapx"
	"awesomeProject1/internal/repository"
	"awesomeProject1/internal/service"
	"awesomeProject1/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitApiServer() *gin.Engine {
	wire.Build(ioc.InitZapLogger,
		logger2.NewZapLogger,
		ioc.Initkafka,
		ioc.InitSyncProducer,
		ioc.InitSyncConsumer,
		repository.NewLocateRepository,
		service.NewNodeService,
		api.NewApiServerHandler,
		InitGinEngine,
	)
	return new(gin.Engine)
}

func InitGinEngine(ApiHandler *api.ApiServerHandler) *gin.Engine {
	server := gin.Default()
	ApiHandler.RegisterApiServerRoute(server)
	return server
}
