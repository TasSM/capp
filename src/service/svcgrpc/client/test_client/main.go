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

	testKey, testTTL := "test", 20

	//test creation
	r, err := c.CreateRecord(testKey, int32(testTTL))
	if err != nil {
		log.Printf("Creation of record failed - exiting")
		return
	}
	log.Printf("record created with key %v and ttl %d", r.Key, r.Ttl)

	// test writing messages
	messages := []string{"abc", "def", "hij"}
	for i, m := range messages {
		if _, err := c.StoreMessage(testKey, m); err != nil {
			log.Printf("Write operation for %d, %v failed", i, m)
			return
		}
		log.Printf("Writing message to record %v: %v", testKey, m)
	}

	// get statistics
	s, err := c.GetStatistics()
	if err != nil {
		log.Printf("Retrieval of statistics failed")
		return
	}
	log.Printf("Statistics received: %v", s)

	stream, err := c.GetRecord(testKey)
	if err != nil {
		log.Printf("Retrieving record stream %v failed", testKey)
	}
	a, err := c.StreamToArray(stream)
	if err != nil {
		log.Printf("Failed to unmarshall message stream to array")
	}
	log.Printf("Retrieved messages %v", a)

	return
}
