package dbresolver

import (
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

	}
}
