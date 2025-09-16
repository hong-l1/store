package service

import (
	"context"
	"time"
)

type Node struct {
	Addr     string
	LastBeat time.Time
	Alive    bool
}
type NodeService interface {
	listenHeartbeats()
	cleanDeadNodes()
	ChooseNode(ctx context.Context) *Node
}
type nodeService struct {
}

func (s *nodeService) listenHeartbeats() {
	//从消息队列获取心跳

}
func (s *nodeService) cleanDeadNodes() {
	//清楚掉长时间不心跳的节点
}
func (s *nodeService) ChooseNode(ctx context.Context) *Node {
	//选择出来一个节点
}
