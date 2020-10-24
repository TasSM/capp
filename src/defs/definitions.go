package defs

type CacheClientService interface {
	KeyExists(key string) bool
	CreateCacheArrayRecord(key string, ttl int64) error
	Start(key string, expiry int64, dc chan string)
}

type Request struct {
	key string
	ttl int32
}

type Message struct {
	key     string
	message string
}
