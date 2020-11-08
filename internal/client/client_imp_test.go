package client

import (
	"context"
	"log"
	"net"
	"sync"
	"testing"

	"github.com/TasSM/capp/internal/controller"
	"github.com/TasSM/capp/internal/service"
	pb "github.com/TasSM/capp/internal/svcgrpc"
	"github.com/alicebob/miniredis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var testKey string = "test12s"
var testTTL int32 = 122
var testMessage string = "hello, testing"

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	redisServer, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	cacheClient := service.NewCacheClient(redisServer.Addr())
	pb.RegisterArrayBasedCacheServer(server, controller.NewCacheClientController(cacheClient))

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func connectTestClient(ctx context.Context) *GrpcService {
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		panic(err)
	}
	return &GrpcService{grpcClient: pb.NewArrayBasedCacheClient(conn)}
}

func initializeTestRecord(c *GrpcService, t *testing.T) {
	res, err := c.CreateRecord(testKey, testTTL)
	if err != nil {
		t.Errorf(err.Error())
	}
	if res.Key != testKey {
		t.Errorf("invalid key in create record response")
	}
	if res.Ttl != testTTL {
		t.Errorf("invalid ttl in create record response")
	}
}

func TestRecordCreation(t *testing.T) {
	ctx := context.Background()
	client := connectTestClient(ctx)
	initializeTestRecord(client, t)
}

func TestGetInvalid(t *testing.T) {
	ctx := context.Background()
	client := connectTestClient(ctx)
	res, err := client.GetRecord(testKey)
	if err != nil {
		t.Errorf(err.Error())
	}
	if _, err := client.StreamToArray(res); err == nil {
		t.Errorf("Expected error to occur converting invalid stream to array")
	}
}

func TestGetValid(t *testing.T) {
	ctx := context.Background()
	client := connectTestClient(ctx)
	initializeTestRecord(client, t)
	get, err := client.GetRecord(testKey)
	if err != nil {
		t.Errorf(err.Error())
	}
	arr, err := client.StreamToArray(get)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(arr) != 0 {
		t.Errorf("Response should be empty for invalid key")
	}
}

func TestStoreMessageInvalidKey(t *testing.T) {
	ctx := context.Background()
	client := connectTestClient(ctx)
	if _, err := client.StoreMessage(testKey, testMessage); err == nil {
		t.Errorf("Expected error storing message at invalid key")
	}
}

func TestStoreMessageValidKey(t *testing.T) {
	ctx := context.Background()
	client := connectTestClient(ctx)
	initializeTestRecord(client, t)
	res, err := client.StoreMessage(testKey, testMessage)
	if err != nil {
		t.Errorf(err.Error())
	}
	if res.Status != true {
		t.Errorf("Message should be stored with success status")
	}
}

func TestGetStatistics(t *testing.T) {
	ctx := context.Background()
	client := connectTestClient(ctx)
	res, err := client.GetStatistics()
	if err != nil {
		t.Errorf(err.Error())
	}
	if res.RecordCount != 0 {
		t.Errorf("Expected record count = 0")
	}
	if res.ActiveConnections != 1 {
		t.Errorf("Expected active connections = 1")
	}
	if res.LastUpdate == "" {
		t.Errorf("Expected last update timestamp to exist")
	}
}

func TestE2EFlow(t *testing.T) {
	ctx := context.Background()
	client := connectTestClient(ctx)
	initializeTestRecord(client, t)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := client.StoreMessage(testKey, testMessage); err != nil {
			t.Errorf(err.Error())
		}
	}()
	wg.Wait()
	res, err := client.GetRecord(testKey)
	if err != nil {
		t.Errorf(err.Error())
	}
	arr, err := client.StreamToArray(res)
	if err != nil {
		t.Errorf(err.Error())
	}
	if len(arr) != 1 && arr[0] != testMessage {
		t.Errorf("Invalid messages retrieved")
	}
	stats, err := client.GetStatistics()
	if err != nil {
		t.Errorf(err.Error())
	}
	if stats.RecordCount != 1 {
		t.Errorf("Expected record count = 0")
	}
	if stats.ActiveConnections != 1 {
		t.Errorf("Expected active connections = 1")
	}
	if stats.LastUpdate == "" {
		t.Errorf("Expected last update timestamp to exist")
	}
}
