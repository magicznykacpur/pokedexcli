package pokecache

import (
	"testing"
	"time"
)
func TestNewCache(t *testing.T) {
	cases := []struct{
		input time.Duration
		expected Cache
	}{
		{
			input: time.Second * 5,
			expected: Cache{entries: map[string]cacheEntry{}, interval: time.Second * 5},
		},
		{
			input: time.Millisecond * 50,
			expected: Cache{entries: map[string]cacheEntry{}, interval: time.Millisecond * 50},
		},
		{
			input: time.Minute * 5,
			expected: Cache{entries: map[string]cacheEntry{}, interval: time.Minute * 5},
		},
	}

	for _, c := range cases {
		cache := NewCache(c.input)
		if cache.interval != c.input {
			t.Errorf("new cache interval different than expected: %v != %v", cache.interval, c.input)
			return
		}
	}
}

func TestReapLoop(t *testing.T) {
	interval := time.Millisecond * 50
	cache := NewCache(interval)

	cache.Add("item", []byte("test data"))

	_, ok := cache.Get("item")
	if !ok {
		t.Errorf("expected to find the key")
		return
	}

	time.Sleep(interval + time.Millisecond * 20)

	_, ok = cache.Get("item")
	if ok {
		t.Errorf("expected to not find the key")
		return
	}
}

func TestAddGet(t *testing.T) {
	cache := NewCache(time.Millisecond * 50)
	cache.Add("item", []byte("new item"))

	_, ok := cache.Get("item")
	if !ok {
		t.Errorf("expected to get item")
	}
}