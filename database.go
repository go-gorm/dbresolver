package dbresolver

import (
	"time"

	"gorm.io/gorm"
)

func (dr *DBResolver) SetConnMaxIdleTime(d time.Duration) *DBResolver {
	dr.Call(func(connPool gorm.ConnPool) error {
		connPool.(interface{ SetConnMaxIdleTime(time.Duration) }).SetConnMaxIdleTime(d)
		return nil
	})

	return dr
}

func (dr *DBResolver) SetConnMaxLifetime(d time.Duration) *DBResolver {
	dr.Call(func(connPool gorm.ConnPool) error {
		connPool.(interface{ SetConnMaxLifetime(time.Duration) }).SetConnMaxLifetime(d)
		return nil
	})

	return dr
}

func (dr *DBResolver) SetMaxIdleConns(n int) *DBResolver {
	dr.Call(func(connPool gorm.ConnPool) error {
		connPool.(interface{ SetMaxIdleConns(int) }).SetMaxIdleConns(n)
		return nil
	})

	return dr
}

func (dr *DBResolver) SetMaxOpenConns(n int) *DBResolver {
	dr.Call(func(connPool gorm.ConnPool) error {
		connPool.(interface{ SetMaxOpenConns(int) }).SetMaxOpenConns(n)
		return nil
	})

	return dr
}

func (dr *DBResolver) Call(fc func(connPool gorm.ConnPool) error) error {
	if dr.DB != nil {
		for _, r := range dr.resolvers {
			if err := r.call(fc); err != nil {
				return err
			}
		}
	} else {
		dr.compileCallbacks = append(dr.compileCallbacks, fc)
	}

	return nil
}
