package api

import (
	"context"
	"fmt"
	"sync"

	pb "github.com/bertguan96/cosf/cosf"
	"github.com/bertguan96/cosf/model"
	"github.com/bertguan96/cosf/service"
)

type CosfServer struct {
	pb.UnimplementedCosfServiceServer
	mu sync.Mutex // protects routeNotes
}

var (
	cosfService     = service.NewCosfService()     // cosf service
	downloadService = service.NewDownloadService() // download service
)

// HealthCheck 健康检查
func (s *CosfServer) HealthCheck(ctx context.Context, request *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	return &pb.HealthCheckResponse{
		Base: &pb.BaseResponse{
			Code:    int32(model.CodeOK),
			Message: "OK",
		},
	}, nil
}

// AllocateQPS 申请QPS
func (s *CosfServer) AllocateQPS(ctx context.Context, request *pb.AllocateQPSRequest) (*pb.AllocateQPSResponse, error) {
	resp, err := cosfService.AllocateQPS(ctx, request)
	if err != nil {
		return &pb.AllocateQPSResponse{
			Base: &pb.BaseResponse{
				Code:    int32(model.CodeError),
				Message: fmt.Sprintf("AllocateQPS failed: %s", err.Error()),
			},
		}, nil
	}

	return &pb.AllocateQPSResponse{
		Base: &pb.BaseResponse{
			Code:    int32(model.CodeOK),
			Message: "AllocateQPS success",
		},
		Data: resp,
	}, nil
}

// Download 下载文件
func (s *CosfServer) Download(ctx context.Context, request *pb.DownloadRequest) (*pb.DownloadResponse, error) {
	// 验证请求参数
	if request.Key == "" {
		return &pb.DownloadResponse{
			Base: &pb.BaseResponse{
				Code:    int32(model.CodeInvalidKey),
				Message: string(model.CodeMessageInvalidKey),
			},
		}, nil
	}
	if request.CosKey == "" {
		return &pb.DownloadResponse{
			Base: &pb.BaseResponse{
				Code:    int32(model.CodeInvalidCosKey),
				Message: string(model.CodeMessageInvalidCosKey),
			},
		}, nil
	}
	if request.BucketId == "" {
		return &pb.DownloadResponse{
			Base: &pb.BaseResponse{
				Code:    int32(model.CodeInvalidBucketId),
				Message: string(model.CodeMessageInvalidBucketId),
			},
		}, nil
	}
	resp, err := downloadService.Download(ctx, request)
	if err != nil {
		return &pb.DownloadResponse{
			Base: &pb.BaseResponse{
				Code:    int32(model.CodeError),
				Message: fmt.Sprintf("Download failed: %s", err.Error()),
			},
		}, nil
	}
	return &pb.DownloadResponse{
		Base: &pb.BaseResponse{
			Code:    int32(model.CodeOK),
			Message: "Download success",
		},
		Content: resp,
	}, nil
}
