package dbresolver

import (
	"errors"

	"gorm.io/gorm"
)

const (
	Write Operation = "write"
	Read  Operation = "read"
)

type DBResolver struct {
	*gorm.DB
	configs   []Config
	resolvers map[string]*resolver
	global    *resolver
}

type Config struct {
	Sources  []gorm.Dialector
	Replicas []gorm.Dialector
	Policy   Policy
	datas    []interface{}
}

func Register(config Config, datas ...interface{}) *DBResolver {
	return (&DBResolver{}).Register(config, datas...)
}

func (dr *DBResolver) Register(config Config, datas ...interface{}) *DBResolver {
	if dr.resolvers == nil {
		dr.resolvers = map[string]*resolver{}
	}

	config.datas = datas
	dr.configs = append(dr.configs, config)
	if dr.DB != nil {
		dr.compileConfig(config)
	}
	return dr
}

func (dr *DBResolver) Name() string {
	return "gorm:db_resolver"
}

func (dr *DBResolver) Initialize(db *gorm.DB) error {
	dr.DB = db
	return dr.compile()
}

func (dr *DBResolver) compile() error {
	for _, config := range dr.configs {
		if err := dr.compileConfig(config); err != nil {
			return err
		}
	}

	return nil
}

func (dr *DBResolver) convertToConnPool(dialectors []gorm.Dialector) (connPools []gorm.ConnPool, err error) {
	config := *dr.DB.Config
	for _, dialector := range dialectors {
		if db, err := gorm.Open(dialector, &config); err == nil {
			connPools = append(connPools, db.Config.ConnPool)
		} else {
			return nil, err
		}
	}

	return connPools, err
}

func (dr *DBResolver) compileConfig(config Config) (err error) {
	connPool := dr.DB.Config.ConnPool
	r := resolver{}

	if len(config.Sources) == 0 {
		r.sources = []gorm.ConnPool{connPool}
	} else if r.sources, err = dr.convertToConnPool(config.Sources); err != nil {
		return err
	}

	if len(config.Replicas) == 0 {
		r.replicas = r.sources
	} else if r.replicas, err = dr.convertToConnPool(config.Replicas); err != nil {
		return err
	}

	if len(config.datas) > 0 {
		for _, data := range config.datas {
			if t, ok := data.(string); ok {
				dr.resolvers[t] = &r
			} else {
				stmt := &gorm.Statement{DB: dr.DB}
				if err := stmt.Parse(data); err == nil {
					dr.resolvers[stmt.Table] = &r
				} else {
					return err
				}
			}
		}
	} else if dr.global != nil {
		dr.global = &r
	} else {
		return errors.New("conflicted global resolver")
	}

	return nil
}

func (dr *DBResolver) resolve(stmt *gorm.Statement, op Operation) gorm.ConnPool {
	if len(dr.resolvers) > 0 {
		if stmt.Table != "" {
			if r, ok := dr.resolvers[stmt.Table]; ok {
				return r.resolve(op)
			}
		}

		if stmt.Schema != nil {
			if r, ok := dr.resolvers[stmt.Schema.Table]; ok {
				return r.resolve(op)
			}
		}

		if rawSQL := stmt.SQL.String(); rawSQL != "" {
			if r, ok := dr.resolvers[getTableFromRawSQL(rawSQL)]; ok {
				return r.resolve(op)
			}
		}
	}

	return dr.global.resolve(op)
}

type resolver struct {
	sources  []gorm.ConnPool
	replicas []gorm.ConnPool
	policy   Policy
}

func (r *resolver) resolve(op Operation) gorm.ConnPool {
	if op == Read {
		if len(r.replicas) == 1 {
			return r.replicas[0]
		}
		return r.policy.Resolve(r.replicas)
	}

	if len(r.sources) == 1 {
		return r.sources[0]
	}
	return r.policy.Resolve(r.sources)
}
