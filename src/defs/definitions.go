package defs

import (
	"net/http"
)

type CacheStartFunction func(string, int64, chan string)

type CacheClientService interface {
	KeyExists(key string) bool
	CreateCacheArrayRecord(key string, ttl int64) error
	DisposePool()
	GetStatistics() (*StatisticResponse, error)
	GetTTL(key string) (int, error)
	Ping() error
	Start(key string, expiry int64, dc chan string)
	ReadArrayRecord(key string) ([]string, error)
}

type CacheClientRouter interface {
	HandleHealthcheck(w http.ResponseWriter, r *http.Request)
	HandleStatistics(w http.ResponseWriter, r *http.Request)
}

type StatisticResponse struct {
	RecordCount       int
	ActiveConnections int
	Timestamp         string
}

type TimedChannel struct {
	DataChannel chan string
	Expiry      int64
}

type Request struct {
	key string
	ttl int32
}

type Message struct {
	key     string
	message string
}
