package main

import (
	"log"
	"time"

	client "github.com/TasSM/appCache/svcclient"
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

	tk1, tttl1 := "testalpha", 20
	tk2, tttl2 := "testbeta", 32

	//test creation
	r1, err := c.CreateRecord(tk1, int32(tttl1))
	if err != nil {
		log.Printf("Creation of record failed - exiting")
		return
	}
	log.Printf("record created with key %v and ttl %d", r1.Key, r1.Ttl)

	r2, err := c.CreateRecord(tk2, int32(tttl2))
	if err != nil {
		log.Printf("Creation of record failed - exiting")
		return
	}
	log.Printf("record created with key %v and ttl %d", r2.Key, r2.Ttl)

	// test writing messages
	messages := []string{"abc", "def", "hij"}
	for i, m := range messages {
		if _, err := c.StoreMessage(tk1, m); err != nil {
			log.Printf("Write operation for %d, %v failed", i, m)
			return
		}
		log.Printf("Writing message to record %v: %v", tk1, m)

		if _, err := c.StoreMessage(tk2, m); err != nil {
			log.Printf("Write operation for %d, %v failed", i, m)
			return
		}
		log.Printf("Writing message to record %v: %v", tk2, m)
	}

	// get statistics
	s1, err := c.GetStatistics()
	if err != nil {
		log.Printf("Retrieval of statistics failed")
		return
	}
	log.Printf("Statistics received: %v", s1)

	time.Sleep(1 * time.Second)
	stream, err := c.GetRecord(tk1)
	if err != nil {
		log.Printf("Retrieving record stream %v failed", tk1)
	}
	arr, err := c.StreamToArray(stream)
	if err != nil {
		log.Printf("Failed to unmarshall message stream to array")
	}
	log.Printf("Retrieved messages %v", arr)

	time.Sleep(3 * time.Second)

	s2, err := c.GetStatistics()
	if err != nil {
		log.Printf("Retrieval of statistics failed")
		return
	}
	log.Printf("Statistics received: %v", s2)
	return
}
