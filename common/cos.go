package common

import (
	"net/http"
	"net/url"
	"os"

	"github.com/tencentyun/cos-go-sdk-v5"
	"google.golang.org/grpc/grpclog"
)

var Client *cos.Client

func initCos() *cos.Client {
	secretID := os.Getenv("COS_SECRET_ID")
	secretKey := os.Getenv("COS_SECRET_KEY")
	bucket := os.Getenv("COS_BUCKET")
	region := os.Getenv("COS_REGION")

	u, _ := url.Parse("https://" + bucket + ".cos." + region + ".myqcloud.com")
	baseURL := &cos.BaseURL{BucketURL: u}
	Client = cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	grpclog.Info("COS connected successfully")
	return Client
}

func InitCos() {
	once.Do(func() {
		Client = initCos()
	})
	grpclog.Info("COS connected successfully")
}

func GetCosClient() *cos.Client {
	// 如果CLient端开
	if Client == nil {
		once.Do(func() {
			Client = initCos()
		})
		grpclog.Info("COS reconnected successfully")
	}
	return Client
}
