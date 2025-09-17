package service

import (
	logger2 "awesomeProject1/internal/pkg/zapx"
	"awesomeProject1/internal/repository"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"sync"
	"time"
)

type Node struct {
	Name     string
	Addr     string
	LastBeat time.Time
	Alive    bool
}

//NodeService 管理全局节点
type NodeService interface {
	ListenHeartbeats(partition int32)
	CleanDeadNodes()
	Locate(ctx context.Context, object string) string
	ChooseNode(ctx context.Context) string
	ListenLocateMsg(partition int32)
	SendMsg() error
	AddNode(node *Node)
}
type nodeService struct {
	lock  sync.RWMutex
	nodes map[string]*Node
	p     sarama.SyncProducer
	c     sarama.Consumer
	l     logger2.Loggerv1
	repo  repository.LocateRepository
}

//给post用的，选择活节点
func (s *nodeService) ChooseNode(ctx context.Context) string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, v := range s.nodes {
		return v.Addr
	}
	return ""
}

func NewNodeService(l logger2.Loggerv1) NodeService {
	return &nodeService{l: l}
}
func (s *nodeService) ListenHeartbeats(partition int32) {
	//从消息队列获取心跳
	pc, err := s.c.ConsumePartition("hearbeat", partition, sarama.OffsetNewest)
	if err != nil {
		s.l.Error("创建分区消费者失败", logger2.Error(err), logger2.Int32("partition", partition))
		return
	}
	defer pc.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	for {
		select {
		case msg, ok := <-pc.Messages():
			if !ok {
				s.l.Info("消息通道已关闭", logger2.Int32("partition", partition))
			}
			var node *Node
			err = json.Unmarshal(msg.Value, node)
			s.AddNode(node)
		case <-ctx.Done():
			cancel()
		}
	}
}
func (s *nodeService) CleanDeadNodes() {
	//清掉长时间不心跳的节点
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		s.lock.Lock()
		now := time.Now()
		for _, node := range s.nodes {
			if now.Sub(node.LastBeat) > time.Second*10 {
				node.Alive = false
			}
		}
		s.lock.Unlock()
	}
}
func (s *nodeService) Locate(ctx context.Context, object string) string {
	//定位所需资源在哪个节点
	s.lock.RLock()
	for _, node := range s.nodes {
		if node.Alive {
			_, _, err := s.p.SendMessage(&sarama.ProducerMessage{
				Topic: "locate",
				Value: sarama.StringEncoder(object),
			})
			if err != nil {
				//打日志
				continue
			}
		}
	}
	s.lock.RUnlock()
	var results []string
	go func() {
		pc, err1 := s.c.ConsumePartition("locate", int32(0), sarama.OffsetNewest)
		if err1 != nil {
			s.l.Error("创建分区消费者失败", logger2.Error(err1), logger2.Int32("partition", int32(0)))
			return
		}
		c, cancel := context.WithTimeout(context.Background(), time.Second*2)
		for {
			select {
			case <-c.Done():
				cancel()
			case msg, ok := <-pc.Messages():
				if !ok {

				}
				var addr *string
				err1 = json.Unmarshal(msg.Value, &addr)
				if err1 != nil {
				}
				results = append(results, string(msg.Value))
			}
		}
	}()
	return results[0]
}
func (s *nodeService) ListenLocateMsg(partition int32) {
	pc, err := s.c.ConsumePartition("locate", partition, sarama.OffsetNewest)
	if err != nil {
		s.l.Error("创建分区消费者失败", logger2.Error(err), logger2.Int32("partition", partition))
		return
	}
	defer pc.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	for {
		select {
		case msg, ok := <-pc.Messages():
			if !ok {
			}
			var name string
			err = json.Unmarshal(msg.Value, &name)
			if s.repo.Locate(name) {
				err = s.SendMsg()
				if err != nil {
					s.l.Error("发送消息失败", logger2.Error(err), logger2.Int32("partition", partition))
				}
			}
		case <-ctx.Done():
			cancel()
		}
	}
}
func (s *nodeService) AddNode(node *Node) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.nodes[node.Name] = node
}
func (s *nodeService) SendMsg() error {
	ip := s.repo.Getip()
	_, _, err := s.p.SendMessage(&sarama.ProducerMessage{
		Topic: "object",
		Value: sarama.StringEncoder(ip),
	})
	return err
}
