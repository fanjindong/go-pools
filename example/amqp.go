package main

import (
	"context"
	"github.com/fanjindong/go-pools"
	"github.com/streadway/amqp"
	"time"
)

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
