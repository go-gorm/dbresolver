package dbresolver

import "testing"

func TestGetTableFromRawSQL(t *testing.T) {
	datas := [][2]string{
		{"select * from users as u", "users"},
		{"select * from (select * from users) as u", "users"},
		{"select * from (select * from users)", "users"},
		{"select * from (select * from users), (select * from products)", "users"},
		{"select * from users, products", "users"},
		{"select * from users as u, products as p", "users"},
	}

	for _, data := range datas {
		if getTableFromRawSQL(data[0]) != data[1] {
			t.Errorf("Failed to get table name from %v, expect: %v, got: %v", data[0], data[1], getTableFromRawSQL(data[0]))
		}
	}
}
