package common

import (
	"sync"
	"testing"
)

func TestHealthStatus_ConcurrentAccess(t *testing.T) {
	status := NewHealthStatus()

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() {
			status.Set(i%2 == 0)
		})
		wg.Go(func() {
			_, _ = status.Get()
		})
	}
	wg.Wait()
}
