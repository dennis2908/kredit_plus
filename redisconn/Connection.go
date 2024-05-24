package redisconn

import (
	redis "github.com/redis/go-redis/v9"
)

func Connect() (*redis.Client, error) {

	return redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		Protocol: 3,
	}), nil
}
