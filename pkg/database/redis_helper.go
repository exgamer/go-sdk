package database

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

// NewRedisHelper - Новый Хелпер для работы с редисом
func NewRedisHelper[E interface{}](redisClient *redis.Client) *RedisHelper[E] {
	return &RedisHelper[E]{
		redisClient: redisClient,
	}
}

// RedisHelper - Хелпер для работы с редисом
type RedisHelper[E interface{}] struct {
	redisClient *redis.Client
	result      E
}

// GetByModel Возвращает значение по ключу
func (redisHelper *RedisHelper[E]) GetByModel(key string) (*E, error) {
	ctx := context.Background()
	val, err := redisHelper.redisClient.Get(ctx, key).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	if val == "" {
		return nil, nil
	}

	unMarshErr := json.Unmarshal([]byte(val), &redisHelper.result)

	if unMarshErr != nil {
		return nil, unMarshErr
	}

	return &redisHelper.result, nil
}

// SetByModel Записывает значение по ключу
func (redisHelper *RedisHelper[E]) SetByModel(key string, model *E, ttl time.Duration) error {
	jsonModel, err := json.Marshal(model)

	if err != nil {

		return err
	}

	ctx := context.Background()
	err = redisHelper.redisClient.Set(ctx, key, jsonModel, ttl).Err()

	if err != nil {
		return err
	}

	return nil
}
