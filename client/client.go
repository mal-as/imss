package client

import (
	"context"
	"errors"

	"github.com/mal-as/imss/proto"
	"google.golang.org/grpc"
)

var errBadStatus = errors.New("returned bad status")

type Client struct {
	router proto.StorageClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{router: proto.NewStorageClient(conn)}, nil
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	resp, err := c.router.Get(ctx, &proto.GetReq{Key: key})
	if err != nil {
		return nil, err
	}
	return resp.Value, nil
}

func (c *Client) Store(ctx context.Context, key string, value []byte) error {
	resp, err := c.router.Store(ctx, &proto.StoreReq{Key: key, Value: value})
	if err != nil {
		return err
	}
	if resp.Status == proto.StoreStatus_Bad {
		return errBadStatus
	}
	return nil
}
