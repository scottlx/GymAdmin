# TASK007 - 系统基础设施搭建

## 一、功能概述

### 1.1 功能目标
搭建健身房管理系统的完整基础设施，包括后端服务架构、数据库、缓存、文件存储、消息队列、日志监控等核心组件，为业务功能提供稳定可靠的技术支撑。

### 1.2 核心价值
- 提供高可用、高性能的系统架构
- 建立完善的监控和日志体系
- 实现系统的可扩展性和可维护性
- 保障数据安全和系统稳定性
- 支持快速迭代和部署

### 1.3 技术栈
- **后端**: Golang (Gin框架)
- **数据库**: MySQL 8.0
- **缓存**: Redis 6.0
- **文件存储**: 阿里云OSS / MinIO
- **消息队列**: RabbitMQ / Kafka (可选)
- **日志**: ELK (Elasticsearch + Logstash + Kibana)
- **监控**: Prometheus + Grafana
- **容器化**: Docker + Docker Compose
- **CI/CD**: GitLab CI / GitHub Actions

## 二、功能详细拆解

### 2.1 后端服务架构搭建

#### 2.1.1 项目结构设计
**输入**: 业务需求和技术选型
**输出**: 完整的项目目录结构

**执行内容**:
```
gym-admin-backend/
├── cmd/
│   └── server/
│       └── main.go                 # 程序入口
├── internal/
│   ├── config/                     # 配置管理
│   │   ├── config.go
│   │   └── config.yaml
│   ├── models/                     # 数据模型
│   │   ├── user.go
│   │   ├── coach.go
│   │   ├── membership_card.go
│   │   └── ...
│   ├── repository/                 # 数据访问层
│   │   ├── user_repo.go
│   │   ├── coach_repo.go
│   │   └── ...
│   ├── service/                    # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── coach_service.go
│   │   └── ...
│   ├── controller/                 # 控制器层
│   │   ├── user_controller.go
│   │   ├── coach_controller.go
│   │   └── ...
│   ├── middleware/                 # 中间件
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── rate_limit.go
│   ├── router/                     # 路由配置
│   │   └── router.go
│   └── pkg/                        # 公共包
│       ├── utils/                  # 工具函数
│       ├── errors/                 # 错误定义
│       ├── response/               # 响应封装
│       └── validator/              # 参数校验
├── pkg/                            # 外部可用的包
│   ├── cache/                      # 缓存封装
│   ├── database/                   # 数据库封装
│   ├── logger/                     # 日志封装
│   ├── oss/                        # 对象存储封装
│   └── mq/                         # 消息队列封装
├── scripts/                        # 脚本文件
│   ├── init_db.sql                 # 数据库初始化
│   └── deploy.sh                   # 部署脚本
├── tests/                          # 测试文件
│   ├── unit/                       # 单元测试
│   └── integration/                # 集成测试
├── docs/                           # 文档
│   ├── api/                        # API文档
│   └── design/                     # 设计文档
├── Dockerfile                      # Docker配置
├── docker-compose.yml              # Docker Compose配置
├── Makefile                        # Make命令
├── go.mod                          # Go模块
├── go.sum
└── README.md
```

**验收标准**:
- 项目结构清晰，职责分明
- 符合Go语言项目规范
- 支持模块化开发
- 易于维护和扩展

---

#### 2.1.2 核心框架搭建
**输入**: 项目结构
**输出**: 可运行的基础框架

**执行内容**:

**A. 主程序入口 (cmd/server/main.go)**
```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "gym-admin/internal/config"
    "gym-admin/internal/router"
    "gym-admin/pkg/database"
    "gym-admin/pkg/cache"
    "gym-admin/pkg/logger"
)

func main() {
    // 加载配置
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 初始化日志
    if err := logger.Init(cfg.Log); err != nil {
        log.Fatalf("Failed to init logger: %v", err)
    }

    // 初始化数据库
    if err := database.Init(cfg.Database); err != nil {
        logger.Fatal("Failed to init database", err)
    }
    defer database.Close()

    // 初始化Redis
    if err := cache.Init(cfg.Redis); err != nil {
        logger.Fatal("Failed to init redis", err)
    }
    defer cache.Close()

    // 初始化路由
    r := router.Setup(cfg)

    // 启动HTTP服务器
    srv := &http.Server{
        Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
        Handler:        r,
        ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
        WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

    // 优雅关闭
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("Failed to start server", err)
        }
    }()

    logger.Info(fmt.Sprintf("Server started on port %d", cfg.Server.Port))

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    logger.Info("Shutting down server...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal("Server forced to shutdown", err)
    }

    logger.Info("Server exited")
}
```

