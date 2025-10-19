// logger.go
package common

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"google.golang.org/grpc/grpclog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// fileLogger 实现 grpclog.LoggerV2 接口
type fileLogger struct {
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
	mu      sync.Mutex // 确保并发安全（可选，log.Logger 本身是并发安全的）
}

// Errorln implements grpclog.LoggerV2.
func (l *fileLogger) Errorln(args ...any) {
	l.error.Println(args...)
}

// Fatal implements grpclog.LoggerV2.
func (l *fileLogger) Fatal(args ...any) {
	l.error.Println(args...)
}

// Fatalf implements grpclog.LoggerV2.
func (l *fileLogger) Fatalf(format string, args ...any) {
	l.error.Printf(format, args...)
}

// Fatalln implements grpclog.LoggerV2.
func (l *fileLogger) Fatalln(args ...any) {
	l.error.Println(args...)
}

// Infoln implements grpclog.LoggerV2.
func (l *fileLogger) Infoln(args ...any) {
	l.info.Println(args...)
}

// Warningln implements grpclog.LoggerV2.
func (l *fileLogger) Warningln(args ...any) {
	l.warning.Println(args...)
}

func NewFileLogger(filePath string) (*fileLogger, error) {
	// lumberjack 自动轮转
	fileWriter := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    100, // MB
		MaxBackups: 5,
		MaxAge:     28, // days
		Compress:   true,
	}

	// 创建同时写入控制台和文件的 MultiWriter
	multiWriter := io.MultiWriter(os.Stdout, fileWriter)

	flags := log.LstdFlags | log.Lmicroseconds | log.Lshortfile
	return &fileLogger{
		info:    log.New(multiWriter, "[INFO] ", flags),
		warning: log.New(multiWriter, "[WARN] ", flags),
		error:   log.New(multiWriter, "[ERROR] ", flags),
	}, nil
}

// 实现 grpclog.LoggerV2 接口
func (l *fileLogger) Info(args ...interface{}) {
	l.info.Output(2, fmt.Sprint(args...))
}

func (l *fileLogger) Infof(format string, args ...interface{}) {
	l.info.Output(2, fmt.Sprintf(format, args...))
}

func (l *fileLogger) Warning(args ...interface{}) {
	l.warning.Output(2, fmt.Sprint(args...))
}

func (l *fileLogger) Warningf(format string, args ...interface{}) {
	l.warning.Output(2, fmt.Sprintf(format, args...))
}

func (l *fileLogger) Error(args ...interface{}) {
	l.error.Output(2, fmt.Sprint(args...))
}

func (l *fileLogger) Errorf(format string, args ...interface{}) {
	l.error.Output(2, fmt.Sprintf(format, args...))
}

// V 控制 verbose 日志是否输出（返回 true 表示输出）
// 你可以根据需要返回 true/false，或读取配置
func (l *fileLogger) V(level int) bool {
	// 例如：只输出 level <= 2 的 verbose 日志
	return level <= 2
}

// initLogger 初始化日志
func InitLogger() {
	logger, err := NewFileLogger("logs/grpc.log")
	if err != nil {
		grpclog.Fatalf("failed to create logger: %v", err)
	}
	grpclog.SetLoggerV2(logger)
}
