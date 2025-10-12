package repository

import (
	"net"
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
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return ""
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
