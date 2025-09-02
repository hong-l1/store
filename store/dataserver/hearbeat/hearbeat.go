package hearbeat

import (
	"awesomeProject1/store/event"
	"github.com/IBM/sarama"
	"time"
)

func SendHeartbeat(producer sarama.SyncProducer, ip string) {
	p := event.NewSaramaProducer(producer)
	for {
		p.SendMessage(ip, "apiserver")
		time.Sleep(5 * time.Second)
	}
}
