package bootstrap

import "github.com/redis/go-redis/v9"

func NewRedisClient(addr, password, username string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		Username: username,
	})
}