**B. 配置管理 (internal/config/config.go)**
```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    OSS      OSSConfig      `mapstructure:"oss"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
    Port         int    `mapstructure:"port"`
    Mode         string `mapstructure:"mode"`
    ReadTimeout  int    `mapstructure:"read_timeout"`
    WriteTimeout int    `mapstructure:"write_timeout"`
}

type DatabaseConfig struct {
    Host            string `mapstructure:"host"`
    Port            int    `mapstructure:"port"`
    Username        string `mapstructure:"username"`
    Password        string `mapstructure:"password"`
    Database        string `mapstructure:"database"`
    MaxIdleConns    int    `mapstructure:"max_idle_conns"`
    MaxOpenConns    int    `mapstructure:"max_open_conns"`
    ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}

type OSSConfig struct {
    Endpoint        string `mapstructure:"endpoint"`
    AccessKeyID     string `mapstructure:"access_key_id"`
    AccessKeySecret string `mapstructure:"access_key_secret"`
    BucketName      string `mapstructure:"bucket_name"`
}

type JWTConfig struct {
    Secret     string `mapstructure:"secret"`
    ExpireTime int    `mapstructure:"expire_time"`
}

type LogConfig struct {
    Level      string `mapstructure:"level"`
    Filename   string `mapstructure:"filename"`
    MaxSize    int    `mapstructure:"max_size"`
    MaxBackups int    `mapstructure:"max_backups"`
    MaxAge     int    `mapstructure:"max_age"`
    Compress   bool   `mapstructure:"compress"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./internal/config")
    viper.AddConfigPath(".")

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

**C. 路由配置 (internal/router/router.go)**
```go
package router

import (
    "gym-admin/internal/config"
    "gym-admin/internal/controller"
    "gym-admin/internal/middleware"

    "github.com/gin-gonic/gin"
)

func Setup(cfg *config.Config) *gin.Engine {
    if cfg.Server.Mode == "release" {
        gin.SetMode(gin.ReleaseMode)
    }

    r := gin.New()

    // 全局中间件
    r.Use(gin.Recovery())
    r.Use(middleware.Logger())
    r.Use(middleware.Cors())
    r.Use(middleware.RateLimit())

    // 健康检查
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // API v1
    v1 := r.Group("/api/v1")
    {
        // 公开接口
        public := v1.Group("")
        {
            public.POST("/login", controller.Login)
            public.POST("/register", controller.Register)
        }

        // 需要认证的接口
        auth := v1.Group("")
        auth.Use(middleware.Auth())
        {
            // 用户管理
            users := auth.Group("/users")
            {
                users.POST("", controller.CreateUser)
                users.GET("/:id", controller.GetUser)
                users.PUT("/:id", controller.UpdateUser)
                users.DELETE("/:id", controller.DeleteUser)
                users.GET("", controller.ListUsers)
            }

            // 教练管理
            coaches := auth.Group("/coaches")
            {
                coaches.POST("", controller.CreateCoach)
                coaches.GET("/:id", controller.GetCoach)
                coaches.PUT("/:id", controller.UpdateCoach)
                coaches.DELETE("/:id", controller.DeleteCoach)
                coaches.GET("", controller.ListCoaches)
            }

            // 会员卡管理
            cards := auth.Group("/membership-cards")
            {
                cards.POST("", controller.CreateMembershipCard)
                cards.GET("/:id", controller.GetMembershipCard)
                cards.PUT("/:id", controller.UpdateMembershipCard)
                cards.DELETE("/:id", controller.DeleteMembershipCard)
                cards.GET("", controller.ListMembershipCards)
            }

            // 课程预约
            bookings := auth.Group("/bookings")
            {
                bookings.POST("", controller.CreateBooking)
                bookings.GET("/:id", controller.GetBooking)
                bookings.PUT("/:id/cancel", controller.CancelBooking)
                bookings.PUT("/:id/confirm", controller.ConfirmBooking)
                bookings.GET("", controller.ListBookings)
            }
        }
    }

    return r
}
```

**验收标准**:
- 服务可正常启动
- 配置加载正确
- 路由注册成功
- 中间件生效
- 支持优雅关闭

---

### 2.2 数据库设计与初始化

#### 2.2.1 数据库连接池配置
**输入**: 数据库配置
**输出**: 数据库连接池

**执行内容**:

**A. 数据库封装 (pkg/database/mysql.go)**
```go
package database

import (
    "fmt"
    "time"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

type Config struct {
    Host            string
    Port            int
    Username        string
    Password        string
    Database        string
    MaxIdleConns    int
    MaxOpenConns    int
    ConnMaxLifetime int
}

func Init(cfg Config) error {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.Username,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.Database,
    )

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        return err
    }

    sqlDB, err := db.DB()
    if err != nil {
        return err
    }

    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

    DB = db
    return nil
}

