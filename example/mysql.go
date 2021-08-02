package main

import (
	"context"
	"database/sql"
	"github.com/fanjindong/go-pools"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

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
	_, err = db.Query("SELECT * FROM table")
	if err != nil {
		panic(err.Error())
	}

	rp.Put(db)
}
