package main

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type App struct {
	server   *gin.Engine
	producer sarama.AsyncProducer
	comsumer sarama.Consumer
}
