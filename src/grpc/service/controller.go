package service

import (
	"context"
	"errors"
	"time"

	"github.com/TasSM/appCache/defs"
	"github.com/TasSM/appCache/grpc"
)

type cacheClientController struct {
	client        defs.CacheClientService
	inputChannels map[string](chan string)
	grpc.UnimplementedArrayBasedCacheServer
}

func NewCacheClientController(cacheClient defs.CacheClientService) grpc.ArrayBasedCacheServer {
	return &cacheClientController{
		client:        cacheClient,
		inputChannels: make(map[string](chan string)),
	}
}

func (ctlr *cacheClientController) CreateRecord(ctx context.Context, req *grpc.CreateRecordRequest) (*grpc.CreateRecordResponse, error) {
	key, ttl := req.GetKey(), req.GetTtl()
	if ctlr.client.KeyExists(key) == true {
		return &grpc.CreateRecordResponse{Key: key, Ttl: ttl}, errors.New("key in use")
	}
	err := ctlr.client.CreateCacheArrayRecord(key, ttl)
	if err != nil {
		return &grpc.CreateRecordResponse{Key: key, Ttl: ttl}, err
	}
	ctlr.inputChannels[key] = make(chan string, 24)
	go ctlr.client.Start(key, ttl, ctlr.inputChannels[key])
	go func() {
		time.Sleep(time.Duration(ttl) * time.Second)
		delete(ctlr.inputChannels, key)
	}()
	return &grpc.CreateRecordResponse{Key: key, Ttl: ttl}, nil
}

func (ctlr *cacheClientController) StoreMessage(ctx context.Context, req *grpc.AppendRecordRequest) (*grpc.AppendRecordResponse, error) {
	key, msg := req.GetKey(), req.GetMessage()
	if ctlr.inputChannels[key] == nil {
		return &grpc.AppendRecordResponse{Status: false}, errors.New("unable to write to input channel")
	}
	ctlr.inputChannels[key] <- msg
	return &grpc.AppendRecordResponse{Status: true}, nil
}

func (ctlr *cacheClientController) GetStatistics(ctx context.Context, req *grpc.Empty) (*grpc.StatisticResponse, error) {
	return &grpc.StatisticResponse{ServerCount: int32(len(ctlr.inputChannels)), UserCount: int32(10), LastUpdate: int64(time.Now().Unix())}, nil
}
