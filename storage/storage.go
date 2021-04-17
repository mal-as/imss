package storage

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/mal-as/imss/proto"
	"google.golang.org/grpc"
)

type Storage struct {
	sync.Mutex
	data map[string][]byte
	srv  *grpc.Server
}

func Init(port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	s := &Storage{
		data: make(map[string][]byte),
		srv:  grpc.NewServer(),
	}
	proto.RegisterStorageServer(s.srv, s)
	go s.srv.Serve(listen)

	return nil
}

func (s *Storage) Get(ctx context.Context, req *proto.GetReq) (*proto.Respone, error) {
	s.Lock()
	value := s.data[req.Key]
	s.Unlock()
	fmt.Println(string(value))
	return &proto.Respone{Value: value}, nil
}

func (s *Storage) Store(ctx context.Context, req *proto.StoreReq) (*proto.StoreStatus, error) {
	s.Lock()
	s.data[req.Key] = req.Value
	s.Unlock()
	return &proto.StoreStatus{Status: proto.StoreStatus_Ok}, nil
}
