package dao

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/olivere/elastic/v7"
	"time"
)

//go:embed init_index.go
var objectIndexCfg string

func InitEsIndex(client *elastic.Client) error {
	ctx, canel := context.WithTimeout(context.Background(), 2*time.Second)
	defer canel()
	return tryCreateIndex(ctx, client, ObjectIndex, objectIndexCfg)
}
func tryCreateIndex(ctx context.Context, client *elastic.Client, idxName, idxCfg string) error {
	ok, err := client.IndexExists(idxName).Do(ctx)
	if err != nil {
		return fmt.Errorf("检测 %s 是否存在失败 %w", idxName, err)
	}
	if ok {
		return nil
	}
	_, err = client.CreateIndex(idxName).Body(idxCfg).Do(ctx)
	if err != nil {
		return fmt.Errorf("创建%s失败 %w", idxName, err)
	}
	return nil
}
