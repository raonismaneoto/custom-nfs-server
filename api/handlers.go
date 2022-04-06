package api

import (
	context "context"
	"io"
	"log"

	"github.com/raonismaneoto/custom-nfs-server/server"
)

type Handler struct {
	s *server.Server
}

func (h *Handler) Save(srv NFSS_SaveServer) (error) {
	log.Println("Save call received.")
	ctx := srv.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		req, err := srv.Recv()
		if err == io.EOF {
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		err = h.s.Save(req.Content)
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}
	}
}

func (h *Handler) Ping(ctx context.Context, request *Empty) (*Empty, error) {
	log.Println("Ping received.")
	return &Empty{}, nil
}

func (h *Handler) Mount(ctx context.Context, request *MountRequest) (*MountResponse, error) {
	log.Println("Mount received.")
	return &MountResponse{}, nil
}

func (h *Handler) UnMount(ctx context.Context, request *UnMountRequest) (*Empty, error) {
	log.Println("UnMount received.")
	return &Empty{}, nil
}

func (h *Handler) Read(request *ReadRequest, srv NFSS_ReadServer) (error) {
	log.Println("Read received.")
	return srv.Send(&ReadResponse{})
}