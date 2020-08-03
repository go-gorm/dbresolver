# GORM DBResolver

```go
db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

db.Use(dbresolver.Register(dbresolver.Config{
  Sources:  []gorm.Dialector{}, // or use default
  Replicas: []gorm.Dialector{},
}).Register(dbresolver.Config{
  Sources:  gorm.Dialector, // or use default
  Replicas: []gorm.Dialector{},
}, &User{}, &Product{}).Register(dbresolver.Config{
  Source:   gorm.Dialector, // or use default
  Replicas: []gorm.Dialector{},
}, "users", &Order{}))
```
