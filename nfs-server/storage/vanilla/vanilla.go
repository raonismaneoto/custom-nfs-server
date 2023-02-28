package vanilla

import (
	"io/fs"
	"log"
	"math"
	"os"

	"github.com/raonismaneoto/custom-nfs-server/helpers"
)

type VanillaStorage struct {
	root                string
	MaxBytesPerResponse int32
}

func New(root string) *VanillaStorage {
	return &VanillaStorage{
		root: root,
	}
}

func (s VanillaStorage) Save(id, path string, content <-chan []byte, errors chan<- error) {
	f, err := os.OpenFile(s.root+path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("unable to open/create %v", path)
		errors <- err
		close(errors)
		return
	}

	defer f.Close()

	for {
		currContent, ok := <-content
		if !ok {
			close(errors)
			return
		}

		if _, err := f.Write(currContent); err != nil {
			log.Println("unable to write to %v", path)
		}
	}

}

func (s VanillaStorage) Read(id, path string, content chan<- []byte, errors chan<- error) {
	var f fs.FileInfo
	var err error
	if f, err = os.Stat(s.root + path); err != nil {
		log.Println("unable to open the file %v", s.root+path)
		errors <- err
	}

	chuncks := int32(math.Ceil(float64(f.Size()) / float64(s.MaxBytesPerResponse)))
	for i := int32(0); i < chuncks; i++ {
		currContent, err := helpers.ReadFileChunk(s.root+path, i*s.MaxBytesPerResponse, s.MaxBytesPerResponse)
		if err != nil {
			log.Println("error while reading file content: ", err)
			errors <- err
		}

		content <- currContent
	}

}
