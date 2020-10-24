package service

import (
	"log"
	"time"

	"github.com/TasSM/appCache/defs"
	"github.com/gomodule/redigo/redis"
)

type client struct {
	cp *redis.Pool
}

func NewCacheClient(addr string) defs.CacheClientService {
	return &client{
		cp: &redis.Pool{
			MaxIdle:     5,
			IdleTimeout: 240,
			Dial: func() (redis.Conn, error) {
				conn, err := redis.Dial("tcp", addr)
				if err != nil {
					log.Printf("Failed to dial redis host at %s", addr)
					panic(err)
				}
				log.Printf("Successfully dialed redis host at %s", addr)
				return conn, nil
			},
		},
	}
}

func (c *client) KeyExists(key string) bool {
	conn := c.cp.Get()
	defer conn.Close()
	res, err := redis.Int(conn.Do("EXISTS", key))
	if err != nil {
		panic(err)
	}
	if res == 1 {
		return true
	}
	return false
}

func (c *client) CreateCacheArrayRecord(key string, ttl int64) error {
	conn := c.cp.Get()
	defer conn.Close()

	//rewrite this as a transaction with MULTI

	conn.Send("LPUSH", key, "BEGIN")
	conn.Send("EXPIREAT", key, ttl)
	conn.Flush()
	conn.Receive()
	if _, err := conn.Receive(); err != nil {
		log.Printf("Received Error Status")
		return err
	}
	log.Printf("Created record with key %v and ttl %d", key, ttl)
	return nil
}

// func to create new Redis record - called from a different API route
func (c *client) Start(key string, expiry int64, dc chan string) {
	for {
		select {
		case msg := <-dc:
			if time.Now().Unix() > expiry {
				log.Printf("INFO - Closing cache connection for expired server %s", key)
				return
			}
			conn := c.cp.Get()
			if _, err := conn.Do("RPUSH", key, msg); err != nil {
				panic(err)
				log.Printf("ERROR - Writing message to key %s failed", key)
			}
			conn.Close()
		}
	}
}
