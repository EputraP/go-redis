package main

import (
	"context"
	"fmt"
	"time"

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
	err = rdb.Set(ctx, "key1", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	err = rdb.Set(ctx, "key2", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	err = rdb.Set(ctx, "key3", "value", 0).Err()
	if err != nil {
		panic(err)
	}
	err = rdb.Set(ctx, "halo", "satu", 0).Err()
	if err != nil {
		panic(err)
	}
	if err = rdb.Del(ctx, "key1").Err(); err != nil {
		panic(err)
	}

	if err = rdb.Append(ctx, "halo", " dua, tiga, kucing berlari").Err(); err != nil {
		panic(err)
	}

	valKeyExist, err := rdb.Exists(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key exists:", valKeyExist)

	keys, err := rdb.Keys(ctx, "*key*").Result()
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		fmt.Println("Matched key:", key)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("key:", val)

	if err = rdb.SetRange(ctx, "halo", 14, " empat lima enam tuju").Err(); err != nil {
		panic(err)
	}

	getRange, err := rdb.GetRange(ctx, "halo", 14, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("get range:", getRange)

	if err = rdb.MSet(ctx, "key2", "value2", "key3", "value3").Err(); err != nil {
		panic(err)
	}

	mgetResult, err := rdb.MGet(ctx, "key2", "key3").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("mget result:", mgetResult)

	if err := rdb.Expire(ctx, "key", 1*time.Minute).Err(); err != nil {
		panic(err)
	}

	expireAt := time.Now().Add(1 * time.Hour)
	if err := rdb.ExpireAt(ctx, "key2", expireAt).Err(); err != nil {
		panic(err)
	}

	ttl, err := rdb.TTL(ctx, "key2").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("TTL:", ttl) // Example: 9.899s

	if err = rdb.Incr(ctx, "counter").Err(); err != nil {
		panic(err)
	}

	if err = rdb.Decr(ctx, "counter").Err(); err != nil {
		panic(err)
	}
	if err = rdb.IncrBy(ctx, "counter", 10).Err(); err != nil {
		panic(err)
	}
	if err = rdb.DecrBy(ctx, "counter", 5).Err(); err != nil {
		panic(err)
	}

	// rdb.FlushDB(ctx)
	// rdb.FlushAll(ctx)

	pipe := rdb.Pipeline()

	set1 := pipe.Set(ctx, "x", "100", 0)
	set2 := pipe.Set(ctx, "y", "200", 0)
	get1 := pipe.Get(ctx, "x")

	_, err = pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println("SET x:", set1.Val())
	fmt.Println("SET y:", set2.Val())
	fmt.Println("GET x:", get1.Val())

	// Transaction
	Txpipe := rdb.TxPipeline()

	Txpipe.Set(ctx, "a", 1, 0)
	Txpipe.Incr(ctx, "a")

	cmds, err := Txpipe.Exec(ctx)
	if err != nil {
		panic(err)
	}

	for _, cmd := range cmds {
		fmt.Println(cmd.Name(), cmd.Args())
	}

	err = rdb.Watch(ctx, func(tx *redis.Tx) error {
		val, err := tx.Get(ctx, "balance").Int()
		if err != nil && err != redis.Nil {
			return err
		}

		if val < 100 {
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.IncrBy(ctx, "balance", 50)
				return nil
			})
			return err
		}

		return nil
	}, "balance")

	// If key doesn't exist
	val2, err := rdb.Get(ctx, "key").Result()
	if err == redis.Nil {
		fmt.Println("missing key does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("missing:", val2)
	}
}
