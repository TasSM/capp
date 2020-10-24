package main

import (
	"log"

	"github.com/TasSM/appCache/service/svcgrpc/client"
)

const (
	address = "localhost:9099"
)

// test the grpc client
func main() {
	c, err := client.ConnectGRPCService(address)
	if err != nil {
		log.Printf("Unable to connect to the GRPC server at %v", address)
		return
	}

	testKey, testTTL := "abcdefgh", 1234

	r, err := c.CreateRecord(testKey, int32(testTTL))
	if err != nil {
		log.Printf("Creation of record failed - exiting")
		return
	}
	log.Printf("record created with key %v and ttl %d", r.Key, r.Ttl)
	return
}
