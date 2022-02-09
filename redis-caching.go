package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

const CACHE_DURATION = 2 * time.Second

var rdb *redis.Client
var ctx = context.Background()

func conRedis() {

	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

}

func getFromCache(key string) (album, error) {
	var alb album
	res, err := rdb.Get(ctx, key).Result()

	if err == redis.Nil {
		return alb, redis.Nil
	} else if err != nil {
		panic(err)
	}

	bytes, err := hex.DecodeString(res)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &alb)
	if err != nil {
		panic(err)
	}
	return alb, nil

}

func storeInCache(key string, alb album) {
	bytes, err := json.Marshal(alb)
	if err != nil {
		panic(err)
	}
	val := hex.EncodeToString(bytes)
	err = rdb.Set(ctx, key, val, CACHE_DURATION).Err()
	if err != nil {
		panic(err)
	}
}
