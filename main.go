package main

import (
	"log"
	"net"
	"os"

	"github.com/raonismaneoto/custom-nfs-server/api"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("port")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	s := grpc.NewServer()
	api.RegisterNFSSServer(s, &api.Handler{})

	log.Println("NodeServer listening at %v", lis.Addr())

	log.Println("going to start grpc NodeServer listener")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}