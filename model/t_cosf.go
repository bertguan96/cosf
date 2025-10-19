package model

import "time"

type CosfBucket struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BucketID   string    `gorm:"column:bucket_id" json:"bucket_id"`
	BucketName string    `gorm:"column:bucket_name" json:"bucket_name"`
	Region     string    `gorm:"column:region" json:"region"`
	BucketQps  int       `gorm:"column:bucket_qps" json:"bucket_qps"`
	RegionId   int       `gorm:"column:region_id" json:"region_id"`
	Owner      string    `gorm:"column:owner" json:"owner"`
	CreateAt   time.Time `gorm:"column:create_at;type:datetime;" json:"create_at"`
	Supplier   string    `gorm:"column:supplier" json:"supplier"` // 供应商
	AccessKey  string    `gorm:"column:access_key" json:"access_key"` // 供应商的access key
	SecretKey  string    `gorm:"column:secret_key" json:"secret_key"` // 供应商的secret key
}

func (*CosfBucket) TableName() string {
	return "cosf_bucket"
}

type CosfBusiness struct {
	ID            int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	BusinessID    string    `gorm:"column:business_id" json:"business_id"`
	BusinessName  string    `gorm:"column:business_name" json:"business_name"`
	BusinessOwner string    `gorm:"column:business_owner" json:"business_owner"`
	CreateAt      time.Time `gorm:"column:create_at;type:datetime;" json:"create_at"`
}

func (*CosfBusiness) TableName() string {
	return "cosf_business"
}

type CosfTask struct {
	ID         int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID     string    `gorm:"column:user_id" json:"user_id"`
	Qps        int       `gorm:"column:qps" json:"qps"`
	ExpireAt   time.Time `gorm:"column:expire_at;type:datetime;" json:"expire_at"`
	BusinessID string    `gorm:"column:business_id" json:"business_id"`
	BucketID   string    `gorm:"column:bucket_id" json:"bucket_id"`
	Key        string    `gorm:"column:key" json:"key"`
}

func (*CosfTask) TableName() string {
	return "cosf_task"
}

type CosfRegion struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Region    string    `gorm:"column:region" json:"region"`
	RegionQps int       `gorm:"column:region_qps" json:"region_qps"`
	CreateAt  time.Time `gorm:"column:create_at;type:datetime;" json:"create_at"`
}

func (*CosfRegion) TableName() string {
	return "cosf_region"
}
