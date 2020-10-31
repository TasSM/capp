package service

import (
	"testing"
	"time"

	"github.com/TasSM/appCache/defs"
	"github.com/alicebob/miniredis"
)

func SetupTest() (*miniredis.Miniredis, defs.CacheClientService) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	cc := NewCacheClient(s.Addr())
	return s, cc
}

func TeardownTest(s *miniredis.Miniredis, cc defs.CacheClientService) {
	s.Close()
	cc.DisposePool()
}

func TestRedisConnection(t *testing.T) {
	s, cc := SetupTest()
	defer TeardownTest(s, cc)
	err := cc.Ping()
	if err != nil {
		t.Errorf("Connection to redis from test client failed")
	}
}

func TestGetStatistics(t *testing.T) {
	s, cc := SetupTest()
	defer TeardownTest(s, cc)
	stats, err := cc.GetStatistics()
	if err != nil {
		t.Errorf("Unexpected error retrieving statistics")
	}
	if stats.ActiveConnections != 1 || stats.RecordCount != 0 {
		t.Errorf("statistic values are incorrect")
	}
}

func TestKeyCheckAbsent(t *testing.T) {
	s, cc := SetupTest()
	defer TeardownTest(s, cc)
	if res := cc.KeyExists("testkeyalpha"); res == true {
		t.Errorf("Key should not exist")
	}
}

func TestCreateRecord(t *testing.T) {
	key := "testkeyalpha"
	s, cc := SetupTest()
	defer TeardownTest(s, cc)
	if err := cc.CreateCacheArrayRecord(key, 100); err != nil {
		t.Errorf("Error creating new record")
	}
	ttl, err := cc.GetTTL(key)
	if err != nil {
		t.Errorf("Error checking ttl of record")
	}
	if ttl < 98 {
		t.Errorf("TTL outside of expected range")
	}
}

func TestCreateDuplicate(t *testing.T) {
	key := "testkeyalpha"
	s, cc := SetupTest()
	defer TeardownTest(s, cc)
	if err := cc.CreateCacheArrayRecord(key, 100); err != nil {
		t.Errorf("Error creating new record")
	}
	ttl, err := cc.GetTTL(key)
	if err != nil {
		t.Errorf("Error checking ttl of record")
	}
	if ttl < 98 {
		t.Errorf("TTL outside of expected range")
	}
	if err := cc.CreateCacheArrayRecord(key, 100); err == nil {
		t.Errorf("Record should not be duplicated")
	}
}

func TestKeyCheckPresent(t *testing.T) {
	key := "testkeyalpha"
	s, cc := SetupTest()
	defer TeardownTest(s, cc)
	if err := cc.CreateCacheArrayRecord(key, 100); err != nil {
		t.Errorf("Error creating new record")
	}
	if exists := cc.KeyExists(key); exists != true {
		t.Errorf("Key should be recorded as present")
	}
}

func TestReadWriteRecord(t *testing.T) {
	key := "testkeyalpha"
	s, cc := SetupTest()
	defer TeardownTest(s, cc)
	if err := cc.CreateCacheArrayRecord(key, 100); err != nil {
		t.Errorf("Error creating new record")
	}
	// Start the cache write values to the channel
	dc := make(chan string, 10)
	go cc.Start(key, time.Now().Unix()+500, dc)
	dc <- "abc"
	dc <- "def"
	dc <- "hij"
	close(dc)

}
