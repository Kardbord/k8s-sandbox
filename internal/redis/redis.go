package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const jobQueueKey = "job_queue"

func NewClient() (*redis.Client, error) {
	addr := "host.docker.internal:6379"
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return rdb, nil
}

func PushJob(ctx context.Context, rdb *redis.Client, jobID string) error {
	return rdb.RPush(ctx, jobQueueKey, jobID).Err()
}

func PopJob(ctx context.Context, rdb *redis.Client, timeout time.Duration) (string, error) {
	result, err := rdb.BLPop(ctx, timeout, jobQueueKey).Result()
	if err != nil {
		// Note that when BLPop times out, err == redis.Nil
		return "", err
	}

	return result[1], nil
}
