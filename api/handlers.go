package api

import (
	"bytes"
	context "context"
	"encoding/json"
	"io"
	"log"

	"github.com/raonismaneoto/custom-nfs-server/models"
	"github.com/raonismaneoto/custom-nfs-server/server"
)

const MaxBytesPerResponse int32 = 10000

type Handler struct {
	s *server.Server
}

func New(server *server.Server) *Handler{
	return &Handler{s : server}
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

		err = h.s.Save(req.Id, req.Path ,req.Content)
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
	ctx := srv.Context()

	var fm models.Metadata
	err := json.NewDecoder(bytes.NewReader(request.MetaData)).Decode(&fm)
	if err != nil {

	}

	chuncks := fm.Size / MaxBytesPerResponse

	for i := int32(0); i <= chuncks; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		content, err := h.s.Read(request.Id, fm.Path, i*MaxBytesPerResponse, MaxBytesPerResponse)
		if err != nil {
			log.Println("error while reading file content: %v", err)
			return err
		}
		
		resp := ReadResponse{
			Content: content,
		}

		if err := srv.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
	}

	return nil
}

func (h *Handler) Remove(ctx context.Context, request *RemoveRequest) (*Empty, error) {
	log.Println("Ping received.")
	return &Empty{}, nil
}