func Close() error {
    sqlDB, err := DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

func GetDB() *gorm.DB {
    return DB
}
```

**B. 数据库迁移 (scripts/init_db.sql)**
```sql
-- 创建数据库
CREATE DATABASE IF NOT EXISTS gym_admin DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE gym_admin;

-- 用户表
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_no VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    gender TINYINT,
    birthday DATE,
    id_card VARCHAR(18),
    phone VARCHAR(11) UNIQUE NOT NULL,
    email VARCHAR(100),
    avatar_url VARCHAR(255),
    address VARCHAR(255),
    emergency_contact VARCHAR(50),
    emergency_phone VARCHAR(11),
    health_status TEXT,
    training_goal TEXT,
    source TINYINT DEFAULT 1,
    status TINYINT DEFAULT 1,
    remark TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_phone (phone),
    INDEX idx_user_no (user_no),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 其他表结构参考TASK001-TASK006的数据库设计
```

**验收标准**:
- 数据库连接成功
- 连接池配置生效
- 支持自动重连
- 查询性能良好

---

### 2.3 缓存系统搭建

#### 2.3.1 Redis配置与封装
**输入**: Redis配置
**输出**: Redis客户端封装

**执行内容**:

**A. Redis封装 (pkg/cache/redis.go)**
```go
package cache

import (
    "context"
    "fmt"
    "time"

    "github.com/go-redis/redis/v8"
)

var Client *redis.Client
var ctx = context.Background()

type Config struct {
    Host     string
    Port     int
    Password string
    DB       int
}

func Init(cfg Config) error {
    Client = redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
        Password: cfg.Password,
        DB:       cfg.DB,
    })

    _, err := Client.Ping(ctx).Result()
    return err
}

func Close() error {
    return Client.Close()
}

