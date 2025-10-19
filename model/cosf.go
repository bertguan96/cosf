package model

// RdsTaskCfg 是Redis中存储的任务配置
type RdsTaskCfg struct {
	Key       string `json:"key"`
	ExpireAt  string `json:"expire_at"`
	UserId    string `json:"user_id"`
	Bucket    string `json:"bucket"`
	Region    string `json:"region"`
	SecretKey string `json:"secret_key"`
	AccessKey string `json:"access_key"`
	Qps       int64  `json:"qps"`
}
