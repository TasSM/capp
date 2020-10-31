package svcclient

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/TasSM/appCache/svcgrpc"
	"google.golang.org/grpc"
)

const (
	defaultTimeout = time.Second * 10
)

type GrpcService struct {
	grpcClient svcgrpc.ArrayBasedCacheClient
}

func ConnectGRPCService(connString string) (*GrpcService, error) {
	conn, err := grpc.Dial(connString, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &GrpcService{grpcClient: svcgrpc.NewArrayBasedCacheClient(conn)}, nil
}

func (s *GrpcService) CreateRecord(key string, ttl int32) (*svcgrpc.CreateRecordResponse, error) {
	req := &svcgrpc.CreateRecordRequest{
		Key: key,
		Ttl: ttl,
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	resp, err := s.grpcClient.CreateRecord(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *GrpcService) StoreMessage(key string, message string) (*svcgrpc.AppendRecordResponse, error) {
	req := &svcgrpc.AppendRecordRequest{
		Key:     key,
		Message: message,
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	resp, err := s.grpcClient.StoreMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *GrpcService) GetStatistics() (*svcgrpc.StatisticResponse, error) {
	req := &svcgrpc.Empty{}
	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	resp, err := s.grpcClient.GetStatistics(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *GrpcService) GetRecord(key string) (svcgrpc.ArrayBasedCache_GetRecordClient, error) {
	req := &svcgrpc.GetRecordRequest{Key: key}
	ctx := context.Background()
	resp, err := s.grpcClient.GetRecord(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *GrpcService) StreamToArray(stream svcgrpc.ArrayBasedCache_GetRecordClient) ([]string, error) {
	out := []string{}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.New("Unexpected error reading stream")
		}
		out = append(out, msg.Message)
	}
	return out, nil
}
