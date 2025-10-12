package service

import (
	"awesomeProject1/internal/domain"
	"awesomeProject1/internal/repository"
	"context"
)

type MetadataService interface {
	Versions(ctx context.Context, name string) ([]domain.Meta, error)
	GetVersionDetail(ctx context.Context, name string, version int) (domain.Meta, error)
	SearchLatestVersion(ctx context.Context, name string) (domain.Meta, error)
	PutNewMetadata(ctx context.Context, filename string, hash string, size int64) error
	GetMetadata(ctx context.Context, name string, version int) (domain.Meta, error)
}
type metadataservice struct {
	repo repository.MetadataRepo
}

func NewMetadataService(repo repository.MetadataRepo) MetadataService {
	return &metadataservice{repo: repo}
}
func (o *metadataservice) GetMetadata(ctx context.Context, name string, version int) (domain.Meta, error) {
	return o.repo.GetMetadata(ctx, name, version)
}

func (o *metadataservice) Versions(ctx context.Context, name string) ([]domain.Meta, error) {
	return o.repo.GetVersions(ctx, name)
}

func (o *metadataservice) GetVersionDetail(ctx context.Context, name string, version int) (domain.Meta, error) {
	return o.repo.GetVersionDetail(ctx, name, version)
}

func (o *metadataservice) SearchLatestVersion(ctx context.Context, name string) (domain.Meta, error) {
	return o.repo.SearchLatestVersion(ctx, name)
}

func (o *metadataservice) PutNewMetadata(ctx context.Context, name string, hash string, size int64) error {
	return o.repo.UpSertMetadata(ctx, name, hash, size)
}
