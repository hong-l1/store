package event

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type SyncProducer interface {
	SendMessage(msg string, topic string) error
}
type SaramaProducer struct {
	producer sarama.SyncProducer
}

func (s SaramaProducer) SendMessage(msg string, topic string) error {
	data, err := json.Marshal(msg)
	_, _, err = s.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(data),
	})
	return err
}
func NewSaramaProducer(producer sarama.SyncProducer) SyncProducer {
	return SaramaProducer{
		producer: producer,
	}
}
