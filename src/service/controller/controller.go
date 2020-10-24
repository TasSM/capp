package controller

import (
	"context"
	"errors"
	"time"

	"github.com/TasSM/appCache/defs"
	"github.com/TasSM/appCache/service/svcgrpc"
)

type cacheClientController struct {
	client        defs.CacheClientService
	inputChannels map[string](chan string)
	svcgrpc.UnimplementedArrayBasedCacheServer
}

func NewCacheClientController(cacheClient defs.CacheClientService) svcgrpc.ArrayBasedCacheServer {
	return &cacheClientController{
		client:        cacheClient,
		inputChannels: make(map[string](chan string)),
	}
}

func (ctlr *cacheClientController) CreateRecord(ctx context.Context, req *svcgrpc.CreateRecordRequest) (*svcgrpc.CreateRecordResponse, error) {
	key, ttl := req.GetKey(), req.GetTtl()
	if ctlr.client.KeyExists(key) == true {
		return &svcgrpc.CreateRecordResponse{Key: key, Ttl: ttl}, errors.New("key in use")
	}
	err := ctlr.client.CreateCacheArrayRecord(key, int64(ttl))
	if err != nil {
		return nil, err
	}
	ctlr.inputChannels[key] = make(chan string, 24)
	go ctlr.client.Start(key, int64(ttl), ctlr.inputChannels[key])
	go func() {
		time.Sleep(time.Duration(ttl) * time.Second)
		delete(ctlr.inputChannels, key)
	}()
	return &svcgrpc.CreateRecordResponse{Key: key, Ttl: ttl}, nil
}

func (ctlr *cacheClientController) StoreMessage(ctx context.Context, req *svcgrpc.AppendRecordRequest) (*svcgrpc.AppendRecordResponse, error) {
	key, msg := req.GetKey(), req.GetMessage()
	if ctlr.inputChannels[key] == nil {
		return &svcgrpc.AppendRecordResponse{Status: false}, errors.New("unable to write to input channel")
	}
	ctlr.inputChannels[key] <- msg
	return &svcgrpc.AppendRecordResponse{Status: true}, nil
}

func (ctlr *cacheClientController) GetStatistics(ctx context.Context, req *svcgrpc.Empty) (*svcgrpc.StatisticResponse, error) {
	return &svcgrpc.StatisticResponse{ServerCount: int32(len(ctlr.inputChannels)), UserCount: int32(10), LastUpdate: int64(time.Now().Unix())}, nil
}
