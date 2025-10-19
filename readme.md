# CosF - 云存储QPS分配与限流服务

CosF 是一个基于 gRPC 和 HTTP 的云存储 QPS（每秒查询数）分配与限流服务，支持多区域、多桶的 QPS 管理和文件下载限流。

## 功能特性

- **QPS 分配管理**: 支持按用户、业务、桶和区域进行 QPS 分配
- **多级限流**: 支持桶级别和区域级别的 QPS 限制
- **Redis 限流**: 基于 Redis 的分布式限流实现
- **文件下载**: 支持云存储文件下载，带 QPS 限流
- **多协议支持**: 同时支持 gRPC 和 HTTP REST API
- **健康检查**: 内置健康检查接口

## 技术栈

- **Go 1.24+**: 主要开发语言
- **gRPC**: 高性能 RPC 框架
- **Protocol Buffers**: 接口定义
- **MySQL**: 数据持久化存储
- **Redis**: 缓存和限流
- **GORM**: ORM 框架
- **腾讯云 COS**: 云存储服务

## 项目结构

```
CosF/
├── api/                    # API 层
│   └── cosf.go           # gRPC 服务实现
├── common/                # 公共组件
│   ├── cos.go            # 腾讯云 COS 客户端
│   ├── db.go             # 数据库连接管理
│   ├── key_gen.go        # 密钥生成工具
│   ├── logger.go         # 日志配置
│   ├── rate_limit.go     # Redis 限流实现
│   └── rate_limit.lua    # Lua 限流脚本
├── cosf/                  # Protocol Buffers 定义
│   ├── cosf.proto        # 接口定义
│   ├── cosf.pb.go        # 生成的 Go 代码
│   ├── cosf_grpc.pb.go   # gRPC 服务代码
│   └── cosf.pb.gw.go     # HTTP 网关代码
├── model/                 # 数据模型
│   ├── base.go           # 基础模型
│   ├── cosf.go           # Redis 任务配置
│   └── t_cosf.go         # 数据库表模型
├── service/               # 业务逻辑层
│   ├── cosf.go           # QPS 分配服务
│   └── download.go       # 文件下载服务
├── main.go               # 应用入口
├── start.sh              # 启动脚本
├── code_gen.sh           # 代码生成脚本
├── Dockerfile            # Docker 镜像构建文件
├── docker-compose.yml    # Docker Compose 开发环境配置
├── docker-compose.prod.yml # Docker Compose 生产环境配置
├── .dockerignore         # Docker 构建忽略文件
├── env.example           # 环境变量配置示例
└── README.md             # 项目说明
```

## 数据模型

### 核心实体

- **CosfBucket**: 存储桶信息，包含 QPS 限制和访问密钥
- **CosfRegion**: 区域信息，包含区域级 QPS 限制
- **CosfBusiness**: 业务信息
- **CosfTask**: QPS 分配任务，包含用户、QPS、过期时间等

### 限流机制

- **桶级限流**: 每个存储桶有独立的 QPS 限制
- **区域级限流**: 每个区域有总的 QPS 限制
- **Redis 限流**: 基于滑动窗口的分布式限流

## API 接口

### gRPC 服务

#### 1. 健康检查
```protobuf
rpc HealthCheck(HealthCheckRequest) returns (HealthCheckResponse)
```

#### 2. QPS 分配
```protobuf
rpc AllocateQPS(AllocateQPSRequest) returns (AllocateQPSResponse)
```

**请求参数**:
- `user_id`: 用户ID
- `qps`: 申请的QPS数量
- `expire_at`: 过期时间（支持时间戳和RFC3339格式）
- `bucket_id`: 存储桶ID
- `business_id`: 业务ID

**响应数据**:
- `key`: 分配的唯一密钥
- `qps`: 分配的QPS数量
- `expire_at`: 过期时间

#### 3. 文件下载
```protobuf
rpc Download(DownloadRequest) returns (DownloadResponse)
```

**请求参数**:
- `key`: 分配密钥
- `cos_key`: 云存储文件路径
- `bucket_id`: 存储桶ID

### HTTP REST API

服务同时提供 HTTP REST API，通过 gRPC-Gateway 自动生成：

- `GET /cosf/healthCheck` - 健康检查
- `POST /cosf/allocateQPS` - QPS 分配
- `POST /cosf/download` - 文件下载

## 快速开始

### 环境要求

- Go 1.24+
- MySQL 5.7+
- Redis 6.0+
- Protocol Buffers 编译器

### Docker 部署（推荐）

使用 Docker 可以快速部署整个服务栈，包括 MySQL、Redis 和 CosF 应用。

#### 使用 Docker Compose

```bash
# 克隆项目
git clone <repository-url>
cd CosF

# 开发环境快速启动
docker-compose up -d

# 生产环境部署
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f cosf

# 停止服务
docker-compose down
```

#### 环境配置

```bash
# 复制环境变量模板
cp env.example .env

# 编辑环境变量
vim .env

# 使用环境变量启动
docker-compose --env-file .env up -d
```

#### 单独构建 Docker 镜像

