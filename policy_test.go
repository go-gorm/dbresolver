package dbresolver

import (
	"sync/atomic"
	"testing"

	"gorm.io/gorm"
)

func TestPolicy_RoundRobinPolicy(t *testing.T) {
	var p1, p2, p3 gorm.ConnPool
	var pools = []gorm.ConnPool{
		p1, p2, p3,
	}

	for i := 0; i < 10; i++ {
		if pools[i%3] != RoundRobinPolicy().Resolve(pools) {
			t.Errorf("RoundRobinPolicy failed")
		}
		if pools[i%3] != StrictRoundRobinPolicy().Resolve(pools) {
			t.Errorf("StrictRoundRobinPolicy failed")
		}
	}
}

func BenchmarkPolicy_StrictRoundRobinPolicy(b *testing.B) {
	var p1, p2, p3 gorm.ConnPool
	var pools = []gorm.ConnPool{
		p1, p2, p3,
	}

	var i int64
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if pools[int(atomic.AddInt64(&i, 1))%3] != StrictRoundRobinPolicy().Resolve(pools) {
				b.Errorf("RoundRobinPolicy failed")
			}
		}
	})
}
