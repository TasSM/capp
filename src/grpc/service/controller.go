package service

import (
	"context"
	"errors"

	"github.com/TasSM/appCache/defs"
	"github.com/TasSM/appCache/grpc"
)

type cacheClientController struct {
	client defs.CacheClientService
}

func NewCacheClientController(cacheClient defs.CacheClientService) grpc.ArrayBasedCacheServer {
	return &cacheClientController{
		client: cacheClient,
	}
}

func (ctlr *cacheClientController) CreateRecord(ctx context.Context, req *grpc.CreateRecordRequest) (*grpc.CreateRecordResponse, error) {
	key, ttl := req.GetKey(), req.GetTtl()
	if ctlr.client.KeyExists(key) == true {
		return &grpc.CreateRecordResponse{Key: key, Ttl: ttl}, errors.New("key in use")
	}
	// sub with the _, err := format
	ctlr.client.CreateCacheArrayRecord(key, ttl)
}
