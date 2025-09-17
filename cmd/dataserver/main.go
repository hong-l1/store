package main

import (
	"awesomeProject1/internal/api"
	logger2 "awesomeProject1/internal/pkg/zapx"
	"awesomeProject1/internal/repository"
	"awesomeProject1/internal/service"
	"awesomeProject1/ioc"
	"encoding/json"
	"github.com/IBM/sarama"
	"strings"
	"time"
)

func main() {
	l := ioc.InitZapLogger()
	logger := logger2.NewZapLogger(l)
	reposto := repository.NewStorageRepository(logger)
	repoloa := repository.NewLocateRepository()
	objSvc := service.NewObjectService(reposto)
	node := service.NewNodeService(logger)
	api.NewObjectHandler(objSvc, logger)
	client := ioc.Initkafka()
	p := ioc.InitSyncProducer(client)
	go heartbeat(p, repoloa.Getip())
	go node.ListenLocateMsg(int32(0))
}
func heartbeat(p sarama.SyncProducer, ip string) {
	res := strings.Split(ip, ":")
	var node = service.Node{
		Name:     res[len(res)-1],
		Addr:     ip,
		LastBeat: time.Now(),
	}
	data, err := json.Marshal(&node)
	if err != nil {
		panic(err)
	}
	_, _, err = p.SendMessage(&sarama.ProducerMessage{
		Topic: "heartbeat",
		Value: sarama.StringEncoder(data),
	})
	if err != nil {
		panic(err)
	}
}
