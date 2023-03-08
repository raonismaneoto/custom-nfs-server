package client

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/raonismaneoto/custom-nfs-server/nfs-server/api"
	"google.golang.org/grpc"
)

type Client struct {
	address    string
	connection *grpc.ClientConn
}

func NewClient(address string) *Client {
	return &Client{address: address}
}

func (c *Client) Save(id, path string, content []byte) error {
	lc := c.getGrpcClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*30)
	defer cancel()

	client, err := lc.Save(ctx)
	if err != nil {
		log.Println(err.Error())
	}

	select {
	case <-client.Context().Done():
		return client.Context().Err()
	default:
		req := api.SaveRequest{
			Id:      id,
			Path:    path,
			Content: content,
		}
		log.Println("content size in client")
		log.Println(len(req.Content))

		if err := client.Send(&req); err != nil {
			log.Printf("send error %v", err)
		}

		if _, err := client.CloseAndRecv(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) SaveAsync(id, path string, content <-chan []byte, proceed chan<- string) error {
	lc := c.getGrpcClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*30)
	defer cancel()

	client, err := lc.Save(ctx)
	if err != nil {
		log.Println(err.Error())
	}

	for {
		select {
		case <-client.Context().Done():
			return client.Context().Err()
		default:
		}

		proceed <- "proceed"
		currContent, ok := <-content
		if !ok {
			log.Println("entered not ok")
			if _, err := client.CloseAndRecv(); err != nil {
				return err
			}
			close(proceed)
			return nil
		}

		req := api.SaveRequest{
			Id:      id,
			Path:    path,
			Content: currContent,
		}
		log.Println("content size in client")
		log.Println(len(req.Content))

		if err := client.Send(&req); err != nil {
			log.Printf("send error %v", err)
		}
	}
}

func (c *Client) Read(id, path string, content chan<- []byte, proceed <-chan string) error {
	lc := c.getGrpcClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*30)
	defer cancel()

	srv, err := lc.Read(ctx, &api.ReadRequest{Path: path, Id: id})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for {
		select {
		case <-srv.Context().Done():
			return srv.Context().Err()
		default:
		}

		<-proceed
		data, err := srv.Recv()
		if err == io.EOF {
			close(content)
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			return err
		}

		content <- data.Content
	}
}

func (c *Client) Mount(id, path string) ([]byte, error) {
	lc := c.getGrpcClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	response, err := lc.Mount(ctx, &api.MountRequest{Path: path, Id: id})

	if err != nil {
		return nil, err
	}

	return response.MetaData, nil
}

func (c *Client) Chpem(ownerId, user, path, op string) error {
	lc := c.getGrpcClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	_, err := lc.Chpem(ctx, &api.ChpemRequest{OwnerId: ownerId, User: user, Path: path, Op: op})

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) getGrpcClient() api.NFSSClient {
	if c.connection == nil {
		log.Println("Starting grpc connection")
		conn, err := grpc.Dial(c.address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		c.connection = conn
	}

	return api.NewNFSSClient(c.connection)
}
