package defs

type CacheClientService interface {
	KeyExists(key string) bool
	CreateCacheArrayRecord(key string, ttl int32) error
	Start(key string, expiry int32, dc chan string)
}

type Request struct {
	key string
	ttl int32
}

type Message struct {
	key     string
	message string
}
