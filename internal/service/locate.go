package service

import (
	logger2 "awesomeProject1/internal/pkg/zapx"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type LocateService interface {
	Locate(ctx context.Context, objectName string) ([]string, error)
}
type locateService struct {
	storageNodes []string //存活的节点
	client       *http.Client
	l            logger2.Loggerv1
}

func NewLocateService(storageNodes []string, l logger2.Loggerv1) LocateService {
	return &locateService{
		storageNodes: storageNodes,
		client:       &http.Client{Timeout: time.Second * 2},
		l:            l,
	}
}
func (l *locateService) Locate(ctx context.Context, objectName string) ([]string, error) {
	var result []string
	for _, node := range l.storageNodes {
		//构造向存活节点发送定位请求的url
		url := fmt.Sprintf("http://%s:%s/locate", node, objectName)
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := l.client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			continue
		}
		result = append(result, node)
	}
	if len(result) == 0 {
		return nil, errors.New("在所有节点上找不到数据")
	}
	return result, nil
}
