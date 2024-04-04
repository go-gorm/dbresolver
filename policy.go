package dbresolver

import (
	"math/rand"
	"sync"

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

type RoundRobinPolicy struct {
	lock sync.Mutex
	i    int
}

func (p *RoundRobinPolicy) Resolve(connPools []gorm.ConnPool) gorm.ConnPool {
	p.lock.Lock()
	defer func() {
		p.i = (p.i + 1) % len(connPools)
		p.lock.Unlock()
	}()
	return connPools[p.i]
}