func Set(key string, value interface{}, expiration time.Duration) error {
    return Client.Set(ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
    return Client.Get(ctx, key).Result()
}

func Del(keys ...string) error {
    return Client.Del(ctx, keys...).Err()
}

func Exists(keys ...string) (int64, error) {
    return Client.Exists(ctx, keys...).Result()
}

func Expire(key string, expiration time.Duration) error {
    return Client.Expire(ctx, key, expiration).Err()
}

// 缓存常用场景封装
func GetOrSet(key string, expiration time.Duration, fn func() (interface{}, error)) (interface{}, error) {
    // 先从缓存获取
    val, err := Get(key)
    if err == nil {
        return val, nil
    }

    // 缓存未命中，执行函数获取数据
    data, err := fn()
    if err != nil {
        return nil, err
    }

    // 写入缓存
    if err := Set(key, data, expiration); err != nil {
        return data, err
    }

    return data, nil
}
```

**B. 缓存策略**
```go
// 缓存键命名规范
const (
    UserCacheKey          = "user:%d"                    // 用户信息
    CoachCacheKey         = "coach:%d"                   // 教练信息
    MembershipCardKey     = "membership_card:%d"         // 会员卡信息
    BookingCacheKey       = "booking:%d"                 // 预约信息
    TokenCacheKey         = "token:%s"                   // 用户Token
    VerifyCodeKey         = "verify_code:%s"             // 验证码
)

// 缓存过期时间
const (
    UserCacheExpire          = 1 * time.Hour
    CoachCacheExpire         = 1 * time.Hour
    MembershipCardExpire     = 30 * time.Minute
    BookingCacheExpire       = 30 * time.Minute
    TokenExpire              = 24 * time.Hour
    VerifyCodeExpire         = 5 * time.Minute
)
```

**验收标准**:
- Redis连接成功
- 缓存读写正常
- 过期时间生效
- 支持常用操作

---

### 2.4 文件存储服务

#### 2.4.1 OSS对象存储配置
**输入**: OSS配置
**输出**: 文件上传下载服务

**执行内容**:

**A. OSS封装 (pkg/oss/oss.go)**
```go
package oss

import (
    "fmt"
    "io"
    "path"
    "time"

    "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var client *oss.Client
var bucket *oss.Bucket

type Config struct {
    Endpoint        string
    AccessKeyID     string
    AccessKeySecret string
    BucketName      string
}

func Init(cfg Config) error {
    var err error
    client, err = oss.New(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret)
    if err != nil {
        return err
    }

    bucket, err = client.Bucket(cfg.BucketName)
    return err
}

// UploadFile 上传文件
func UploadFile(objectKey string, reader io.Reader) (string, error) {
    err := bucket.PutObject(objectKey, reader)
    if err != nil {
        return "", err
    }

    // 返回文件URL
    url := fmt.Sprintf("https://%s.%s/%s", bucket.BucketName, client.Config.Endpoint, objectKey)
    return url, nil
}

// UploadFileWithPath 上传文件并自动生成路径
func UploadFileWithPath(filename string, reader io.Reader) (string, error) {
    // 生成路径：年/月/日/时间戳_文件名
    now := time.Now()
    objectKey := fmt.Sprintf("%d/%02d/%02d/%d_%s",
        now.Year(), now.Month(), now.Day(), now.Unix(), filename)

    return UploadFile(objectKey, reader)
}

// DeleteFile 删除文件
func DeleteFile(objectKey string) error {
    return bucket.DeleteObject(objectKey)
}

// GetSignedURL 获取签名URL（用于私有文件访问）
func GetSignedURL(objectKey string, expireTime time.Duration) (string, error) {
    return bucket.SignURL(objectKey, oss.HTTPGet, int64(expireTime.Seconds()))
}

// IsFileExist 检查文件是否存在
func IsFileExist(objectKey string) (bool, error) {
    return bucket.IsObjectExist(objectKey)
}
```

**B. 文件上传接口 (internal/controller/upload_controller.go)**
```go
package controller

import (
    "gym-admin/pkg/oss"
    "gym-admin/pkg/response"
    "path/filepath"

    "github.com/gin-gonic/gin"
)

// UploadImage 上传图片
func UploadImage(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        response.Error(c, "文件上传失败")
        return
    }

    // 验证文件类型
    ext := filepath.Ext(file.Filename)
    if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
        response.Error(c, "只支持jpg、jpeg、png格式")
        return
    }

    // 验证文件大小（最大5MB）
    if file.Size > 5*1024*1024 {
        response.Error(c, "文件大小不能超过5MB")
        return
    }

    // 打开文件
    src, err := file.Open()
    if err != nil {
        response.Error(c, "文件打开失败")
        return
    }
    defer src.Close()

    // 上传到OSS
    url, err := oss.UploadFileWithPath(file.Filename, src)
    if err != nil {
        response.Error(c, "文件上传失败")
        return
    }

    response.Success(c, gin.H{"url": url})
}
```

**验收标准**:
- 文件上传成功
- 支持多种文件类型
- 文件大小限制生效
- 返回可访问的URL

---

### 2.5 中间件开发

#### 2.5.1 认证中间件
**输入**: JWT配置
**输出**: 认证中间件

**执行内容**:

**A. JWT工具 (pkg/utils/jwt.go)**
```go
package utils

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v4"
)

