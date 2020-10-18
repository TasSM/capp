package server

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/TasSM/appCache/service"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedArrayBasedCacheServer
}

func (s *server) CreateRecord(ctx context.Context, in *pb.CreateRecordRequest) (*pb.CreateRecordResponse, error) {
	key, ttl := in.GetKey(), in.GetTtl()
	log.Printf("Received Record Create Request: %v | %v", key, ttl)
	//Validate and create the record
	return &pb.CreateRecordResponse{Key: in.GetKey(), Ttl: in.GetTtl()}, nil
}

func (s *server) StoreMessage(ctx context.Context, in *pb.AppendRecordRequest) (*pb.AppendRecordResponse, error) {
	log.Printf("Received Message Request for server: %v", in.GetKey())
	//Store the Message
	return &pb.AppendRecordResponse{Status: true}, nil
}

func (s *server) GetStatistics(ctx context.Context, in *pb.Empty) (*pb.StatisticResponse, error) {
	log.Printf("Received Statistics Request")
	// Retrieve statistics
	update := time.Now().Unix()
	return &pb.StatisticResponse{ServerCount: 5, UserCount: 6, LastUpdate: update}, nil
}

func ServeGRPC(port string, cp *redis.Pool) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to start GRPC server: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterArrayBasedCacheServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
