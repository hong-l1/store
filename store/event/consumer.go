package event

import (
	"awesomeProject1/store/pkg/saramax"
	logger "awesomeProject1/store/pkg/zapx"
	"context"
	"github.com/IBM/sarama"
	"time"
)

type Consumer interface {
	Start(fn saramax.Fn) error
}
type KafkaConsumer struct {
	l       logger.Logger
	client  sarama.Client
	groupId string
	topic   string
}

func (k *KafkaConsumer) Start(fn saramax.Fn) error {
	client, err := sarama.NewConsumerGroupFromClient(k.groupId, k.client)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	go func() {
		er := client.Consume(ctx, []string{k.topic}, saramax.NewHandler(k.l, fn))
		if er != nil {
			k.l.Error("退出消费循环异常", logger.Error(err))
		}
	}()
	return err
}
func NewKafkaConsumer(l logger.Logger, client sarama.Client, groupId string, topic string) *KafkaConsumer {
	return &KafkaConsumer{
		l:       l,
		client:  client,
		groupId: groupId,
		topic:   topic,
	}
}
