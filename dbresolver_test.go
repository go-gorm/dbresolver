package dbresolver_test

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type User struct {
	ID   uint
	Name string
}

type Product struct {
	ID   uint
	Name string
}

type Order struct {
	ID      uint
	OrderNo string
}

func GetDB(port int) *gorm.DB {
	DB, err := gorm.Open(mysql.Open(fmt.Sprintf("gorm:gorm@tcp(localhost:%v)/gorm?charset=utf8&parseTime=True&loc=Local", port)), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect db, got error: %v, port: %v", err, port))
	}
	return DB
}

func init() {
	for _, port := range []int{9911, 9912, 9913, 9914} {
		DB := GetDB(port)
		DB.Migrator().DropTable(&User{}, &Product{}, &Order{})
		DB.AutoMigrate(&User{}, &Product{}, &Order{})

		DB.Create(&User{Name: fmt.Sprintf("%v", port)})
		DB.Create(&Product{Name: fmt.Sprintf("%v", port)})
		DB.Create(&Order{OrderNo: fmt.Sprintf("%v", port)})
	}
}

func TestDBResolver(t *testing.T) {
	for i := 0; i < 2; i++ {
		DB, err := gorm.Open(mysql.Open("gorm:gorm@tcp(localhost:9911)/gorm?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{PrepareStmt: i%2 == 0})
		if err != nil {
			t.Fatalf("failed to connect db, got error: %v", err)
		}

		if err := DB.Use(dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{mysql.Open("gorm:gorm@tcp(localhost:9911)/gorm?charset=utf8&parseTime=True&loc=Local")},
			Replicas: []gorm.Dialector{
				mysql.Open("gorm:gorm@tcp(localhost:9912)/gorm?charset=utf8&parseTime=True&loc=Local"),
				mysql.Open("gorm:gorm@tcp(localhost:9913)/gorm?charset=utf8&parseTime=True&loc=Local"),
			},
		}).Register(dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open("gorm:gorm@tcp(localhost:9914)/gorm?charset=utf8&parseTime=True&loc=Local")},
			Replicas: []gorm.Dialector{mysql.Open("gorm:gorm@tcp(localhost:9913)/gorm?charset=utf8&parseTime=True&loc=Local")},
		}, "users", &Product{}).SetMaxOpenConns(5)); err != nil {
			t.Fatalf("failed to use plugin, got error: %v", err)
		}

		for j := 0; j < 20; j++ {
			// test query
			var order Order
			DB.First(&order)
			if order.OrderNo != "9912" && order.OrderNo != "9913" {
				t.Fatalf("idx: %v: order should comes from read db, but got order %v", j, order.OrderNo)
			}

			DB.Clauses(dbresolver.Write).First(&order)
			if order.OrderNo != "9911" {
				t.Fatalf("idx: %v: order should comes from write db, but got order %v", j, order.OrderNo)
			}

			DB.Clauses(dbresolver.Use("users")).First(&order)
			if order.OrderNo != "9913" {
				t.Fatalf("idx: %v: order should comes from write db @ users, but got order %v", j, order.OrderNo)
			}

			DB.Clauses(dbresolver.Use("users"), dbresolver.Write).First(&order)
			if order.OrderNo != "9914" {
				t.Fatalf("idx: %v: order should comes from write db @ users, but got order %v", j, order.OrderNo)
			}

			var user User
			DB.First(&user)
			if user.Name != "9913" {
				t.Fatalf("idx: %v: user should comes from read db, but got %v", j, user.Name)
			}

			DB.Clauses(dbresolver.Write).First(&user)
			if user.Name != "9914" {
				t.Fatalf("idx: %v: user should comes from read db, but got %v", j, user.Name)
			}

			var product Product
			DB.First(&product)
			if product.Name != "9913" {
				t.Fatalf("idx: %v: product should comes from read db, but got %v", j, product.Name)
			}

			DB.Clauses(dbresolver.Write).First(&product)
			if product.Name != "9914" {
				t.Fatalf("idx: %v: product should comes from write db, but got %v", j, product.Name)
			}

			// test create
			DB.Create(&User{Name: "create"})
			if err := DB.First(&User{}, "name = ?", "create").Error; err == nil {
				t.Fatalf("can't read user from read db, got no error happened")
			}

			if err := DB.Clauses(dbresolver.Write).First(&User{}, "name = ?", "create").Error; err != nil {
				t.Fatalf("read user from write db, got error: %v", err)
			}

			DB9914 := GetDB(9914)

			if err := DB9914.First(&User{}, "name = ?", "create").Error; err != nil {
				t.Fatalf("read user from write db, got error: %v", err)
			}

			var name string
			if err := DB.Raw("select name from users").Row().Scan(&name); err != nil || name != "9913" {
				t.Fatalf("read users from read db, name %v", name)
			}

			if err := DB.Debug().Raw("select name from users where name = ? for update", "create").Row().Scan(&name); err != nil || name != "create" {
				t.Fatalf("read users from write db, name %v, err %v", name, err)
			}

			// test update
			if err := DB.Model(&User{}).Where("name = ?", "create").Update("name", "update").Error; err != nil {
				t.Fatalf("failed to update users, got error: %v", err)
			}

			if err := DB9914.First(&User{}, "name = ?", "update").Error; err != nil {
				t.Fatalf("read user from write db, got error: %v", err)
			}

			// test raw sql
			name = ""
			if err := DB.Raw("select name from users where name = ?", "update").Row().Scan(&name); err == nil || name != "" {
				t.Fatalf("can't read users from read db, name %v", name)
			}

			if err := DB.Clauses(dbresolver.Write).Raw("select name from users where name = ?", "update").Row().Scan(&name); err != nil || name != "update" {
				t.Fatalf("read users from write db, error %v, name %v", err, name)
			}

			// test delete
			if err := DB.Where("name = ?", "update").Delete(&User{}).Error; err != nil {
				t.Fatalf("failed to delete users, got error: %v", err)
			}

			if err := DB9914.First(&User{}, "name = ?", "update").Error; err != gorm.ErrRecordNotFound {
				t.Fatalf("read user from write db after delete, got error: %v", err)
			}
		}
	}
}
