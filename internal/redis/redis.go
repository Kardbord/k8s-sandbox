package redis

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const jobQueueKey = "job_queue"

func NewClient() (*redis.Client, error) {
	rdsAddr := os.Getenv("REDIS_ADDR")
	if rdsAddr == "" {
		rdsAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{Addr: rdsAddr})
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
