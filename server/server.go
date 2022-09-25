package server

import (
	"errors"
	"log"
	"os"
	"strings"
)

const MetaFileSuffix string = "meta"

type Server struct {
	root string
}

func New(root string) *Server {
	if _, err := os.Stat(root); err != nil {
		os.Mkdir(root, 0777)
	}
	return &Server{
		root: root,
	}
}

func (s *Server) Save(id, path string, content []byte) error {
	log.Println("Save call received.")
	log.Println("saving")
	f, err := os.OpenFile(s.root+path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("unable to open/create %v", path)
		return err
	}

	defer f.Close()

	if _, err := os.Stat(s.root + path + MetaFileSuffix); err != nil {
		fm, err := os.OpenFile(s.root+path+MetaFileSuffix, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("unable to open/create %v", path+MetaFileSuffix)
			return err
		}

		defer fm.Close()

		if _, err := fm.Write([]byte(id)); err != nil {
			log.Println("unable to write to %v", path+MetaFileSuffix)
			return err
		}
	}

	// TODO: check if it appends the content
	if _, err := f.Write(content); err != nil {
		log.Println("unable to write to %v", path)
	}

	return err
}

func (s *Server) Read(id, path string, offset, limit int32) ([]byte, error) {
	//check if id is allowed to access its content
	fm, err := os.Open(s.root + path + MetaFileSuffix)
	if err != nil {
		log.Println("unable to read %v", path+MetaFileSuffix)
		return nil, err
	}

	defer fm.Close()

	mcontent := make([]byte, limit)
	if _, err := fm.Read(mcontent); err != nil {
		log.Println("unable to read meta file: %v", err)
		return nil, err
	}

	if !strings.Contains(string(mcontent), id) {
		log.Println("The id %v does not have permission to read the file %v", id, path)
		return nil, errors.New("Permission Denied")
	}

	f, err := os.Open(s.root + path)
	if err != nil {
		log.Println("unable to open file %v", path)
		log.Println(err)
		return nil, err
	}

	stat, err := f.Stat()

	if err != nil {
		log.Println("unable to get file stat")
		return nil, err
	}

	if int32(stat.Size()) <= offset {
		log.Println("the file size is smaller than or equal to the offset")
		return nil, errors.New("the file size is smaller than or equal to the offset")
	}

	if int32(stat.Size()) < limit {
		limit = int32(stat.Size())
	}

	content := make([]byte, limit)

	if _, err := f.ReadAt(content, int64(offset)); err != nil {
		log.Println("unable to read file chunck: ", err)
		return nil, err
	}

	return content, nil
}
