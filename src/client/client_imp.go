package client

import (
	"context"
	"errors"
	"io"
	"time"

	pb "github.com/TasSM/appCache/svcgrpc"
	"google.golang.org/grpc"
)

const (
	defaultTimeout = time.Second * 10
)

type GrpcService struct {
	grpcClient pb.ArrayBasedCacheClient
}

func ConnectGRPCService(connection string) (*GrpcService, error) {
	conn, err := grpc.Dial(connection, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &GrpcService{grpcClient: pb.NewArrayBasedCacheClient(conn)}, nil
}

func (s *GrpcService) CreateRecord(key string, ttl int32) (*pb.CreateRecordResponse, error) {
	req := &pb.CreateRecordRequest{
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

func (s *GrpcService) StoreMessage(key string, message string) (*pb.AppendRecordResponse, error) {
	req := &pb.AppendRecordRequest{
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

func (s *GrpcService) GetStatistics() (*pb.StatisticResponse, error) {
	req := &pb.Empty{}
	ctx, cancelFunc := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	resp, err := s.grpcClient.GetStatistics(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *GrpcService) GetRecord(key string) (pb.ArrayBasedCache_GetRecordClient, error) {
	req := &pb.GetRecordRequest{Key: key}
	ctx := context.Background()
	resp, err := s.grpcClient.GetRecord(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *GrpcService) StreamToArray(stream pb.ArrayBasedCache_GetRecordClient) ([]string, error) {
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
