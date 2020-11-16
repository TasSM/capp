package service

import (
	"errors"
	"log"
	"time"

	"github.com/TasSM/capp/internal/defs"
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
					log.Printf("ERROR - Failed to dial redis host at %s", addr)
					panic(err)
				}
				log.Printf("INFO - Dialed redis host at %v", addr)
				return conn, nil
			},
		},
	}
}

func (c *client) Ping() error {
	conn := c.cp.Get()
	defer c.cp.Close()
	_, err := conn.Do("PING")
	if err != nil {
		return err
	}
	return nil
}

func (c *client) CreateCacheArrayRecord(key string, ttl int64) error {
	if check := c.KeyExists(key); check == true {
		return errors.New("A record with this key already exists")
	}
	conn := c.cp.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("LPUSH", key, "BEGIN")
	conn.Send("EXPIRE", key, ttl)
	if _, err := conn.Do("EXEC"); err != nil {
		log.Printf("ERROR - Received Error Status")
		return err
	}
	return nil
}

func (c *client) DisposePool() {
	c.cp.Close()
}

func (c *client) GetStatistics() (*defs.StatisticResponse, error) {
	conn := c.cp.Get()
	defer conn.Close()
	rc, err := redis.Values(conn.Do("KEYS", "*"))
	if err != nil {
		return nil, err
	}
	return &defs.StatisticResponse{
		RecordCount:       len(rc),
		ActiveConnections: c.cp.ActiveCount(),
		Timestamp:         time.Now().Format(time.RFC3339),
	}, nil
}

func (c *client) GetTTL(key string) (int, error) {
	conn := c.cp.Get()
	defer conn.Close()
	ttl, err := redis.Int(conn.Do("TTL", key))
	if err != nil {
		return -2, err
	}
	return ttl, nil
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

func (c *client) ReadArrayRecord(key string) ([]string, error) {
	conn := c.cp.Get()
	defer conn.Close()
	res, err := redis.Strings(conn.Do("LRANGE", key, 1, -1))
	if err != nil {
		log.Printf("ERROR - unable to read record %v", key)
		return nil, err
	}
	return res, nil
}

func (c *client) Start(key string, expiry int64, dc chan string) {
	for {
		select {
		case msg, ok := <-dc:
			if !ok {
				log.Printf("INFO -Data channel for %v received close signal", key)
				return
			}
			if time.Now().Unix() > expiry {
				log.Printf("INFO - Closing cache connection for expired server %s", key)
				close(dc)
				return
			}
			conn := c.cp.Get()
			if _, err := conn.Do("RPUSH", key, msg); err != nil {
				log.Printf("ERROR - Writing message to key %s failed", key)
			}
			conn.Close()
		}
	}
}
