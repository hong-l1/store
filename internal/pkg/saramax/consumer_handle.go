package saramax

import (
	logger2 "awesomeProject1/internal/pkg/zapx"
	"encoding/json"
	"github.com/IBM/sarama"
)

type Fn interface {
	Handle(msg *sarama.ConsumerMessage, t string) error
}

type Handler struct {
	l  logger2.Loggerv1
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
				logger2.Int32("partition", msg.Partition),
				logger2.Int64("offset", msg.Offset),
				logger2.String("topic", msg.Topic),
				logger2.Error(err))
			continue
		}
		for i := 0; i < 3; i++ {
			err = s.fn.Handle(msg, t)
			if err == nil {
				break
			}
			s.l.Error("处理消息失败",
				logger2.Int32("partition", msg.Partition),
				logger2.Int64("offset", msg.Offset),
				logger2.String("topic", msg.Topic),
				logger2.Error(err))
		}
		if err != nil {
			s.l.Error("处理消息失败-重试上限",
				logger2.Int32("partition", msg.Partition),
				logger2.Int64("offset", msg.Offset),
				logger2.String("topic", msg.Topic),
				logger2.Error(err))
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
func NewHandler(l logger2.Loggerv1, fn Fn) *Handler {
	return &Handler{
		l:  l,
		fn: fn,
	}
}
