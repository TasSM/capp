package defs

type CacheClientService interface {
	KeyExists(key string) bool
	GetActiveConnections() int
	CreateCacheArrayRecord(key string, ttl int64) error
	Start(key string, expiry int64, dc chan string)
	ReadArrayRecord(key string) ([]string, error)
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
