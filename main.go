package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // or use your container IP if needed
		// Password: "",         // set password if needed
		DB: 0, // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("key:", val)

	// If key doesn't exist
	val2, err := rdb.Get(ctx, "missing").Result()
	if err == redis.Nil {
		fmt.Println("missing key does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("missing:", val2)
	}
}
