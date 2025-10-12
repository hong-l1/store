package repository

import (
	logger2 "awesomeProject1/internal/pkg/zapx"
	"context"
	"io"
	"os"
	"path/filepath"
)

type StorageRepository interface {
	Upload(ctx context.Context, name string, body io.Reader) error
	Download(ctx context.Context, name string) (io.ReadCloser, error)
}
type storageRepository struct {
	l logger2.Loggerv1
}

func NewStorageRepository(l logger2.Loggerv1) StorageRepository {
	return &storageRepository{l: l}
}
func (s *storageRepository) Upload(ctx context.Context, name string, body io.Reader) error {
	f, err := os.Create(filepath.Join(os.Getenv("STORAGE_ROOT"), name))
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, body)
	return err
}

func (s *storageRepository) Download(ctx context.Context, name string) (io.ReadCloser, error) {
	f, err := os.Open(filepath.Join(os.Getenv("STORAGE_ROOT"), name))
	if err != nil {
		return nil, err
	}
	f.Seek(0, io.SeekStart)
	return f, nil
}
