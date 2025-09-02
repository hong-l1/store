package main

import (
	"awesomeProject1/store/dataserver/handle"
	"awesomeProject1/store/dataserver/hearbeat"
	"awesomeProject1/store/dataserver/locate"
	"awesomeProject1/store/event"
	"awesomeProject1/store/ioc"
	"github.com/gin-gonic/gin"
	"net"
	"os"
)

func main() {
	l := ioc.InitLogger()
	client := ioc.Initkafka()
	pro := ioc.InitSyncProducer(client)
	late := locate.NewLocate()
	consumer := event.NewKafkaConsumer(l, client, "data", "apiserver")
	cos := ioc.InitConsumers(consumer)
	go func() {
		cos.Start(late)
	}()
	go hearbeat.SendHeartbeat(pro, GetLocalIp())
	server := gin.Default()
	RegisterRoute(server)
}
func RegisterRoute(server *gin.Engine) {
	group := server.Group("/objects")
	group.GET("/:filename", handle.GetObjects)
	group.PUT("/:filename", handle.PutObjects)
}
func GetLocalIp() string {
	host, _ := os.Hostname()
	ips, _ := net.LookupIP(host)
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return ""
}
