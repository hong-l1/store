package repository

import (
	"os"
	"path/filepath"
)

type LocateRepository interface {
	Locate(name string) bool
	Getip() string
}
type locateRepository struct {
}

func NewLocateRepository() LocateRepository {
	return &locateRepository{}
}
func (r *locateRepository) Locate(name string) bool {
	root := os.Getenv("STORAGE_ROOT")
	if root == "" {
		root = "."
	}
	full := filepath.Join(root, name)
	_, err := os.Stat(full)
	return !os.IsNotExist(err)
}
func (r *locateRepository) Getip() string {
	//envIP := os.Getenv("NODE_IP")
	//return envIP
	return "127.0.0.1:8080"
}
