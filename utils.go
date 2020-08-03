package dbresolver

import (
	"regexp"

	"gorm.io/gorm/clause"
)

type Operation string

func (op Operation) Name() string {
	if op == Write {
		return "gorm:db_resolver:write"
	} else {
		return "gorm:db_resolver:read"
	}
}

func (op Operation) Build(clause.Builder) {
}

func (op Operation) MergeClause(*clause.Clause) {
}

var fromTableRegexp = regexp.MustCompile("(?i)FROM ['`\"]?([a-zA-Z0-9_]+)[ '`\"$,)]")

func getTableFromRawSQL(sql string) string {
	if matches := fromTableRegexp.FindAllStringSubmatch(sql, -1); len(matches) > 0 {
		return matches[0][1]
	}

	return ""
}
