package pokecache

import (
	"testing"
	"time"
)

func TestReadLoop(t *testing.T) {
	var ttl time.Duration = 1
	cache := NewCache(ttl * time.Second)
	cache.Add("item", []byte{})

	if len(cache.Data) != 1 {
		t.Errorf("Could not add a value in cache")
	} else {
		ticker := time.NewTicker((ttl + 1) * time.Second)
		<-ticker.C
		if len(cache.Data) != 0 {
			t.Errorf("Read loop did not clear data from cache.")
		}
	}
}
