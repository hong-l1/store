package service

import (
	"awesomeProject1/internal/repository"
	"context"
	"io"
)

type ObjectService interface {
	Upload(context context.Context, name string, body io.ReadCloser) error
	Download(context context.Context, name string) (io.ReadCloser, error)
}
type objectService struct {
	repo repository.StorageRepository
}

func NewObjectService(repo repository.StorageRepository) ObjectService {
	return &objectService{repo: repo}
}
func (o *objectService) Upload(context context.Context, name string, body io.ReadCloser) error {
	return o.repo.Upload(context, name, body)
}

func (o *objectService) Download(context context.Context, name string) (io.ReadCloser, error) {
	return o.repo.Download(context, name)
}
