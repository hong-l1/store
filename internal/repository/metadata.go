package repository

import (
	"awesomeProject1/internal/domain"
	"awesomeProject1/internal/repository/dao"
	"context"
)

type MetadataRepo interface {
	GetVersions(ctx context.Context, name string) ([]domain.Meta, error)
	GetVersionDetail(ctx context.Context, name string, version int) (domain.Meta, error)
	SearchLatestVersion(ctx context.Context, name string) (domain.Meta, error)
	UpSertMetadata(ctx context.Context, name string, hash string, size int64) error
	GetMetadata(ctx context.Context, name string, version int) (domain.Meta, error)
}
type metadataRepo struct {
	dao dao.ObjectDao
}

func (m *metadataRepo) GetVersions(ctx context.Context, name string) ([]domain.Meta, error) {
	metas, err := m.dao.SearchAllVersions(ctx, name)
	if err != nil || len(metas) == 0 {
		return []domain.Meta{}, err
	}
	meta := make([]domain.Meta, len(metas))
	for k := range metas {
		meta = append(meta, m.ToDomain(metas[k]))
	}
	return meta, nil
}

func (m *metadataRepo) GetMetadata(ctx context.Context, name string, version int) (domain.Meta, error) {
	meta, err := m.dao.GetMetadata(ctx, name, version)
	if err != nil {
		return domain.Meta{}, err
	}
	return m.ToDomain(meta), nil
}

func NewMetadataRepo(dao dao.ObjectDao) MetadataRepo {
	return &metadataRepo{dao}
}
func (m *metadataRepo) UpSertMetadata(ctx context.Context, name string, hash string, size int64) error {
	return m.dao.AddVersion(ctx, name, hash, size)
}

func (m *metadataRepo) SearchLatestVersion(ctx context.Context, name string) (domain.Meta, error) {
	meta, err := m.dao.SearchLatestVersion(ctx, name)
	if err != nil {
		return domain.Meta{}, err
	}
	return m.ToDomain(meta), nil
}

func (m *metadataRepo) GetVersionDetail(ctx context.Context, name string, version int) (domain.Meta, error) {
	meta, err := m.dao.GetMetadata(ctx, name, version)
	if err != nil {
		return domain.Meta{}, err
	}
	return m.ToDomain(meta), nil
}
func (m *metadataRepo) ToDomain(meta dao.Metadata) domain.Meta {
	return domain.Meta{
		Filename: meta.Name,
		Size:     meta.Size,
		Version:  meta.Version,
		Hash:     meta.Hash,
	}
}
