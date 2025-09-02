package main

import (
	"awesomeProject1/store/apiserver"
	"awesomeProject1/store/event"
	"awesomeProject1/store/ioc"
	"github.com/gin-gonic/gin"
)

func main() {
	l := ioc.InitLogger()
	server := gin.Default()
	RegisterRoute(server)
	ioc.InitSyncProducer()
	client := ioc.Initkafka()
	consumer := event.NewKafkaConsumer(l, client, "api", "dataerver")
	cos := ioc.InitConsumers(consumer)
	fn := apiserver.NewFn()
	go func() {
		cos.Start(fn)
	}()
	go apiserver.RemoveExpiredDataServer()

}
func RegisterRoute(server *gin.Engine) {

}
