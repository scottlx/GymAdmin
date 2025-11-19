# 健身房管理系统 (GymAdmin)

基于 Golang + React 的健身房管理系统，支持会员管理、会员卡管理、教练管理、课程预约、人脸识别签到、第三方券核销等功能。

## 技术栈

### 后端
- **语言**: Golang 1.21+
- **框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0+
- **缓存**: Redis 6.0+
- **配置**: Viper
- **日志**: Zap + Lumberjack
- **认证**: JWT
- **容器化**: Docker

### 前端
- **框架**: React 18+ with TypeScript
- **路由**: React Router v6
- **UI组件**: Ant Design 5
- **状态管理**: Zustand
- **构建工具**: Vite
- **HTTP客户端**: Axios

## 项目结构

```
GymAdmin/
├── backend/                    # 后端代码
│   ├── cmd/
│   │   └── server/
│   │       └── main.go        # 程序入口
│   ├── internal/
│   │   ├── config/            # 配置管理
│   │   ├── models/            # 数据模型
│   │   ├── repository/        # 数据访问层
│   │   ├── service/           # 业务逻辑层
│   │   ├── controller/        # 控制器层
│   │   ├── middleware/        # 中间件
│   │   └── router/            # 路由配置
│   ├── pkg/
│   │   ├── database/          # 数据库封装
│   │   ├── cache/             # 缓存封装
│   │   ├── logger/            # 日志封装
│   │   ├── jwt/               # JWT工具
│   │   └── response/          # 响应封装
│   ├── go.mod
│   └── Dockerfile
├── frontend/                   # 前端代码
│   ├── src/
│   │   ├── layouts/           # 布局组件
│   │   ├── pages/             # 页面组件
│   │   ├── services/          # API服务层
│   │   ├── store/             # 状态管理
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── package.json
│   └── vite.config.ts
├── docker-compose.yml          # Docker编排
└── README.md
```

## 快速开始

### 使用 Docker Compose（推荐）

1. 克隆项目
```bash
git clone <repository-url>
cd GymAdmin
```

2. 启动所有服务
```bash
docker-compose up -d
```

3. 访问应用
- 前端: http://localhost:3000
- 后端API: http://localhost:8080
- MySQL: localhost:3306
- Redis: localhost:6379

### 本地开发

#### 后端开发

1. 安装依赖
```bash
cd backend
go mod download
```

2. 配置数据库
编辑 `backend/internal/config/config.yaml`，修改数据库连接信息

3. 运行
```bash
go run cmd/server/main.go
```

#### 前端开发

1. 安装依赖
```bash
cd frontend
npm install
```

2. 运行开发服务器
```bash
npm run dev
```

3. 构建生产版本
```bash
npm run build
```

## 已实现功能

### ✅ 基础架构
- 项目基础架构搭建
- 数据库模型定义（10个核心表）
- RESTful API路由配置
- 中间件（CORS、日志、JWT认证）
- 统一响应格式
- Docker容器化部署

### ✅ JWT认证
- Token生成与验证
- 登录/注册接口
- 请求拦截器自动添加Token
- Token过期自动跳转登录

### ✅ 用户管理
- 用户列表查询（分页）
- 用户创建
- 用户详情查看
- 用户信息更新
- 用户删除
- 用户训练统计

### ✅ 教练管理
- 教练列表查询（分页）
- 教练创建
- 教练详情查看
- 教练信息更新
- 教练删除

### ✅ 前端功能
- 登录页面（集成API）
- 用户管理页面（完整CRUD）
- 教练管理页面（完整CRUD）
- 仪表盘页面
- 响应式布局
- 状态管理（Zustand）

### ⏳ 待实现功能
- 会员卡管理（开卡、续费、冻结、转卡）
- 课程管理与预约
- 人脸识别签到
- 美团/抖音券核销
- 数据统计与报表
- 微信小程序

## API文档

API基础路径: `/api/v1`

### 认证相关
- `POST /login` - 用户登录
- `POST /register` - 用户注册

### 会员管理
- `GET /users` - 获取会员列表（支持分页、状态筛选）
- `POST /users` - 创建会员
- `GET /users/:id` - 获取会员详情
- `PUT /users/:id` - 更新会员信息
- `DELETE /users/:id` - 删除会员
- `GET /users/:id/stats` - 获取会员训练统计

### 教练管理
- `GET /coaches` - 获取教练列表（支持分页、状态筛选）
- `POST /coaches` - 添加教练
- `GET /coaches/:id` - 获取教练详情
- `PUT /coaches/:id` - 更新教练信息
- `DELETE /coaches/:id` - 删除教练

### 会员卡管理（待实现）
- `GET /cards` - 获取会员卡列表
- `POST /cards` - 办理会员卡
- `GET /cards/:id` - 获取会员卡详情
- `PUT /cards/:id` - 更新会员卡

### 课程管理（待实现）
- `GET /courses` - 获取课程列表
- `POST /courses` - 创建课程
- `GET /courses/:id` - 获取课程详情
- `PUT /courses/:id` - 更新课程
- `DELETE /courses/:id` - 删除课程

### 预约管理（待实现）
- `GET /bookings` - 获取预约列表
- `POST /bookings` - 创建预约
- `DELETE /bookings/:id` - 取消预约

### 签到管理（待实现）
- `GET /checkins` - 获取签到记录
- `POST /checkins` - 创建签到记录

### 券核销（待实现）
- `GET /vouchers` - 获取券列表
- `POST /vouchers/verify` - 核销券

## 数据库设计

主要数据表：
- `users` - 用户表
- `user_training_stats` - 用户训练统计
- `card_types` - 会员卡类型
- `membership_cards` - 会员卡
- `coaches` - 教练
- `courses` - 课程
- `bookings` - 预约记录
- `check_ins` - 签到记录
- `face_records` - 人脸信息
- `voucher_records` - 券核销记录

## 开发规范

### Git提交规范
- `feat`: 新功能
- `fix`: 修复bug
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建/工具相关

### 代码规范
- 后端遵循 Effective Go 规范
- 前端遵循 Airbnb React Style Guide
- 使用有意义的变量和函数命名
- 添加必要的注释

## 项目特点

1. **分层架构**: 严格遵循MVC模式，Repository-Service-Controller三层架构
2. **RESTful API**: 标准的REST接口设计
3. **统一响应**: 所有API返回统一的JSON格式
4. **JWT认证**: 完整的JWT token生成和验证机制
5. **自动迁移**: 使用GORM的AutoMigrate功能自动创建表
6. **容器化**: 完整的Docker Compose配置，一键启动
7. **TypeScript**: 前端使用TypeScript提供类型安全
8. **状态管理**: 使用Zustand进行轻量级状态管理
9. **API封装**: 完整的API服务层封装，统一错误处理

## 许可证

MIT License

## 联系方式

如有问题，请提交 Issue 或 Pull Request。
