package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
}

func Open(addr, password string, db int) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Redis{Client: client}
}

func (r *Redis) Close() error {
	return r.Client.Close()
}

func (r *Redis) SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	return r.Client.SetNX(ctx, key, value, ttl).Result()
}

func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.Client.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *Redis) Del(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *Redis) Incr(ctx context.Context, key string, ttl time.Duration) (int64, error) {
	pipe := r.Client.TxPipeline()
	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Val(), nil
}
