package config

import "github.com/redis/go-redis/v9"

func NewRedisClient(env Env) *redis.Client {
	rd := redis.NewClient(&redis.Options{
		Addr:     env.RedisAddress,
		Password: env.RedisPassword,
		DB:       env.RedisDB,
	})

	return rd
}
