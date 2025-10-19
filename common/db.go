package common

import (
	"context"
	"sync"
	"time"

	"github.com/bertguan96/cosf/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/grpclog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once sync.Once
	db   *gorm.DB
	rdb  *redis.Client
	dsn  string = "root:123456@tcp(127.0.0.1:3306)/cosf?charset=utf8mb4&parseTime=True&loc=Local" // 链接地址
)

// InitMySQL 初始化 MySQL 连接
func InitMySQL(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Info),
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		grpclog.Fatalf("failed to connect MySQL: %v", err)
	}
	err = db.AutoMigrate(&model.CosfBucket{}, &model.CosfBusiness{}, &model.CosfTask{}, &model.CosfRegion{})
	if err != nil {
		grpclog.Fatalf("failed to migrate database: %v", err)
	}

	grpclog.Info("MySQL connected successfully")
	return db
}

// InitRedis 初始化 Redis 连接
func InitRedis() *redis.Client {
	Rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379", // Redis 地址
		Password: "",               // 无密码
		DB:       0,                // 默认 DB 0
		PoolSize: 20,               // 连接池大小（默认 10 * CPU 数）
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Rdb.Ping(ctx).Err(); err != nil {
		grpclog.Fatalf("failed to connect Redis: %v", err)
	}
	grpclog.Info("Redis connected successfully")
	return Rdb
}

func InitDB() {
	once.Do(func() {
		db = InitMySQL(dsn)
		rdb = InitRedis()
	})
	grpclog.Info("MySQL connected successfully")
}

func GetMysql() *gorm.DB {
	return db
}

func GetRedis() *redis.Client {
	return rdb
}
