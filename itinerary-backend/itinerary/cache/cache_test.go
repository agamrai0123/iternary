package cache

import (
	"testing"
	"time"
)

// TestMemoryCacheSet tests the Set operation
func TestMemoryCacheSet(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	err := cache.Set("key1", "value1", time.Minute)
	if err != nil {
		t.Errorf("Set failed: %v", err)
	}

	// Verify key exists
	exists := cache.Exists("key1")
	if !exists {
		t.Error("Key should exist after Set")
	}
}

// TestMemoryCacheGet tests the Get operation
func TestMemoryCacheGet(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	cache.Set("key1", "value1", time.Minute)
	value, err := cache.Get("key1")

	if err != nil {
		t.Errorf("Get failed: %v", err)
	}

	if value != "value1" {
		t.Errorf("Expected 'value1', got %v", value)
	}
}

// TestMemoryCacheGetMiss tests Get on non-existent key
func TestMemoryCacheGetMiss(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	_, err := cache.Get("nonexistent")
	if err == nil {
		t.Error("Get should fail for non-existent key")
	}
}

// TestMemoryCacheDelete tests the Delete operation
func TestMemoryCacheDelete(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	cache.Set("key1", "value1", time.Minute)
	cache.Delete("key1")

	exists := cache.Exists("key1")
	if exists {
		t.Error("Key should not exist after Delete")
	}
}

// TestMemoryCacheExpiration tests key expiration
func TestMemoryCacheExpiration(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	cache.Set("key1", "value1", 100*time.Millisecond)

	// Should exist immediately
	if !cache.Exists("key1") {
		t.Error("Key should exist immediately after Set")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should not exist after expiration
	_, err := cache.Get("key1")
	if err == nil {
		t.Error("Key should not exist after expiration")
	}
}

// TestMemoryCacheSetIfNotExists tests SetIfNotExists
func TestMemoryCacheSetIfNotExists(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	// First call should succeed
	set, err := cache.SetIfNotExists("key1", "value1", time.Minute)
	if !set || err != nil {
		t.Error("First SetIfNotExists should succeed")
	}

	// Second call should fail
	set, err = cache.SetIfNotExists("key1", "value2", time.Minute)
	if set || err != nil {
		t.Error("Second SetIfNotExists should fail")
	}

	// Verify value is unchanged
	value, _ := cache.Get("key1")
	if value != "value1" {
		t.Errorf("Value should be 'value1', got %v", value)
	}
}

// TestMemoryCacheIncrement tests Increment operation
func TestMemoryCacheIncrement(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	// Set initial value
	cache.Set("counter", 10, time.Minute)

	// Increment
	cache.Increment("counter")

	value, _ := cache.Get("counter")
	if value != 11 {
		t.Errorf("Expected 11, got %v", value)
	}
}

// TestMemoryCacheDecrement tests Decrement operation
func TestMemoryCacheDecrement(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	// Set initial value
	cache.Set("counter", 10, time.Minute)

	// Decrement
	cache.Decrement("counter")

	value, _ := cache.Get("counter")
	if value != 9 {
		t.Errorf("Expected 9, got %v", value)
	}
}

// TestMemoryCacheClear tests Clear operation
func TestMemoryCacheClear(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	// Set multiple values
	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)
	cache.Set("key3", "value3", time.Minute)

	// Clear all
	cache.Clear()

	if cache.Size() != 0 {
		t.Error("Cache should be empty after Clear")
	}
}

// TestMemoryCacheSize tests Size operation
func TestMemoryCacheSize(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	for i := 1; i <= 5; i++ {
		cache.Set("key"+string(rune('0'+i)), "value", time.Minute)
		if cache.Size() != i {
			t.Errorf("Expected size %d, got %d", i, cache.Size())
		}
	}
}

// TestMemoryCacheGetAll tests GetAll operation
func TestMemoryCacheGetAll(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)
	cache.Set("key3", "value3", time.Minute)

	all := cache.GetAll()
	if len(all) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(all))
	}
}

