# go-pools
Provides functionality to manage and reuse resources like connections.

Base on [vitess](!https://github.com/vitessio/vitess)

## Fast Start

```go
func main() {
    factory := func(context.Context) (pools.Resource, error) {
        return sql.Open("mysql", "user:password@/dbname")
    }
    
    rp := pools.NewResourcePool(factory, 10, 30, 1*time.Hour, 1, nil)
    defer rp.Close()
    
    resource, err := rp.Get(context.Background())
    if err != nil {
        panic(err)
    }
    db := resource.(*sql.DB)
    
    // Execute the query
    _, err = db.Query("SELECT * FROM table Limit 10")
    if err != nil {
        panic(err.Error())
    }
    
    rp.Put(db)
}

```

# Example

- [mysql](./example/mysql.go)
- [gorm](./example/gorm.go)
- [amqp](./example/amqp.go)
- [hive](./example/hive.go)