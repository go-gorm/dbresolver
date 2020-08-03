# GORM DBResolver

```go
db, err := gorm.Open(mysql.Open("dsn"), gorm.Config{})

dbresolver := dbresolver.Register(dbresolver.Config{
  Sources:  []gorm.Dialector{}, // or use default
  Replicas: []gorm.Dialector{},
}).Register(dbresolver.Config{
  Sources:  gorm.Dialector, // or use default
  Replicas: []gorm.Dialector{},
}, &User{}, &Product{}).Register(dbresolver.Config{
  Source:   gorm.Dialector, // or use default
  Replicas: []gorm.Dialector{},
}, "users", &Order{})
```

dbresolver
