package cache

import (
	"log"
	"time"

	"github.com/TasSM/appCache/defs"
	"github.com/gomodule/redigo/redis"
)

type client struct {
	cp *redis.Pool
}

func createConnPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func NewCacheClient(addr string) defs.CacheClientService {
	return &client{
		cp: createConnPool(addr),
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

func (c *client) CreateCacheArrayRecord(key string, expiry int64) error {
	conn := c.cp.Get()
	defer conn.Close()
	conn.Send("LPUSH", key, "BEGIN")
	conn.Send("EXPIREAT", key, expiry)
	conn.Flush()
	conn.Receive()
	if _, err := conn.Receive(); err != nil {
		return err
	}
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
