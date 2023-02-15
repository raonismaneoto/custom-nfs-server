package client

import (
	"context"
	"log"
	"time"

	"github.com/raonismaneoto/custom-nfs-server/nfs-server/api"
	"google.golang.org/grpc"
)

type Client struct {
	address string
}

func NewClient(address string) *Client {
	return &Client{address: address}
}

func (c *Client) Mount(id, path string) (*api.MountResponse, error) {
	lc, conn := grpcClient(c.address)
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	response, err := lc.Mount(ctx, &api.MountRequest{Path: path, Id: id})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func grpcClient(address string) (api.NFSSClient, *grpc.ClientConn) {
	log.Println("Starting grpc connection")
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	log.Println("Grpc connection started.")
	return api.NewNFSSClient(conn), conn
}
