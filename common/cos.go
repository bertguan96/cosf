package common

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
	"google.golang.org/grpc/grpclog"
)

var Client *cos.Client

func InitCos(secretID, secretKey, bucket, region string) *cos.Client {
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

func Download(ctx context.Context, client *cos.Client, key string) (string, error) {
	resp, err := client.Object.Get(ctx, key, nil)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return string(body), nil
}
