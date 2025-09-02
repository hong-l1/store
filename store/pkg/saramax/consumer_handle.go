package saramax

import (
	logger "awesomeProject1/store/pkg/zapx"
	"encoding/json"
	"github.com/IBM/sarama"
)

type Fn interface {
	Handle(msg *sarama.ConsumerMessage, t string) error
}

type Handler struct {
	l  logger.Logger
	fn Fn
}

func (s Handler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s Handler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t string
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			s.l.Error("反序列化失败",
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset),
				logger.String("topic", msg.Topic),
				logger.Error(err))
			continue
		}
		for i := 0; i < 3; i++ {
			err = s.fn.Handle(msg, t)
			if err == nil {
				break
			}
			s.l.Error("处理消息失败",
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset),
				logger.String("topic", msg.Topic),
				logger.Error(err))
		}
		if err != nil {
			s.l.Error("处理消息失败-重试上限",
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset),
				logger.String("topic", msg.Topic),
				logger.Error(err))
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
func NewHandler(l logger.Logger, fn Fn) *Handler {
	return &Handler{
		l:  l,
		fn: fn,
	}
}
