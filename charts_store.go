package itl

import "github.com/go-redis/redis"
import "log"

type ChartsStore interface {
	update(key, url string, weigth float64)
}

type RedisChartsStore struct {
	redisClient *redis.Client
}

func NewRedisChartsStore(redisURL string) ChartsStore {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
		DB:   0,
	})

	return &RedisChartsStore{
		redisClient: redisClient,
	}
}

func (r RedisChartsStore) update(key, url string, weigth float64) {
	_, err := r.redisClient.ZIncr(key, redis.Z{weigth, url}).Result()
	if err != nil {
		log.Println(err)
	}
}
