package itinerary

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/yourusername/itinerary-backend/itinerary/cache"
)

// ============================================================================
// DAY 6 INTEGRATION TESTS
// ============================================================================

// TestCacheHitsReduceDatabaseQueries verifies that cache hits reduce database round trips
func TestCacheHitsReduceDatabaseQueries(t *testing.T) {
	cacheManager := cache.NewMemoryCache()
	defer cacheManager.Close()

	testKey := "user:123:profile"
	testData := map[string]interface{}{
		"id":    123,
		"name":  "Test User",
		"email": "test@example.com",
	}

	queryCount := 0

	// First access - cache miss, count as DB query
	queryCount++
	cacheManager.Set(testKey, testData, 5*time.Minute)

	// Subsequent accesses - should hit cache
	for i := 0; i < 10; i++ {
		if _, err := cacheManager.Get(testKey); err != nil {
			queryCount++
		}
	}

	// Expected: 1 DB query + 10 cache hits
	if queryCount != 1 {
		t.Errorf("Expected 1 DB query, got %d", queryCount)
	}
	t.Log("✓ Cache hits reduce database queries")
}

// TestCacheInvalidationOnUpdates verifies that cache is properly invalidated on data updates
func TestCacheInvalidationOnUpdates(t *testing.T) {
	cacheManager := cache.NewMemoryCache()
	defer cacheManager.Close()

	testKey := "itinerary:456"
	originalData := map[string]interface{}{
		"id":       456,
		"title":    "Original Title",
		"modified": false,
	}

	// Set initial cache
	cacheManager.Set(testKey, originalData, 5*time.Minute)

	// Verify cache contains original data
	if value, err := cacheManager.Get(testKey); err == nil {
		if m, ok := value.(map[string]interface{}); ok && m["title"] == "Original Title" {
			t.Log("✓ Initial cache data stored correctly")
		}
	}

	// Simulate update - invalidate cache
	cacheManager.Delete(testKey)

	// Verify cache is invalidated
	if _, err := cacheManager.Get(testKey); err == nil {
		t.Error("Expected cache miss after invalidation")
	} else {
		t.Log("✓ Cache invalidation works")
	}

	// Cache new data
	updatedData := map[string]interface{}{
		"id":       456,
		"title":    "Updated Title",
		"modified": true,
	}
	cacheManager.Set(testKey, updatedData, 5*time.Minute)

	// Verify updated data is in cache
	if value, err := cacheManager.Get(testKey); err == nil {
		if m, ok := value.(map[string]interface{}); ok && m["title"] == "Updated Title" {
			t.Log("✓ Updated cache data stored correctly")
		}
	}
}

// TestMultiUserSessionManagement verifies that multiple user sessions are properly managed
func TestMultiUserSessionManagement(t *testing.T) {
	cacheManager := cache.NewMemoryCache()
	defer cacheManager.Close()

	numUsers := 10
	sessionTimeout := time.Minute

	// Create sessions for multiple users
	for i := 1; i <= numUsers; i++ {
		sessionKey := fmt.Sprintf("session:%d", i)
		sessionData := map[string]interface{}{
			"user_id":  i,
			"username": fmt.Sprintf("user%d", i),
			"login_at": time.Now(),
		}
		cacheManager.Set(sessionKey, sessionData, sessionTimeout)
	}

	// Verify all sessions exist
	activeCount := 0
	for i := 1; i <= numUsers; i++ {
		if _, err := cacheManager.Get(fmt.Sprintf("session:%d", i)); err == nil {
			activeCount++
		}
	}

	if activeCount != numUsers {
		t.Errorf("Expected %d active sessions, got %d", numUsers, activeCount)
	} else {
		t.Logf("✓ All %d sessions created and active", numUsers)
	}

	// Simulate session expiration for one user
	cacheManager.Delete("session:5")

	// Verify other sessions still exist
	remainingCount := 0
	for i := 1; i <= numUsers; i++ {
		if _, err := cacheManager.Get(fmt.Sprintf("session:%d", i)); err == nil {
			remainingCount++
		}
	}

	if remainingCount != numUsers-1 {
		t.Errorf("Expected %d sessions after deletion, got %d", numUsers-1, remainingCount)
	} else {
		t.Log("✓ Session deletion works correctly")
	}
}

