package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/bertguan96/cosf/common"
	"golang.org/x/sync/errgroup"

	pb "github.com/bertguan96/cosf/cosf"
)

type DownloadService struct {
}

func NewDownloadService() *DownloadService {
	return &DownloadService{}
}

func (s *DownloadService) Download(ctx context.Context, request *pb.DownloadRequest) (string, error) {

}

// DownloadTask 下载任务
type DownloadTask struct {
	ObjectKey string // COS 中的 key，如 "folder/file.jpg"
	LocalPath string // 本地保存路径，如 "./downloads/file.jpg"
}

// DownloadOptions 下载选项
type DownloadOptions struct {
	MaxConcurrency int           // 最大并发数，建议 3~10
	Timeout        time.Duration // 单个下载超时，默认 5 分钟
	RetryCount     int           // 失败重试次数（可选）
}

func DefaultDownloadOptions() *DownloadOptions {
	return &DownloadOptions{
		MaxConcurrency: 5,
		Timeout:        5 * time.Minute,
		RetryCount:     2,
	}
}

// DownloadFiles 并发下载多个文件，控制并发数
func DownloadFiles(ctx context.Context, tasks []DownloadTask, opts *DownloadOptions) error {
	if opts == nil {
		opts = DefaultDownloadOptions()
	}

	// 使用 errgroup 限制并发 + 收集错误
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(opts.MaxConcurrency) // Go 1.19+ 支持 SetLimit，否则用 channel 限流

	for _, task := range tasks {
		task := task // 避免闭包陷阱
		g.Go(func() error {
			downloadCtx, cancel := context.WithTimeout(ctx, opts.Timeout)
			defer cancel()

			return downloadWithRetry(downloadCtx, task, opts.RetryCount)
		})
	}

	return g.Wait()
}

// downloadWithRetry 带重试的单文件下载
func downloadWithRetry(ctx context.Context, task DownloadTask, retry int) error {
	var lastErr error
	for i := 0; i <= retry; i++ {
		if i > 0 {
			time.Sleep(time.Second * time.Duration(i)) // 指数退避可选
		}

		err := downloadSingle(ctx, task)
		if err == nil {
			fmt.Printf("✅ 下载成功: %s -> %s\n", task.ObjectKey, task.LocalPath)
			return nil
		}
		lastErr = err
		fmt.Printf("⚠️ 下载失败 (第 %d 次): %s, 错误: %v\n", i+1, task.ObjectKey, err)
	}
	return fmt.Errorf("下载 %s 失败，重试 %d 次后仍失败: %w", task.ObjectKey, retry, lastErr)
}

// downloadSingle 单文件下载
func downloadSingle(ctx context.Context, task DownloadTask) error {
	// 确保本地目录存在
	dir := filepath.Dir(task.LocalPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	resp, err := common.GetCosClient().Object.Get(ctx, task.ObjectKey, nil)
	if err != nil {
		return fmt.Errorf("COS 获取对象失败: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(task.LocalPath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
