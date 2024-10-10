package internal

import (
	"sync"
	"testing"
	"time"
)

func TestIDStoreLoad(t *testing.T) {
	store, err := NewIDStore()
	if err != nil {
		t.Fatalf("failed to create IDStore: %v", err)
	}
	defer store.Close()

	numWorkers := 100
	numRequestsPerWorker := 1000

	var wg sync.WaitGroup
	var mu sync.Mutex

	errorsOccurred := 0

	startTime := time.Now()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for j := 0; j < numRequestsPerWorker; j++ {
				// Generate a new ID
				id := store.GetId()
				if id == "" {
					mu.Lock()
					errorsOccurred++
					mu.Unlock()
					t.Errorf("worker %d failed to generate ID: %v", workerID, err)
				}
			}
		}(i)
	}

	wg.Wait()

	duration := time.Since(startTime)
	t.Logf("Completed %d requests with %d workers in %v", numRequestsPerWorker*numWorkers, numWorkers, duration)

	// Report any errors that occurred during the test
	if errorsOccurred > 0 {
		t.Errorf("Encountered %d errors during load test", errorsOccurred)
	} else {
		t.Log("No errors encountered during load test")
	}
}

func BenchmarkIDStoreLoad(b *testing.B) {
	store, err := NewIDStore()
	if err != nil {
		b.Fatalf("failed to create IDStore: %v", err)
	}
	defer store.Close()

	for i := 0; i < b.N; i++ {
		id := store.GetId()
		if id == "" {
			b.Errorf("failed to generate ID: %v", err)
		}
	}
}
