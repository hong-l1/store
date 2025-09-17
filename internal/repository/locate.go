package repository

import (
	"os"
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
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
func (r *locateRepository) Getip() string {
	envIP := os.Getenv("NODE_IP")
	return envIP
}
