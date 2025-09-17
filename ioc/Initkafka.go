package ioc

import "github.com/IBM/sarama"

var addres = []string{"127.0.0.1:9094"}

func Initkafka() sarama.Client {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.AutoCommit.Enable = true
	client, err := sarama.NewClient(addres, cfg)
	if err != nil {
		panic(err)
	}
	return client
}
func InitSyncProducer(c sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	return p
}
func InitSyncConsumer(client sarama.Client) sarama.Consumer {
	c, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		panic(err)
	}
	return c
}