// TestMemoryCacheGetKeys tests GetKeys operation
func TestMemoryCacheGetKeys(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)
	cache.Set("key3", "value3", time.Minute)

	keys := cache.GetKeys()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}
}

// TestMemoryCacheConcurrency tests concurrent operations
func TestMemoryCacheConcurrency(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	done := make(chan bool)

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func(id int) {
			for j := 0; j < 100; j++ {
				cache.Set("key"+string(rune('0'+(id*10+j)%10)), id, time.Minute)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	if cache.Size() <= 0 {
		t.Error("Cache should have entries after concurrent writes")
	}
}

// TestCacheFactoryMemory tests Factory with Memory type
func TestCacheFactoryMemory(t *testing.T) {
	config := &Config{
		Type: Memory,
	}

	factory := NewFactory(config)
	cache, err := factory.Create()

	if err != nil {
		t.Errorf("Create failed: %v", err)
	}

	if cache == nil {
		t.Error("Cache instance should not be nil")
	}

	cache.Close()
}

// TestCacheBuilder tests CacheBuilder
func TestCacheBuilder(t *testing.T) {
	cache, err := NewCacheBuilder().
		WithType(Memory).
		WithDefaultTTL(time.Hour).
		WithMaxSize(10000).
		Build()

	if err != nil {
		t.Errorf("Build failed: %v", err)
	}

	if cache == nil {
		t.Error("Cache instance should not be nil")
	}

	cache.Close()
}

// TestCacheManager tests CacheManager
func TestCacheManager(t *testing.T) {
	manager := NewCacheManager()
	defer manager.Close()

	cache1 := NewMemoryCache()
	cache2 := NewMemoryCache()

	// Register caches
	err := manager.Register("cache1", cache1)
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}

	err = manager.Register("cache2", cache2)
	if err != nil {
		t.Errorf("Register failed: %v", err)
	}

	// Get cache
	retrieved, err := manager.Get("cache1")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}

	if retrieved != cache1 {
		t.Error("Retrieved cache should be cache1")
	}

	// Test duplicate registration
	err = manager.Register("cache1", cache1)
	if err == nil {
		t.Error("Duplicate registration should fail")
	}

	// Get all
	all := manager.GetAll()
	if len(all) != 2 {
		t.Errorf("Expected 2 caches, got %d", len(all))
	}
}

// TestMemoryCacheTypes tests different value types
func TestMemoryCacheTypes(t *testing.T) {
	cache := NewMemoryCache()
	defer cache.Close()

	// String
	cache.Set("string", "value", time.Minute)

	// Number
	cache.Set("number", 42, time.Minute)

	// Boolean
	cache.Set("bool", true, time.Minute)

	// Slice
	cache.Set("slice", []int{1, 2, 3}, time.Minute)

	// Map
	cache.Set("map", map[string]string{"key": "value"}, time.Minute)

	// Verify all values
	if val, _ := cache.Get("string"); val != "value" {
		t.Error("String value mismatch")
	}

	if val, _ := cache.Get("number"); val != 42 {
		t.Error("Number value mismatch")
	}

	if val, _ := cache.Get("bool"); val != true {
		t.Error("Boolean value mismatch")
	}
}

// BenchmarkMemoryCacheSet benchmarks Set operation
func BenchmarkMemoryCacheSet(b *testing.B) {
	cache := NewMemoryCache()
	defer cache.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Set("key", "value", time.Minute)
	}
}

// BenchmarkMemoryCacheGet benchmarks Get operation
func BenchmarkMemoryCacheGet(b *testing.B) {
	cache := NewMemoryCache()
	defer cache.Close()

	cache.Set("key", "value", time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get("key")
	}
}

// BenchmarkMemoryCacheExists benchmarks Exists operation
func BenchmarkMemoryCacheExists(b *testing.B) {
	cache := NewMemoryCache()
	defer cache.Close()

	cache.Set("key", "value", time.Minute)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Exists("key")
	}
}
