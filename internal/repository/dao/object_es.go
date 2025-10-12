package dao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strconv"
)

const ObjectIndex = "object_index"

type objectDao struct {
	client *elastic.Client
}

type ObjectDao interface {
	GetMetadata(ctx context.Context, name string, version int) (Metadata, error)
	SearchLatestVersion(ctx context.Context, name string) (Metadata, error)
	PutMetadata(ctx context.Context, name string, version int, size int64, hash string) error
	AddVersion(ctx context.Context, name, hash string, size int64) error
	SearchAllVersions(ctx context.Context, name string) ([]Metadata, error)
}

func NewObjectDao(client *elastic.Client) ObjectDao {
	return &objectDao{client: client}
}
func (o *objectDao) SearchAllVersions(ctx context.Context, name string) ([]Metadata, error) {
	result, err := o.client.Search(ObjectIndex).Sort("version", true).Size(10).
		Query(elastic.NewTermQuery("name", name)).Do(ctx)
	if err != nil {
		return []Metadata{}, err
	}
	if result.Hits.TotalHits == nil || result.Hits.TotalHits.Value == 0 {
		return nil, fmt.Errorf("no metadata found for name=%s", name)
	}
	metas := make([]Metadata, 0, result.TotalHits())
	for _, hit := range result.Hits.Hits {
		var m Metadata
		err = json.Unmarshal(hit.Source, &m)
		if err != nil {
			return []Metadata{}, err
		}
		metas = append(metas, m)
	}
	return metas, nil
}
func (o *objectDao) SearchLatestVersion(ctx context.Context, name string) (Metadata, error) {
	result, err := o.client.Search(ObjectIndex).Sort("version", false).Size(1).Query(elastic.NewTermQuery("name", name)).Do(ctx)
	if err != nil {
		return Metadata{}, err
	}
	if len(result.Hits.Hits) == 0 {
		return Metadata{}, fmt.Errorf("no metadata found for name=%s", name)
	}
	var m Metadata
	err = json.Unmarshal(result.Hits.Hits[0].Source, &m)
	if err != nil {
		return Metadata{}, err
	}
	return m, nil
}
func (o *objectDao) GetMetadata(ctx context.Context, name string, version int) (Metadata, error) {
	if version == 0 {
		return o.SearchLatestVersion(ctx, name)
	}
	boolquery := elastic.NewBoolQuery().Must(elastic.NewTermQuery("name", name),
		elastic.NewTermQuery("version", version))
	res, err := o.client.Search(ObjectIndex).Query(boolquery).Size(1).Do(ctx)
	if err != nil {
		return Metadata{}, err
	}
	if res.Hits.TotalHits.Value == 0 {
		return Metadata{}, errors.New("没找到元数据")
	}
	var meta Metadata
	if err := json.Unmarshal(res.Hits.Hits[0].Source, &meta); err != nil {
		return Metadata{}, fmt.Errorf("反序列化失败: %w", err)
	}
	return meta, nil
}
func (o *objectDao) PutMetadata(ctx context.Context, name string, version int, size int64, hash string) error {
	meta := Metadata{
		Name:    name,
		Version: version,
		Size:    size,
		Hash:    hash,
	}
	id := name + strconv.Itoa(version)
	_, err := o.client.Index().Index(ObjectIndex).Id(id).BodyJson(meta).Do(ctx)
	return err
}
func (o *objectDao) AddVersion(ctx context.Context, name, hash string, size int64) error {
	meta, err := o.SearchLatestVersion(ctx, name)
	if err != nil {
		return err
	}
	return o.PutMetadata(ctx, name, meta.Version+1, size, hash)
}
