package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	TimeWindow time.Duration
	Limit      int
	Key        string
	Redis      *redis.Client
}

func NewRateLimiter(redisClient *redis.Client, key string, timeWindow time.Duration, limit int) *RateLimiter {
	return &RateLimiter{
		TimeWindow: timeWindow,
		Limit:      limit,
		Key:        key,
		Redis:      redisClient,
	}
}

// AllowRequest checks if the request can be processed using Redis
func (r *RateLimiter) AllowRequest(ctx context.Context) bool {
	now := time.Now().UnixNano()

	// Start a pipeline for Redis commands
	pipe := r.Redis.TxPipeline()

	// Remove expired timestamps
	expiryThreshold := now - r.TimeWindow.Nanoseconds()
	pipe.ZRemRangeByScore(ctx, r.Key, "0", fmt.Sprint(expiryThreshold))

	// Count the number of valid timestamps
	countCmd := pipe.ZCard(ctx, r.Key)

	// Add the current timestamp to the sorted set with a TTL
	pipe.ZAdd(ctx, r.Key, redis.Z{
		Score:  float64(now),
		Member: now,
	})
	pipe.Expire(ctx, r.Key, r.TimeWindow)

	// Execute the pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		return false // Consider throttling in case of Redis errors
	}

	// Check if the number of requests exceeds the limit
	return countCmd.Val() < int64(r.Limit)
}
