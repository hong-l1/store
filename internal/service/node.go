package service

import (
	logger2 "awesomeProject1/internal/pkg/zapx"
	"awesomeProject1/internal/repository"
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type Node struct {
	Name     string
	Addr     string
	LastBeat time.Time
	Alive    bool
}

// NodeService 管理全局节点
type NodeService interface {
	ListenHeartbeats(partition int32)
	CleanDeadNodes()
	Locate(ctx context.Context, object string) string
	ChooseNode(ctx context.Context) string
	ListenLocateMsg(partition int32)
	SendMsg() error
	AddNode(node *Node)
	HeartbeatTicker()
}
type nodeService struct {
	lock  *sync.RWMutex
	nodes map[string]*Node
	p     sarama.SyncProducer
	c     sarama.Consumer
	l     logger2.Loggerv1
	repo  repository.LocateRepository
}

// 给post用的，选择活节点
func (s *nodeService) ChooseNode(ctx context.Context) string {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, v := range s.nodes {
		if v.Alive {
			return v.Addr
		}
	}
	return ""
}
func NewNodeService(l logger2.Loggerv1, p sarama.SyncProducer, c sarama.Consumer, repo repository.LocateRepository) NodeService {
	return &nodeService{
		l:     l,
		lock:  &sync.RWMutex{},
		p:     p,
		c:     c,
		repo:  repo,
		nodes: make(map[string]*Node),
	}
}
func (s *nodeService) ListenHeartbeats(partition int32) {
	//从消息队列获取心跳
	pc, err := s.c.ConsumePartition("heartbeat", partition, sarama.OffsetNewest)
	if err != nil {
		s.l.Error("创建分区消费者失败", logger2.Error(err), logger2.Int32("partition", partition))
		return
	}
	defer pc.Close()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		select {
		case msg, ok := <-pc.Messages():
			if !ok {
				s.l.Info("消息通道已关闭", logger2.Int32("partition", partition))
				return
			}
			var node Node
			err = json.Unmarshal(msg.Value, &node)
			if err != nil {
				s.l.Error("心跳反序列化失败", logger2.Error(err))
				continue
			}
			node.LastBeat = time.Now()
			node.Alive = true
			s.AddNode(&node)
		case <-ctx.Done():
			return
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
	// 快速实现：返回任一活跃节点地址
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, node := range s.nodes {
		if node.Alive {
			return node.Addr
		}
	}
	return ""
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
func (s *nodeService) HeartbeatTicker() {
	ip := s.repo.Getip()
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		res := strings.Split(ip, ":")
		node := Node{
			Name:     res[len(res)-1],
			Addr:     ip,
			LastBeat: time.Now(),
			Alive:    true,
		}
		data, err := json.Marshal(&node)
		if err != nil {
			panic(err)
		}
		_, _, err = s.p.SendMessage(&sarama.ProducerMessage{
			Topic: "heartbeat",
			Value: sarama.ByteEncoder(data),
		})
		if err != nil {
			panic(err)
		}
	}
}
