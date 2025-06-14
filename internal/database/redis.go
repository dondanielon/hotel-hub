package database

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

type Redis struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedis() *Redis {
	return &Redis{
		client: redis.NewClient(
			&redis.Options{
				Addr:     os.Getenv("REDIS_CONNECTION_STRING"),
				Username: os.Getenv("REDIS_USERNAME"),
				Password: os.Getenv("REDIS_PASSWORD"),
				DB:       0,
			},
		),
		ctx: context.Background(),
	}
}

func (r *Redis) SetCache(key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("Error marshaling value: %v", err)
	}

	err = r.client.Set(r.ctx, key, jsonValue, expiration).Err()
	if err != nil {
		return fmt.Errorf("Error setting cache: %v", err)
	}

	return nil
}

func (r *Redis) GetCache(key string, dest interface{}) error {
	val, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// key doesn't exist
			return nil
		}
		return fmt.Errorf("Error getting cache: %v", err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("Error unmarshaling value: %v", err)
	}

	return nil
}

func (r *Redis) DeleteCache(key string) error {
	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		return fmt.Errorf("Error deleting cache: %v", err)
	}

	return nil
}
