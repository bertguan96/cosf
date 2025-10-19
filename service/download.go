package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/bertguan96/cosf/common"
	pb "github.com/bertguan96/cosf/cosf"
	"github.com/bertguan96/cosf/model"
	"google.golang.org/grpc/grpclog"
)

type DownloadService struct {
}

func NewDownloadService() *DownloadService {
	return &DownloadService{}
}

func (s *DownloadService) Download(ctx context.Context, request *pb.DownloadRequest) (string, error) {
	rds := common.GetRedis()
	// 判断ke是否过期
	expireAt, err := rds.TTL(ctx, request.Key).Result()
	if err != nil {
		return "", errors.New("get key ttl failed," + err.Error())
	}
	if expireAt <= 0 {
		return "", errors.New("key expired")
	}
	result, err := rds.Get(ctx, request.Key).Result()
	grpclog.Info("download result:", result, err)
	if err != nil {
		return "", err
	}
	if result == "" {
		return "", errors.New("task cfg not found")
	}
	var taskCfg *model.RdsTaskCfg
	if err = json.Unmarshal([]byte(result), &taskCfg); err != nil {
		return "", errors.New("unmarshal task cfg failed," + err.Error())
	}

	// 限流
	limiter := common.NewRedisRateLimiter(rds, taskCfg.Qps) //
	allowed, err := limiter.Allow(ctx, request.Key)
	if err != nil {
		return "", errors.New("rate limit check failed," + err.Error())
	}
	if !allowed {
		return "", errors.New("rate limit exceeded, please slow down")
	}

	cosClient := common.InitCos(taskCfg.SecretKey, taskCfg.AccessKey, taskCfg.Bucket, taskCfg.Region)
	return common.Download(ctx, cosClient, request.CosKey)
}
