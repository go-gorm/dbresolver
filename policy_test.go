package dbresolver

import (
	"sync"
	"testing"

	"gorm.io/gorm"
)

type P1ConnPool struct {
	gorm.ConnPool
}

type P2ConnPool struct {
	gorm.ConnPool
}

type P3ConnPool struct {
	gorm.ConnPool
}

func TestPolicy_RoundRobinPolicy(t *testing.T) {
	var p1 P1ConnPool
	var p2 P2ConnPool
	var p3 P3ConnPool
	var pools = []gorm.ConnPool{
		p1, p2, p3,
	}

	t.Run("single", func(t *testing.T) {
		policy := RoundRobinPolicy{}
		for i := 0; i < 10; i++ {
			pool1 := pools[i%3]
			pool2 := policy.Resolve(pools)
			if pool1 != pool2 {
				t.Errorf("RoundRobinPolicy failed")
			}
		}
	})

	// checks for race condition, enable with -race
	t.Run("concurrent", func(t *testing.T) {
		policy := RoundRobinPolicy{}
		var wg sync.WaitGroup
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func(i int) {
				defer wg.Done()
				policy.Resolve(pools)
			}(i)
		}
		wg.Wait()
	})
}
