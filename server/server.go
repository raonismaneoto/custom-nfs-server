package server

import (
	"log"
	"os"
)

type Server struct {
	root string
}

func New(root string) *Server{
	return &Server{
		root: root,
	}
}

func (s *Server) Save(id string, path string, content []byte) (error) {
	log.Println("Save call received.")
	log.Println("saving")
	f, err := os.OpenFile(s.root + path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("unable to open/create %v", path)
	}

	defer f.Close()

	if _, err := f.Write(content); err != nil {
        log.Println("unable to write to %v", path)
    }

	return err
}

func (s *Server) Read(id string, path string, offset, limit int32) ([]byte, error) {
	//check if the file exists, if id is allowed to access its content and if offset and limit are ok regarding the file size
	f, err := os.Open(path)
	if err != nil {

	}

	content := make([]byte, limit)

	if _, err := f.ReadAt(content, int64(offset)); err != nil {
		return nil, err
	}
	return content, nil
}