type Claims struct {
    UserID int64  `json:"user_id"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

var jwtSecret []byte

func InitJWT(secret string) {
    jwtSecret = []byte(secret)
}

// GenerateToken 生成Token
func GenerateToken(userID int64, role string, expireTime int) (string, error) {
    claims := Claims{
        UserID: userID,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireTime) * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}
```

**B. 认证中间件 (internal/middleware/auth.go)**
```go
package middleware

import (
    "gym-admin/pkg/response"
    "gym-admin/pkg/utils"
    "strings"

    "github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.Unauthorized(c, "未登录")
            c.Abort()
            return
        }

        // 验证Bearer token格式
        parts := strings.SplitN(authHeader, " ", 2)
        if !(len(parts) == 2 && parts[0] == "Bearer") {
            response.Unauthorized(c, "Token格式错误")
            c.Abort()
            return
        }

        // 解析Token
        claims, err := utils.ParseToken(parts[1])
        if err != nil {
            response.Unauthorized(c, "Token无效")
            c.Abort()
            return
        }

        // 将用户信息存入上下文
        c.Set("user_id", claims.UserID)
        c.Set("role", claims.Role)

        c.Next()
    }
}

// RequireRole 角色权限中间件
func RequireRole(roles ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        role, exists := c.Get("role")
        if !exists {
            response.Forbidden(c, "无权限")
            c.Abort()
            return
        }

        userRole := role.(string)
        for _, r := range roles {
            if userRole == r {
                c.Next()
                return
            }
        }

        response.Forbidden(c, "无权限")
        c.Abort()
    }
}
```

**C. 其他中间件**
```go
// CORS中间件 (internal/middleware/cors.go)
func Cors() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

// 限流中间件 (internal/middleware/rate_limit.go)
func RateLimit() gin.HandlerFunc {
    // 使用令牌桶算法实现限流
    return func(c *gin.Context) {
        // 实现限流逻辑
        c.Next()
    }
}

// 日志中间件 (internal/middleware/logger.go)
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery

        c.Next()

        latency := time.Since(start)
        clientIP := c.ClientIP()
        method := c.Request.Method
        statusCode := c.Writer.Status()

        if raw != "" {
            path = path + "?" + raw
        }

        logger.Info(fmt.Sprintf("[GIN] %s | %3d | %13v | %15s | %-7s %s",
            time.Now().Format("2006/01/02 - 15:04:05"),
            statusCode,
            latency,
            clientIP,
            method,
            path,
        ))
    }
}
```

**验收标准**:
- Token生成和解析正确
- 认证中间件生效
- 权限控制准确
- CORS配置正确
- 限流功能正常

---

### 2.6 日志系统

#### 2.6.1 日志配置
**输入**: 日志配置
**输出**: 日志系统

**执行内容**:

**A. 日志封装 (pkg/logger/logger.go)**
```go
package logger

import (
    "os"

    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.Logger

type Config struct {
    Level      string
    Filename   string
    MaxSize    int
    MaxBackups int
    MaxAge     int
    Compress   bool
}

func Init(cfg Config) error {
    // 日志级别
    var level zapcore.Level
    switch cfg.Level {
    case "debug":
        level = zapcore.DebugLevel
    case "info":
        level = zapcore.InfoLevel
    case "warn":
        level = zapcore.WarnLevel
    case "error":
        level = zapcore.ErrorLevel
    default:
        level = zapcore.InfoLevel
    }

    // 编码器配置
    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "time",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "caller",
        MessageKey:     "msg",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,
        EncodeTime:     zapcore.ISO8601TimeEncoder,
        EncodeDuration: zapcore.SecondsDurationEncoder,
        EncodeCaller:   zapcore.ShortCallerEncoder,
    }

    // 日志轮转
    writer := &lumberjack.Logger{
        Filename:   cfg.Filename,
        MaxSize:    cfg.MaxSize,
        MaxBackups: cfg.MaxBackups,
        MaxAge:     cfg.MaxAge,
        Compress:   cfg.Compress,
    }

    // 同时输出到文件和控制台
    core := zapcore.NewTee(
        zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(writer), level),
        zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level),
    )

    log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
    return nil
}

func Debug(msg string, fields ...zap.Field) {
    log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
    log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
    log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
    log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
    log.Fatal(msg, fields...)
}

func Sync() {
    log.Sync()
}
```

**验收标准**:
- 日志正常输出
- 日志轮转生效
- 日志级别控制正确
- 支持结构化日志

---

### 2.7 Docker容器化

#### 2.7.1 Dockerfile编写
**输入**: 应用程序
**输出**: Docker镜像

**执行内容**:

**A. Dockerfile**
```dockerfile
# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制go mod文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
COPY --from=builder /app/internal/config/config.yaml ./internal/config/

# 设置时区
ENV TZ=Asia/Shanghai

EXPOSE 8080

CMD ["./main"]
```

**B. docker-compose.yml**
```yaml
version: '3.8'

services:
  # 后端服务
  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
    depends_on:
      - mysql
      - redis
    volumes:
      - ./logs:/root/logs
    networks:
      - gym-network

  # MySQL数据库
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root123
      MYSQL_DATABASE: gym_admin
    ports:
      - "3306:3306"
    volumes:
      - mysql-data:/var/lib/mysql
      - ./scripts/init_db.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - gym-network

  # Redis缓存
  redis:
    image: redis:6.0-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    networks:
      - gym-network

  # Nginx反向代理
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
    depends_on:
      - backend
    networks:
      - gym-network

volumes:
  mysql-data:
  redis-data:

networks:
  gym-network:
    driver: bridge
```

**C. Makefile**
```makefile
.PHONY: build run test docker-build docker-up docker-down

# 编译
build:
	go build -o bin/server cmd/server/main.go

# 运行
run:
	go run cmd/server/main.go

# 测试
test:
	go test -v ./...

# 构建Docker镜像
docker-build:
	docker build -t gym-admin-backend:latest .

# 启动Docker Compose
docker-up:
	docker-compose up -d

# 停止Docker Compose
docker-down:
	docker-compose down

# 查看日志
logs:
	docker-compose logs -f backend

# 数据库迁移
migrate:
	mysql -h localhost -u root -p < scripts/init_db.sql
```

**验收标准**:
- Docker镜像构建成功
- 容器正常启动
- 服务间网络通信正常
- 数据持久化正常

---

## 三、部署文档

### 3.1 开发环境部署

```bash
# 1. 克隆代码
git clone https://github.com/your-repo/gym-admin-backend.git
cd gym-admin-backend

# 2. 安装依赖
go mod download

# 3. 配置文件
cp internal/config/config.example.yaml internal/config/config.yaml
# 编辑config.yaml，配置数据库、Redis等信息

# 4. 初始化数据库
mysql -u root -p < scripts/init_db.sql

# 5. 运行服务
go run cmd/server/main.go
```

### 3.2 生产环境部署

```bash
# 1. 使用Docker Compose部署
docker-compose up -d

# 2. 查看服务状态
docker-compose ps

# 3. 查看日志
docker-compose logs -f

# 4. 停止服务
docker-compose down
```

### 3.3 监控配置

**Prometheus配置 (prometheus.yml)**
```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'gym-admin'
    static_configs:
      - targets: ['backend:8080']
```

**Grafana Dashboard配置**
- 导入预设的Dashboard模板
- 配置数据源为Prometheus
- 设置告警规则

## 四、测试用例

### 4.1 单元测试
- 配置加载测试
- 数据库连接测试
- Redis连接测试
- JWT生成和解析测试
- 文件上传测试

### 4.2 集成测试
- API接口测试
- 中间件测试
- 数据库事务测试
- 缓存一致性测试

### 4.3 性能测试
- 并发压力测试
- 数据库查询性能测试
- 缓存命中率测试
- 接口响应时间测试

## 五、上线检查清单

### 5.1 基础设施检查
- [ ] 数据库配置正确
- [ ] Redis配置正确
- [ ] OSS配置正确
- [ ] 日志系统正常
- [ ] 监控系统正常

### 5.2 安全检查
- [ ] JWT密钥配置
- [ ] 数据库密码强度
- [ ] API接口鉴权
- [ ] HTTPS配置
- [ ] 防火墙规则

### 5.3 性能检查
- [ ] 数据库连接池配置
- [ ] Redis连接池配置
- [ ] 接口限流配置
- [ ] 缓存策略配置
- [ ] 静态资源CDN

### 5.4 运维检查
- [ ] 日志轮转配置
- [ ] 备份策略配置
- [ ] 监控告警配置
- [ ] 部署文档完善
- [ ] 回滚方案准备

## 六、后续优化方向

1. **微服务拆分**: 将单体应用拆分为多个微服务
2. **服务网格**: 引入Istio等服务网格
3. **消息队列**: 引入RabbitMQ/Kafka处理异步任务
4. **分布式追踪**: 引入Jaeger进行链路追踪
5. **自动化运维**: 完善CI/CD流程，实现自动化部署
6. **高可用**: 实现数据库主从、Redis集群等高可用方案

---
**任务优先级**: P0（基础设施）  
**预计工期**: 1-2周  
**依赖任务**: 无  
**后续任务**: 所有业务功能任务