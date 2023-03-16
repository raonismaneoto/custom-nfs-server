package api

import (
	context "context"
	"encoding/json"
	"io"
	"log"

	"github.com/raonismaneoto/custom-nfs-server/nfs-server/server"
)

const MaxBytesPerResponse int32 = 10000

type Handler struct {
	s *server.Server
}

func New(server *server.Server) *Handler {
	return &Handler{s: server}
}

func (h *Handler) SaveAsync(srv NFSS_SaveAsyncServer) error {
	log.Println("Save call received.")
	ctx := srv.Context()

	content := make(chan []byte, 20)
	errors := make(chan error)

	req, err := srv.Recv()
	if err != nil {
		return err
	}
	go h.s.SaveAsync(req.Id, req.Path, content, errors)
	log.Println("putting content in the channel")
	log.Println(content)
	content <- req.Content
	log.Println("content inserted, going to start the loop")
	for {
		log.Println("entering the loop")
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errors:
			log.Println("got an error")
			log.Println(err.Error())
			return err
		default:
		}

		req, err := srv.Recv()
		if err == io.EOF {
			log.Println("exit")
			if err = srv.SendAndClose(&Empty{}); err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			return err
		}
		log.Println("going to send content in the channel inside loop")
		content <- req.Content
	}
}

func (h *Handler) Save(ctx context.Context, request *SaveRequest) (*Empty, error) {
	if len(request.Content) == 0 {
		err := h.s.Mkdir(request.Id, request.Path)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		return &Empty{}, nil
	}

	err := h.s.Save(request.Id, request.Path, request.Content)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &Empty{}, nil
}

func (h *Handler) Ping(ctx context.Context, request *Empty) (*Empty, error) {
	log.Println("Ping received.")
	return &Empty{}, nil
}

func (h *Handler) Mount(ctx context.Context, request *MountRequest) (*MountResponse, error) {
	log.Println("Mount received.")
	metadata, err := h.s.GetMetaData(request.Id, request.Path)
	if err != nil {
		return nil, err
	}
	serializedMd, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}
	return &MountResponse{MetaData: serializedMd}, nil
}

func (h *Handler) UnMount(ctx context.Context, request *UnMountRequest) (*Empty, error) {
	log.Println("UnMount received.")
	return &Empty{}, nil
}

func (h *Handler) Read(request *ReadRequest, srv NFSS_ReadServer) error {
	ctx := srv.Context()

	content := make(chan []byte)
	errors := make(chan error)

	go h.s.Read(request.Id, request.Path, content, errors)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case currContent, ok := <-content:
			if !ok {
				return nil
			}
			resp := ReadResponse{
				Content: currContent,
			}

			if err := srv.Send(&resp); err != nil {
				log.Printf("send error %v", err)
			}
		case err := <-errors:
			log.Println(err.Error())
			return err
		}
	}
}

func (h *Handler) Remove(ctx context.Context, request *RemoveRequest) (*Empty, error) {
	log.Println("Remove received.")
	err := h.s.Rm(request.Id, request.Path)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &Empty{}, nil
}

func (h *Handler) Chpem(ctx context.Context, request *ChpemRequest) (*Empty, error) {
	err := h.s.Chpem(request.OwnerId, request.User, request.Path, request.Op)
	if err != nil {
		return nil, err
	}
	return &Empty{}, nil
}

func (h *Handler) mustEmbedUnimplementedNFSSServer() {

}
