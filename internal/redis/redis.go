package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

const jobQueueKey = "job_queue"

func NewClient() (*redis.Client, error) {
	addr := "localhost:6379"
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}

func PushJob(ctx context.Context, rdb *redis.Client, jobID string) error {
	return rdb.RPush(ctx, jobQueueKey, jobID).Err()
}