```bash
# 构建镜像
docker build -t cosf:latest .

# 运行容器（需要先启动 MySQL 和 Redis）
docker run -d \
  --name cosf-app \
  -p 8000:8000 \
  -p 8080:8080 \
  --link mysql:mysql \
  --link redis:redis \
  cosf:latest
```

#### Docker 环境变量

可以通过环境变量配置数据库连接：

```bash
docker run -d \
  --name cosf-app \
  -p 8000:8000 \
  -p 8080:8080 \
  -e MYSQL_DSN="root:123456@tcp(mysql:3306)/cosf?charset=utf8mb4&parseTime=True&loc=Local" \
  -e REDIS_ADDR="redis:6379" \
  cosf:latest
```

### 安装依赖

```bash
# 安装 Protocol Buffers 工具
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest

# 安装项目依赖
go mod tidy
```

### 数据库配置

1. 创建 MySQL 数据库：
```sql
CREATE DATABASE cosf CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

2. 修改数据库连接配置（`common/db.go`）：
```go
dsn := "root:password@tcp(127.0.0.1:3306)/cosf?charset=utf8mb4&parseTime=True&loc=Local"
```

### Redis 配置

确保 Redis 服务运行在 `127.0.0.1:6379`，或修改 `common/db.go` 中的 Redis 配置。

### 代码生成

```bash
# 生成 Protocol Buffers 代码
chmod +x code_gen.sh
./code_gen.sh
```

### 启动服务

```bash
# 使用启动脚本
chmod +x start.sh
./start.sh

# 或直接运行
go run main.go
```

服务启动后：
- gRPC 服务运行在 `:8000`
- HTTP 服务运行在 `:8080`

## 使用示例

### 1. 申请 QPS

```bash
curl -X POST http://localhost:8080/cosf/allocateQPS \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "qps": 100,
    "expire_at": "1736688600",
    "bucket_id": "my-bucket",
    "business_id": "business123"
  }'
```

### 2. 下载文件

```bash
curl -X POST http://localhost:8080/cosf/download \
  -H "Content-Type: application/json" \
  -d '{
    "key": "your-allocated-key",
    "cos_key": "path/to/file.txt",
    "bucket_id": "my-bucket"
  }'
```

### 3. 健康检查

```bash
curl http://localhost:8080/cosf/healthCheck
```

## 配置说明

### 时间格式支持

`expire_at` 字段支持多种时间格式：
- Unix 时间戳（秒）: `1736688600`
- Unix 时间戳（毫秒）: `1736688600000`
- RFC3339 格式: `2025-10-12T14:30:00Z`
- ISO8601 格式: `2025-10-12T14:30:00Z`

### 限流配置

- 桶级 QPS 限制在 `CosfBucket` 表中配置
- 区域级 QPS 限制在 `CosfRegion` 表中配置
- Redis 限流使用滑动窗口算法

## 开发指南

### 添加新的 API

1. 在 `cosf/cosf.proto` 中定义新的 RPC 方法
2. 运行 `./code_gen.sh` 生成代码
3. 在 `api/cosf.go` 中实现服务方法
4. 在 `service/` 目录中添加业务逻辑

### 数据库迁移

使用 GORM 的 AutoMigrate 功能自动创建和更新表结构。在 `common/db.go` 中的 `InitMySQL` 函数中配置需要迁移的模型。

## Docker 部署详解

### 项目文件说明

- `Dockerfile`: 多阶段构建，优化镜像大小
- `docker-compose.yml`: 一键部署整个服务栈
- `.dockerignore`: 优化构建过程，排除不必要的文件

### Docker 镜像特性

- **多阶段构建**: 减少最终镜像大小
- **非 root 用户**: 提高安全性
- **健康检查**: 自动监控服务状态
- **时区配置**: 支持中国时区
- **日志持久化**: 日志文件挂载到宿主机

### 生产环境部署

#### 1. 使用 Docker Compose（推荐）

```bash
# 生产环境配置
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

#### 2. 使用 Kubernetes

```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cosf
spec:
  replicas: 3
  selector:
    matchLabels:
      app: cosf
  template:
    metadata:
      labels:
        app: cosf
    spec:
      containers:
      - name: cosf
        image: cosf:latest
        ports:
        - containerPort: 8000
        - containerPort: 8080
        env:
        - name: MYSQL_DSN
          valueFrom:
            secretKeyRef:
              name: cosf-secrets
              key: mysql-dsn
        - name: REDIS_ADDR
          value: "redis-service:6379"
```

#### 3. 环境变量配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `MYSQL_DSN` | `root:123456@tcp(127.0.0.1:3306)/cosf?charset=utf8mb4&parseTime=True&loc=Local` | MySQL 连接字符串 |
| `REDIS_ADDR` | `127.0.0.1:6379` | Redis 地址 |
| `GRPC_GO_LOG_VERBOSITY_LEVEL` | `99` | gRPC 日志详细级别 |
| `GRPC_GO_LOG_SEVERITY_LEVEL` | `info` | gRPC 日志级别 |

### 监控和日志

- 服务日志输出到 `logs/grpc.log`
- 支持 gRPC 日志级别配置
- Redis 限流状态监控
- Docker 健康检查自动监控服务状态

## 许可证

本项目采用 MIT 许可证。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进项目。

## 联系方式

如有问题，请通过 Issue 或邮件联系项目维护者。