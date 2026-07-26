package bootstrap

import "github.com/redis/go-redis/v9"

func NewRedisClient(addr string, user string, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: user,
		Password: password,
	})
}
