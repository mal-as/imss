package router

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/cespare/xxhash"
	"github.com/mal-as/imss/proto"
	"google.golang.org/grpc"
)

var (
	errWrongShardNumber = errors.New("wrong shard number")
	errEmptyShards      = errors.New("empty shards list")
)

type Router struct {
	sync.RWMutex
	srv    *grpc.Server
	shards []proto.StorageClient
	n      uint64
}

func Init(shards []string, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	r := &Router{
		srv: grpc.NewServer(),
	}
	proto.RegisterStorageServer(r.srv, r)
	go r.srv.Serve(listen)

	return r.addShards(shards)
}

func (r *Router) addShards(shards []string) error {
	var broken []string
	var errs []error
	var wg sync.WaitGroup

	for _, shard := range shards {
		wg.Add(1)
		go func(shard string) {
			defer wg.Done()
			if err := r.addShard(shard); err != nil {
				r.Lock()
				broken = append(broken, shard)
				errs = append(errs, err)
				r.Unlock()
			}
		}(shard)
	}
	wg.Wait()
	if len(broken) > 0 {
		return fmt.Errorf("coudn't connect to %v; reasons: %v", broken, errs)
	}

	return nil
}

func (r *Router) addShard(address string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}

	r.Lock()
	r.shards = append(r.shards, proto.NewStorageClient(conn))
	r.n++
	r.Unlock()

	return nil
}

func (r *Router) Get(ctx context.Context, req *proto.GetReq) (*proto.Respone, error) {
	r.RLock()
	if r.n == 0 {
		r.RUnlock()
		return nil, errEmptyShards
	}
	num := xxhash.Sum64([]byte(req.Key)) % r.n
	r.RUnlock()
	if len(r.shards) < int(num) {
		return nil, fmt.Errorf("%w: %d", errWrongShardNumber, num)
	}
	return r.shards[num].Get(ctx, req)
}

func (r *Router) Store(ctx context.Context, req *proto.StoreReq) (*proto.StoreStatus, error) {
	r.RLock()
	if r.n == 0 {
		r.RUnlock()
		return nil, errEmptyShards
	}
	num := xxhash.Sum64([]byte(req.Key)) % r.n
	r.RUnlock()
	if len(r.shards) < int(num) {
		return nil, fmt.Errorf("%w: %d", errWrongShardNumber, num)
	}
	return r.shards[num].Store(ctx, req)
}
