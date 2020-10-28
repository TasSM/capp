package main

import (
	"log"
	"net"

	"github.com/TasSM/appCache/api"
	"github.com/TasSM/appCache/service"
	"github.com/TasSM/appCache/service/controller"
	"github.com/TasSM/appCache/svcgrpc"
	"github.com/TasSM/appCache/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	REDIS_HOST = "REDIS_HOST"
	GRPC_ADDR  = "GRPC_ADDR"
	METHOD     = "tcp"
)

func main() {

	cacheService := service.NewCacheClient(util.GetEnv(REDIS_HOST, "apps.labnet:6379"))
	cacheServiceController := controller.NewCacheClientController(cacheService)

	server := grpc.NewServer()
	svcgrpc.RegisterArrayBasedCacheServer(server, cacheServiceController)
	reflection.Register(server)

	grpcPort := util.GetEnv(GRPC_ADDR, ":9099")
	con, err := net.Listen(METHOD, grpcPort)
	if err != nil {
		panic(err)
	}

	// Serve HTTP Routes
	go api.ServeRoutes("8080", cacheService)

	// Serve GRPC Routes
	log.Printf("Starting GRPC server on port: %s", grpcPort)
	err = server.Serve(con)
	if err != nil {
		panic(err)
	}
}
