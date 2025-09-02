package locate

import (
	"errors"
	"github.com/IBM/sarama"
	"os"
)

var Notvaliderr = errors.New("filename not valid")

const topic = "apiserver"

func (l *locate) Locate(name string) error {
	_, err := os.Stat(name)
	if !os.IsNotExist(err) {
		return Notvaliderr
	}
	return nil
}

type locate struct {
}

func (l *locate) Handle(msg *sarama.ConsumerMessage, t string) error {
	err := l.Locate(t)
	//在这里继续调用producer来发送true/false
	if err != nil {
		return err
	}
	return nil
}
func NewLocate() *locate {
	return &locate{}
}
