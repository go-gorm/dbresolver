package dbresolver

import (
	"database/sql"

	"gorm.io/gorm"
)

// connPool wraps original DB connection pool with database resolver.
type connPool struct {
	gorm.ConnPool
	dr *DBResolver
}

// newConnPool creates a new ConnPool.
func newConnPool(base gorm.ConnPool, dr *DBResolver) connPool {
	return connPool{ConnPool: base, dr: dr}
}

// GetDBConnWithContext gets *sql.DB connection based on the context. If no
// information is available, returns original connection.
func (p connPool) GetDBConnWithContext(db *gorm.DB) (*sql.DB, error) {
	var w *gorm.DB

	if stmt := db.Statement; stmt != nil {
		if r := p.dr.getResolver(stmt); r != nil {
			if _, ok := db.Statement.Settings.Load(readName); ok {
				w = p.wrap(r.resolve(stmt, Read))
			} else if _, ok := db.Statement.Settings.Load(writeName); ok {
				w = p.wrap(r.resolve(stmt, Write))
			}
		}
	}

	if w == nil {
		w = p.wrap(p.dr.original)
	}

	return w.DB()
}

func (connPool) wrap(cp gorm.ConnPool) *gorm.DB {
	return &gorm.DB{Config: &gorm.Config{ConnPool: cp}}
}
