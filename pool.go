package dbresolver

import (
	"context"
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

func (p connPool) BeginTx(ctx context.Context, opt *sql.TxOptions) (conn gorm.ConnPool, err error) {
	switch beginner := p.ConnPool.(type) {
	case gorm.TxBeginner:
		conn, err = beginner.BeginTx(ctx, opt)
	case gorm.ConnPoolBeginner:
		conn, err = beginner.BeginTx(ctx, opt)
	default:
		err = gorm.ErrInvalidTransaction
	}
	return
}

// GetDBConnWithContext gets *sql.DB connection based on the context. If no
// information is available, returns original connection.
func (p connPool) GetDBConnWithContext(db *gorm.DB) (*sql.DB, error) {
	var gormPool gorm.ConnPool = p
	if stmt := db.Statement; stmt != nil {
		if r := p.dr.getResolver(stmt); r != nil {
			if _, ok := db.Statement.Settings.Load(readName); ok {
				gormPool = r.resolve(stmt, Read)
			} else if _, ok := db.Statement.Settings.Load(writeName); ok {
				gormPool = r.resolve(stmt, Write)
			} else {
				gormPool = p.dr.original
			}
		}
	}

	if pool, ok := gormPool.(connPool); ok {
		gormPool = pool.ConnPool
	}

	if sqlDB, ok := gormPool.(*sql.DB); ok {
		return sqlDB, nil
	}

	if pool, ok := gormPool.(gorm.GetDBConnectorWithContext); ok {
		return pool.GetDBConnWithContext(db)
	}

	if connPool, ok := gormPool.(gorm.GetDBConnector); ok {
		return connPool.GetDBConn()
	}

	return nil, gorm.ErrInvalidDB
}
