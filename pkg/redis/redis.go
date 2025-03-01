package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis represents a wrapper around a Redis client.
type Redis struct {
	*redis.Client
}

func connectToRedis(ctx context.Context, conn string) (*redis.Client, error) {
	opts, err := redis.ParseURL(conn)
	if err != nil {
		log.Fatal("cannot parse redis url: ", err)
	}

	client := redis.NewClient(opts)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}

// New initializes a new Redis instance with the provided context and configuration.
func New(ctx context.Context, cfg Config) (*Redis, error) {
	cfg, err := mergeWithDefault(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to merge with default config: %w", err)
	}

	var client *redis.Client
	for cfg.ConnAttempts > 0 {
		client, err = connectToRedis(ctx, cfg.ConnURL)
		if err == nil {
			break
		}

		log.Printf("Redis is trying to connect, attempts left: %d, err: %s\n", cfg.ConnAttempts, err)
		time.Sleep(cfg.ConnTimeout)

		cfg.ConnAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Redis{
		Client: client,
	}, nil
}

// Close closes the Redis client connection.
func (r *Redis) Close() {
	if r.Client != nil {
		r.Client.Close()
	}
}
