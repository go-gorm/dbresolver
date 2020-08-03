package dbresolver

import "gorm.io/gorm"

type Policy interface {
	Resolve([]gorm.ConnPool) gorm.ConnPool
}
