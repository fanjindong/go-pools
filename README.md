# go-pools
Provides functionality to manage and reuse resources like connections.

Base on [vitess](https://github.com/vitessio/vitess)

## Fast Start

An AMPQ connection channel pool demo.
```go
func main() {
    conn, err := amqp.Dial("amqp://root:password@127.0.0.1:5672/")
    if err != nil {
        panic("Failed to connect to RabbitMQ")
    }
    defer conn.Close()
    factory := func(context.Context) (pools.Resource, error) {
        return conn.Channel()
    }
    rp := pools.NewResourcePool(factory, 10, 30, 1*time.Hour, 1, nil)
    defer rp.Close()
    
    resource, err := rp.Get(context.Background())
    if err != nil {
        panic(err)
    }
    channel := resource.(*amqp.Channel)
    _ = channel.Publish("exchange", "routingKey", false, false, amqp.Publishing{
        Headers:         amqp.Table{},
        ContentType:     "text/plain",
        ContentEncoding: "",
        Body:            []byte("body"),
        DeliveryMode:    amqp.Transient,
        Priority:        0,
    })
    rp.Put(channel)
}

```

# Example

- [mysql](./example/mysql.go)
- [gorm](./example/gorm.go)
- [amqp](./example/amqp.go)
- [hive](./example/hive.go)