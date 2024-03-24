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
		var p gorm.ConnPool
		switch i % 3 {
		case 0:
			p = p1
		case 1:
			p = p2
		case 2:
			p = p3
		}
		if p != RoundRobinPolicy().Resolve(pools) {
			t.Errorf("RandomPolicy should return p1")
		}
	}
}
