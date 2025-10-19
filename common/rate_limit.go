package common

import (
	"context"
	"errors"
	"fmt"
	"time"

	_ "embed"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/grpclog"
)

var (
	ErrRateLimited = errors.New("too many requests, please slow down")
)

type RedisRateLimiter struct {
	rdb      *redis.Client
	limit    int64         // 每秒最大请求数
	interval time.Duration // 窗口大小（通常为 1s）
}

func NewRedisRateLimiter(rdb *redis.Client, limit int64) *RedisRateLimiter {
	return &RedisRateLimiter{
		rdb:      rdb,
		limit:    limit,
		interval: time.Second,
	}
}

//go:embed rate_limit.lua
var rateLimitScript string

// Allow 检查是否允许请求，返回当前计数和是否允许
func (r *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	now := time.Now().Unix()
	windowKey := fmt.Sprintf("rate_limit:%s:%d", key, now)

	result, err := r.rdb.Eval(ctx, rateLimitScript, []string{windowKey}, r.limit, int64(r.interval.Seconds())).Result()
	if err != nil {
		return false, fmt.Errorf("限流检查失败: %w", err)
	}

	// Lua 返回 0 表示超限
	if result == int64(0) {
		return false, ErrRateLimited
	}

	// 每10min打印一次当前qps
	if time.Now().Minute()%10 == 0 {
		grpclog.Infof("current qps: %d", result)
	}

	return true, nil
}
