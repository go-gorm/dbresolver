package dbresolver

import (
	"math/rand"

	"gorm.io/gorm"
)

type Policy interface {
	Resolve([]gorm.ConnPool) gorm.ConnPool
}

type PolicyFunc func([]gorm.ConnPool) gorm.ConnPool

func (f PolicyFunc) Resolve(connPools []gorm.ConnPool) gorm.ConnPool {
	return f(connPools)
}

type RandomPolicy struct {
}

func (RandomPolicy) Resolve(connPools []gorm.ConnPool) gorm.ConnPool {
	return connPools[rand.Intn(len(connPools))]
}

func RoundRobinPolicy() Policy {
	var i int
	return PolicyFunc(func(connPools []gorm.ConnPool) gorm.ConnPool {
		i = (i + 1) % len(connPools)
		return connPools[i]
	})
}
