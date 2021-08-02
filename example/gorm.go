package main

import (
	"context"
	"github.com/fanjindong/go-pools"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type DB struct {
	*gorm.DB
}

func (d *DB) Close() error {
	return nil
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	factory := func(context.Context) (pools.Resource, error) {
		gdb, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		return &DB{gdb}, nil
	}

	rp := pools.NewResourcePool(factory, 10, 30, 1*time.Hour, 1, nil)
	defer rp.Close()

	resource, err := rp.Get(context.Background())
	if err != nil {
		panic(err)
	}
	db := resource.(*DB)

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	db.Delete(&product, 1)

	rp.Put(db)
}
