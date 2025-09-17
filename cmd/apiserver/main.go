package main

import (
	"awesomeProject1/internal/api"
	logger2 "awesomeProject1/internal/pkg/zapx"
	"awesomeProject1/internal/service"
	"awesomeProject1/ioc"
	"github.com/gin-gonic/gin"
)

func main() {
	l := ioc.InitZapLogger()
	logger := logger2.NewZapLogger(l)
	node := service.NewNodeService(logger)
	go node.ListenHeartbeats(int32(0))
	go node.CleanDeadNodes()
	handler := api.NewApiServerHandler(node, logger)
	server := gin.Default()
	handler.RegisterApiServerRoute(server)
}
