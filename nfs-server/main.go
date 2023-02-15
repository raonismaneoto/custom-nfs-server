package main

import (
	"log"
	"net"
	"os"

	"github.com/raonismaneoto/custom-nfs-server/nfs-server/api"
	"github.com/raonismaneoto/custom-nfs-server/nfs-server/server"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("port")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	gs := grpc.NewServer()
	ns := server.New("/home/raonismaneoto/custom-nfs/")
	api.RegisterNFSSServer(gs, api.New(ns))

	log.Println("NFSS listening at %v", lis.Addr())

	log.Println("going to start grpc Server listener")
	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
