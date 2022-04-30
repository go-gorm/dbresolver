package dbresolver

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Operation specifies dbresolver mode
type Operation string

const writeName = "gorm:db_resolver:write"

// ModifyStatement modify operation mode
func (op Operation) ModifyStatement(stmt *gorm.Statement) {
	optName := "gorm:db_resolver:read"
	if op == Write {
		optName = writeName
	}

	stmt.Clauses[optName] = clause.Clause{}
}

// Build implements clause.Expression interface
func (op Operation) Build(clause.Builder) {
}

// Use specifies configuration
func Use(str string) clause.Expression {
	return using{Use: str}
}

type using struct {
	Use string
}

const usingName = "gorm:db_resolver:using"

// ModifyStatement modify operation mode
func (u using) ModifyStatement(stmt *gorm.Statement) {
	stmt.Clauses[usingName] = clause.Clause{Expression: u}
}

// Build implements clause.Expression interface
func (u using) Build(clause.Builder) {
}
