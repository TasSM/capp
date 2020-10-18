package cache

import (
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

func CreateConnectionPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 240,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

func KeyExists(key string, cp *redis.Pool) bool {
	conn := cp.Get()
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

func CreateCacheArrayRecord(key string, expiry int64, cp *redis.Pool) bool {
	conn := cp.Get()
	defer conn.Close()
	conn.Send("LPUSH", key, "BEGIN")
	conn.Send("EXPIREAT", key, expiry)
	conn.Flush()
	conn.Receive()
	if _, err := conn.Receive(); err != nil {
		panic(err)
	}
	return true
}

// func to create new Redis record - called from a different API route
func StartCacheLoop(key string, expiry int64, cp *redis.Pool, dc chan string) {
	for {
		select {
		case msg := <-dc:
			if time.Now().Unix() > expiry {
				log.Printf("INFO - Closing cache connection for expired server %s", key)
				return
			}
			conn := cp.Get()
			if _, err := conn.Do("RPUSH", key, msg); err != nil {
				panic(err)
				log.Printf("ERROR - Writing message to key %s failed", key)
			}
			conn.Close()
		}
	}
}
