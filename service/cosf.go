package service

import (
	"context"

	"errors"
	"time"

	"github.com/bertguan96/cosf/common"
	mysql "github.com/bertguan96/cosf/common"
	pb "github.com/bertguan96/cosf/cosf"
	"github.com/bertguan96/cosf/model"
)

type CosfService struct {
}

func NewCosfService() *CosfService {
	return &CosfService{}
}

func (s *CosfService) AllocateQPS(ctx context.Context, request *pb.AllocateQPSRequest) (*pb.AllocateQPSResponseData, error) {
	db := mysql.GetMysql()
	// 获取当前桶的情况
	bucket := model.CosfBucket{}
	if err := db.Model(&model.CosfBucket{}).Where("bucket_id = ?", request.BucketId).First(&bucket).Error; err != nil {
		return nil, errors.New("bucket not found")
	}
	// 单桶限制
	var taskList []*model.CosfTask
	if err := db.Model(&model.CosfTask{}).Where("bucket_id = ?", request.BucketId).Find(&taskList).Error; err != nil {
		return nil, errors.New("get task list failed")
	}
	var totalQps int // 当前已有的总qps
	for _, task := range taskList {
		totalQps += task.Qps
	}
	// 如果超过qps，则返回错误
	if totalQps+int(request.Qps) > bucket.BucketQps {
		return nil, errors.New("qps exceed")
	}

	// 当前区域的QPS是否超过上限
	bucketList := make([]*model.CosfBucket, 0)
	if err := db.Model(&model.CosfBucket{}).Where("region_id = ?", bucket.RegionId).Find(&bucketList).Error; err != nil {
		return nil, errors.New("get bucket list failed")
	}
	// 获取当前地区QPS的限制
	var regionQps int64 // 当前地区的QPS
	if err := db.Model(&model.CosfRegion{}).Where("id = ?", bucket.RegionId).Select("region_qps").First(&regionQps).Error; err != nil {
		return nil, errors.New("get region failed " + err.Error())
	}
	// 如果当前地区的QPS超过上限，则返回错误
	if totalQps+int(request.Qps) > int(regionQps) {
		return nil, errors.New("region qps exceed, please try again later.")
	}

	// 尝试多种时间格式解析
	var expireAt time.Time
	var err error

	// 首先尝试RFC3339格式
	expireAt, err = time.Parse(time.RFC3339, request.ExpireAt)
	if err != nil {
		// 如果失败，尝试RFC3339Nano格式
		expireAt, err = time.Parse(time.RFC3339Nano, request.ExpireAt)
		if err != nil {
			// 如果还是失败，尝试ISO8601格式（只有Z没有时区偏移）
			expireAt, err = time.Parse("2006-01-02T15:04:05Z", request.ExpireAt)
			if err != nil {
				return nil, errors.New("invalid expire at, " + err.Error())
			}
		}
	}
	newTask := model.CosfTask{
		UserID:     request.UserId,
		Qps:        int(request.Qps),
		ExpireAt:   expireAt,
		BucketID:   request.BucketId,
		Key:        common.Generate10CharID(),
		BusinessID: request.BusinessId,
	}
	err = db.Create(&newTask).Error
	if err != nil {
		return nil, errors.New("create task failed")
	}

	resp := &pb.AllocateQPSResponseData{
		Key:      newTask.Key,
		Qps:      int64(newTask.Qps),
		ExpireAt: newTask.ExpireAt.Format(time.RFC3339),
	}
	return resp, nil
}
