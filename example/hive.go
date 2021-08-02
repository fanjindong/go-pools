package main

import (
	"context"
	"github.com/beltran/gohive"
	"github.com/fanjindong/go-pools"
	"time"
)

func main() {
	factory := func(context.Context) (pools.Resource, error) {
		return gohive.Connect("hs2.example.com", 10000, "KERBEROS", gohive.NewConnectConfiguration())
	}

	rp := pools.NewResourcePool(factory, 10, 30, 1*time.Hour, 1, nil)
	defer rp.Close()

	resource, err := rp.Get(context.Background())
	if err != nil {
		panic(err)
	}
	conn := resource.(*gohive.Connection)

	cursor := conn.Cursor()
	defer cursor.Close()
	cursor.Exec(context.Background(), "SELECT * FROM myTable")
	if cursor.Err != nil {
		panic(cursor.Err)
	}

	rp.Put(conn)
}
