package storage

import (
	"os"

	"github.com/raonismaneoto/custom-nfs-server/nfs-server/storage/dht"
	"github.com/raonismaneoto/custom-nfs-server/nfs-server/storage/vanilla"
)

type Storage interface {
	SaveAsync(id, path string, content <-chan []byte, errors chan<- error)
	Save(id, path string, content []byte) error
	Read(id, path string, content chan<- []byte, errors chan<- error)
	Rm(id, path string) error
}

func New(t string) Storage {
	if t == "dht" {
		return dht.New()
	} else if t == "vanilla" {
		root := os.Getenv("ROOT_FOLDER")
		return vanilla.New(root)
	}
	return nil
}