// TestCacheExpiration validates that cache entries expire correctly
func TestCacheExpiration(t *testing.T) {
	cacheManager := cache.NewMemoryCache()
	defer cacheManager.Close()

	key := "expiring:key"
	value := "temporary_value"

	// Set with very short TTL
	cacheManager.Set(key, value, 100*time.Millisecond)

	// Should exist immediately
	if _, err := cacheManager.Get(key); err != nil {
		t.Error("Expected key to exist immediately after set")
	} else {
		t.Log("✓ Key exists after creation")
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be expired
	if _, err := cacheManager.Get(key); err == nil {
		t.Error("Expected key to be expired after TTL")
	} else {
		t.Log("✓ Key expired after TTL")
	}
}

// TestConcurrentCacheAccess verifies cache handles concurrent access safely
func TestConcurrentCacheAccess(t *testing.T) {
	cacheManager := cache.NewMemoryCache()
	defer cacheManager.Close()

	numGoroutines := 20
	operationsPerGoroutine := 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	errorCount := 0
	var mu sync.Mutex

	for g := 0; g < numGoroutines; g++ {
		go func(goroutineID int) {
			defer wg.Done()

			for op := 0; op < operationsPerGoroutine; op++ {
				key := fmt.Sprintf("concurrent:key:%d:%d", goroutineID, op)
				value := fmt.Sprintf("value_%d_%d", goroutineID, op)

				// Set
				cacheManager.Set(key, value, time.Minute)

				// Get
				retrieved, err := cacheManager.Get(key)
				if err != nil {
					mu.Lock()
					errorCount++
					mu.Unlock()
					continue
				}

				if retrieved != value {
					mu.Lock()
					errorCount++
					mu.Unlock()
				}

				// Delete
				cacheManager.Delete(key)
			}
		}(g)
	}

	wg.Wait()

	if errorCount > 0 {
		t.Errorf("Expected 0 errors in concurrent access, got %d", errorCount)
	} else {
		t.Logf("✓ Concurrent access with %d goroutines completed without errors", numGoroutines)
	}
}

// BenchmarkCacheSetVsGet compares cache set vs get performance
func BenchmarkCacheSetVsGet(b *testing.B) {
	cacheManager := cache.NewMemoryCache()
	defer cacheManager.Close()

	key := "bench:key"
	value := "benchmark_value"

	b.Run("Set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cacheManager.Set(fmt.Sprintf("%s:%d", key, i), value, time.Hour)
		}
	})

	// Pre-populate for Get benchmark
	cacheManager.Set(key, value, time.Hour)

	b.Run("Get", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cacheManager.Get(key)
		}
	})
}

// TestCacheExistsCheck verifies the Exists method works
func TestCacheExistsCheck(t *testing.T) {
	cacheManager := cache.NewMemoryCache()
	defer cacheManager.Close()

	key := "test:exists"

	// Should not exist initially
	if cacheManager.Exists(key) {
		t.Error("Key should not exist initially")
	}

	// Set value
	cacheManager.Set(key, "test_value", time.Minute)

	// Should exist now
	if !cacheManager.Exists(key) {
		t.Error("Key should exist after being set")
	} else {
		t.Log("✓ Exists check works correctly")
	}

	// Delete and check again
	cacheManager.Delete(key)
	if cacheManager.Exists(key) {
		t.Error("Key should not exist after being deleted")
	}
}

// TestCacheClear verifies the Clear method
func TestCacheClear(t *testing.T) {
	cacheManager := cache.NewMemoryCache()
	defer cacheManager.Close()

	// Add multiple items
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key:%d", i)
		cacheManager.Set(key, i, time.Minute)
	}

	// Verify items exist
	for i := 0; i < 10; i++ {
		if !cacheManager.Exists(fmt.Sprintf("key:%d", i)) {
			t.Errorf("Key %d should exist", i)
		}
	}

	// Clear all
	cacheManager.Clear()

	// Verify all cleared
	for i := 0; i < 10; i++ {
		if cacheManager.Exists(fmt.Sprintf("key:%d", i)) {
			t.Errorf("Key %d should not exist after clear", i)
		}
	}
	t.Log("✓ Cache clear works")
}
