package server

import (
	"log"
)

type Server struct {

}

func (s *Server) Save(content []byte) (error) {
	log.Println("Save call received.")
	log.Println("saving")
	return nil
}