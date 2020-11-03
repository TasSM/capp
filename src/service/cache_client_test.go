package service

import (
	"sync"
	"testing"
	"time"

	"github.com/TasSM/appCache/defs"
	"github.com/alicebob/miniredis"
)

func setupTest() (*miniredis.Miniredis, defs.CacheClientService) {
	s, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	cc := NewCacheClient(s.Addr())
	return s, cc
}

func teardownTest(s *miniredis.Miniredis, cc defs.CacheClientService) {
	s.Close()
	cc.DisposePool()
}

func TestRedisConnection(t *testing.T) {
	s, cc := setupTest()
	defer teardownTest(s, cc)
	err := cc.Ping()
	if err != nil {
		t.Errorf("Connection to redis from test client failed")
	}
}

func TestGetStatistics(t *testing.T) {
	s, cc := setupTest()
	defer teardownTest(s, cc)
	stats, err := cc.GetStatistics()
	if err != nil {
		t.Errorf("Unexpected error retrieving statistics")
	}
	if stats.ActiveConnections != 1 || stats.RecordCount != 0 {
		t.Errorf("statistic values are incorrect")
	}
}

func TestKeyCheckAbsent(t *testing.T) {
	s, cc := setupTest()
	defer teardownTest(s, cc)
	if res := cc.KeyExists("testkeyalpha"); res == true {
		t.Errorf("Key should not exist")
	}
}

func TestCreateRecord(t *testing.T) {
	key := "testkeyalpha"
	s, cc := setupTest()
	defer teardownTest(s, cc)
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
	s, cc := setupTest()
	defer teardownTest(s, cc)
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
	s, cc := setupTest()
	defer teardownTest(s, cc)
	if err := cc.CreateCacheArrayRecord(key, 100); err != nil {
		t.Errorf("Error creating new record")
	}
	if exists := cc.KeyExists(key); exists != true {
		t.Errorf("Key should be recorded as present")
	}
}

func TestReadWriteRecord(t *testing.T) {
	key := "testkeyalpha"
	model := []string{"abc", "def", "hij"}
	s, cc := setupTest()
	defer teardownTest(s, cc)
	if err := cc.CreateCacheArrayRecord(key, 100); err != nil {
		t.Errorf("Error creating new record")
	}
	// Start the cache write values to the channel
	dc := make(chan string, 10)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cc.Start(key, time.Now().Unix()+500, dc)
	}()
	for i := range model {
		dc <- model[i]
	}
	close(dc)
	wg.Wait()
	res, err := cc.ReadArrayRecord(key)
	if err != nil {
		t.Errorf("Error reading record for key %v", key)
	}
	for i, v := range res {
		if v != model[i] {
			t.Errorf("Mismatched element in read output")
		}
	}
}
