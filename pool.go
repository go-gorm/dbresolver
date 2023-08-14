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
	var connPool gorm.ConnPool = p
	if stmt := db.Statement; stmt != nil {
		if r := p.dr.getResolver(stmt); r != nil {
			if _, ok := db.Statement.Settings.Load(readName); ok {
				connPool = r.resolve(stmt, Read)
			} else if _, ok := db.Statement.Settings.Load(writeName); ok {
				connPool = r.resolve(stmt, Write)
			} else {
				connPool = p.dr.original
			}
		}
	}

	if sqlDB, ok := connPool.(*sql.DB); ok {
		return sqlDB, nil
	}

	if pool, ok := connPool.(gorm.GetDBConnectorWithContext); ok && connPool != p {
		return pool.GetDBConnWithContext(db)
	}

	if connPool, ok := connPool.(gorm.GetDBConnector); ok {
		return connPool.GetDBConn()
	}

	return nil, gorm.ErrInvalidDB
}
