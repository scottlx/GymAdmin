# TASK001 - 用户管理功能

## 一、功能概述

### 1.1 功能目标
实现健身房会员的完整生命周期管理，包括用户注册、信息维护、训练数据统计、状态管理等核心功能。

### 1.2 核心价值
- 建立完整的会员档案体系
- 追踪会员训练数据，提供数据化运营支持
- 为后续会员卡、课程等功能提供用户基础

### 1.3 涉及角色
- **管理员**: 可以查看、创建、编辑、删除所有用户信息
- **前台人员**: 可以查看、创建、编辑用户信息（不能删除）
- **教练**: 可以查看自己学员的信息
- **会员**: 可以查看和编辑自己的基本信息

## 二、功能详细拆解

### 2.1 用户基本信息管理

#### 2.1.1 数据库设计
**输入**: 业务需求分析
**输出**: 数据库表结构SQL文件

**执行内容**:
```sql
-- 用户表 (users)
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    user_no VARCHAR(32) UNIQUE NOT NULL COMMENT '用户编号（自动生成）',
    name VARCHAR(50) NOT NULL COMMENT '姓名',
    gender TINYINT COMMENT '性别：1-男，2-女',
    birthday DATE COMMENT '生日',
    id_card VARCHAR(18) COMMENT '身份证号（加密存储）',
    phone VARCHAR(11) UNIQUE NOT NULL COMMENT '手机号',
    email VARCHAR(100) COMMENT '邮箱',
    avatar_url VARCHAR(255) COMMENT '头像URL',
    address VARCHAR(255) COMMENT '地址',
    emergency_contact VARCHAR(50) COMMENT '紧急联系人',
    emergency_phone VARCHAR(11) COMMENT '紧急联系电话',
    health_status TEXT COMMENT '健康状况说明',
    training_goal TEXT COMMENT '训练目标',
    source TINYINT DEFAULT 1 COMMENT '来源：1-前台录入，2-小程序注册，3-美团，4-抖音',
    status TINYINT DEFAULT 1 COMMENT '状态：1-正常，2-冻结，3-黑名单',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    INDEX idx_phone (phone),
    INDEX idx_user_no (user_no),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 用户训练统计表 (user_training_stats)
CREATE TABLE user_training_stats (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    total_days INT DEFAULT 0 COMMENT '累计训练天数',
    total_times INT DEFAULT 0 COMMENT '累计训练次数',
    continuous_days INT DEFAULT 0 COMMENT '连续训练天数',
    last_check_in_date DATE COMMENT '最后签到日期',
    month_times INT DEFAULT 0 COMMENT '本月训练次数',
    year_times INT DEFAULT 0 COMMENT '本年训练次数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户训练统计表';
```

**验收标准**:
- 表结构符合第三范式
- 包含必要的索引优化查询性能
- 敏感字段（身份证）需要加密存储
- 支持软删除机制

---

#### 2.1.2 后端API开发 - 用户CRUD

**输入**: 数据库表结构
**输出**: 完整的用户管理API接口

**执行内容**:

**A. 定义数据模型 (models/user.go)**
```go
type User struct {
    ID               int64      `json:"id" gorm:"primaryKey"`
    UserNo           string     `json:"user_no" gorm:"uniqueIndex;size:32"`
    Name             string     `json:"name" gorm:"size:50;not null"`
    Gender           int8       `json:"gender"`
    Birthday         *time.Time `json:"birthday"`
    IDCard           string     `json:"id_card" gorm:"size:18"` // 加密存储
    Phone            string     `json:"phone" gorm:"uniqueIndex;size:11;not null"`
    Email            string     `json:"email" gorm:"size:100"`
    AvatarURL        string     `json:"avatar_url" gorm:"size:255"`
    Address          string     `json:"address" gorm:"size:255"`
    EmergencyContact string     `json:"emergency_contact" gorm:"size:50"`
    EmergencyPhone   string     `json:"emergency_phone" gorm:"size:11"`
    HealthStatus     string     `json:"health_status" gorm:"type:text"`
    TrainingGoal     string     `json:"training_goal" gorm:"type:text"`
    Source           int8       `json:"source" gorm:"default:1"`
    Status           int8       `json:"status" gorm:"default:1"`
    Remark           string     `json:"remark" gorm:"type:text"`
    CreatedAt        time.Time  `json:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at"`
    DeletedAt        *time.Time `json:"deleted_at" gorm:"index"`
}

type UserTrainingStats struct {
    ID               int64     `json:"id" gorm:"primaryKey"`
    UserID           int64     `json:"user_id" gorm:"uniqueIndex;not null"`
    TotalDays        int       `json:"total_days" gorm:"default:0"`
    TotalTimes       int       `json:"total_times" gorm:"default:0"`
    ContinuousDays   int       `json:"continuous_days" gorm:"default:0"`
    LastCheckInDate  *time.Time `json:"last_check_in_date"`
    MonthTimes       int       `json:"month_times" gorm:"default:0"`
    YearTimes        int       `json:"year_times" gorm:"default:0"`
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}
```

**B. 实现业务逻辑层 (services/user_service.go)**
```go
// 核心方法列表
- CreateUser(req *CreateUserRequest) (*User, error)
- GetUserByID(id int64) (*User, error)
- GetUserByPhone(phone string) (*User, error)
- UpdateUser(id int64, req *UpdateUserRequest) error
- DeleteUser(id int64) error // 软删除
- ListUsers(req *ListUsersRequest) ([]*User, int64, error)
- GenerateUserNo() string // 生成用户编号
- EncryptIDCard(idCard string) string // 身份证加密
- DecryptIDCard(encrypted string) string // 身份证解密
```

**C. 实现API接口层 (controllers/user_controller.go)**
```go
// API路由定义
POST   /api/v1/users              // 创建用户
GET    /api/v1/users/:id          // 获取用户详情
PUT    /api/v1/users/:id          // 更新用户信息
DELETE /api/v1/users/:id          // 删除用户
GET    /api/v1/users              // 用户列表（支持分页、搜索、筛选）
GET    /api/v1/users/phone/:phone // 根据手机号查询用户
```

**D. 请求/响应结构定义**
```go
// 创建用户请求
type CreateUserRequest struct {
    Name             string `json:"name" binding:"required,max=50"`
    Gender           int8   `json:"gender" binding:"omitempty,oneof=1 2"`
    Birthday         string `json:"birthday" binding:"omitempty"`
    IDCard           string `json:"id_card" binding:"omitempty,len=18"`
    Phone            string `json:"phone" binding:"required,len=11"`
    Email            string `json:"email" binding:"omitempty,email"`
    Address          string `json:"address"`
    EmergencyContact string `json:"emergency_contact"`
    EmergencyPhone   string `json:"emergency_phone" binding:"omitempty,len=11"`
    HealthStatus     string `json:"health_status"`
    TrainingGoal     string `json:"training_goal"`
    Source           int8   `json:"source" binding:"omitempty,oneof=1 2 3 4"`
    Remark           string `json:"remark"`
}

// 用户列表请求
type ListUsersRequest struct {
    Page     int    `form:"page" binding:"required,min=1"`
    PageSize int    `form:"page_size" binding:"required,min=1,max=100"`
    Keyword  string `form:"keyword"`  // 搜索关键词（姓名、手机号）
    Gender   int8   `form:"gender"`   // 性别筛选
    Status   int8   `form:"status"`   // 状态筛选
    Source   int8   `form:"source"`   // 来源筛选
    StartDate string `form:"start_date"` // 注册开始日期
    EndDate   string `form:"end_date"`   // 注册结束日期
}
```

**验收标准**:
- 所有API接口通过Postman测试
- 参数校验完整，错误提示清晰
- 身份证号码加密存储和解密正确
- 用户编号自动生成且唯一
- 支持软删除，删除后数据可恢复
- 列表接口支持分页、搜索、多条件筛选
- 接口响应时间 < 200ms（单条查询）

---

#### 2.1.3 前端页面开发 - 用户列表

**输入**: 后端API接口文档
**输出**: 用户列表页面

**执行内容**:

**A. 页面布局设计**
- 顶部搜索栏：关键词搜索、高级筛选（性别、状态、来源、日期范围）
- 操作按钮区：新增用户、批量导入、导出Excel
- 数据表格：展示用户列表，支持排序
- 分页组件：页码切换、每页条数选择

**B. 核心功能实现**
```typescript
// 组件文件: src/pages/User/UserList.tsx

// 功能点：
1. 用户列表展示
   - 表格列：用户编号、姓名、性别、手机号、会员状态、训练天数、训练次数、注册时间、操作
   - 支持按字段排序
   - 支持行选择（批量操作）

2. 搜索与筛选
   - 关键词搜索（姓名、手机号）
   - 性别筛选下拉框
   - 状态筛选下拉框
   - 来源筛选下拉框
   - 日期范围选择器

3. 操作按钮
   - 查看详情：跳转到用户详情页
   - 编辑：打开编辑弹窗
   - 删除：二次确认后删除
   - 更多：下拉菜单（冻结/解冻、加入黑名单等）

4. 状态标识
   - 正常：绿色标签
   - 冻结：橙色标签
   - 黑名单：红色标签
```

**C. 状态管理**
```typescript
// 使用Redux Toolkit或Zustand管理状态
interface UserListState {
  users: User[];
  total: number;
  loading: boolean;
  filters: FilterParams;
  selectedRowKeys: number[];
}
```

**验收标准**:
- 页面响应式设计，适配不同屏幕尺寸
- 表格数据加载流畅，支持虚拟滚动（数据量大时）
- 搜索和筛选实时生效
- 操作按钮权限控制正确
- 用户体验流畅，无明显卡顿

---

#### 2.1.4 前端页面开发 - 用户详情与编辑

**输入**: 后端API接口文档
**输出**: 用户详情页和编辑表单

**执行内容**:

**A. 用户详情页 (src/pages/User/UserDetail.tsx)**
```typescript
// 页面布局：
1. 顶部信息卡片
   - 头像、姓名、用户编号
   - 会员状态标签
   - 快捷操作按钮（编辑、冻结、删除）

2. 基本信息Tab
   - 个人信息：姓名、性别、生日、身份证、手机号、邮箱、地址
   - 紧急联系人信息
   - 健康状况、训练目标
   - 来源、注册时间

3. 训练数据Tab
   - 累计训练天数、累计训练次数
   - 连续训练天数
   - 本月训练次数、本年训练次数
   - 最后签到时间
   - 训练趋势图表（ECharts）

4. 会员卡Tab
   - 当前会员卡列表
   - 历史会员卡记录

5. 课程记录Tab
   - 已预约课程
   - 已完成课程
   - 课程统计

6. 操作日志Tab
   - 用户相关的所有操作记录
```

**B. 用户编辑表单 (src/components/User/UserForm.tsx)**
```typescript
// 表单字段：
1. 基本信息
   - 姓名（必填）
   - 性别（单选）
   - 生日（日期选择器）
   - 身份证号（输入框，自动校验格式）
   - 手机号（必填，自动校验格式）
   - 邮箱（自动校验格式）
   - 头像上传

2. 详细信息
   - 地址（输入框）
   - 紧急联系人（输入框）
   - 紧急联系电话（输入框）
   - 健康状况（文本域）
   - 训练目标（文本域）

3. 其他信息
   - 来源（下拉选择）
   - 状态（下拉选择）
   - 备注（文本域）

// 表单校验规则：
- 姓名：必填，最大50字符
- 手机号：必填，11位数字，格式校验
- 身份证：18位，格式校验，自动解析生日和性别
- 邮箱：邮箱格式校验
- 紧急联系电话：11位数字
```

**C. 头像上传功能**
```typescript
// 功能实现：
1. 支持点击上传和拖拽上传
2. 图片格式限制：jpg、png、jpeg
3. 图片大小限制：最大2MB
4. 上传前预览
5. 上传到OSS，返回URL
6. 支持裁剪（可选）
```

**验收标准**:
- 详情页信息展示完整准确
- 表单校验规则完整，错误提示友好
- 头像上传功能正常，支持预览
- 身份证号输入后自动解析生日和性别
- 编辑保存后数据更新成功
- 页面加载性能良好

---

### 2.2 用户训练数据统计

#### 2.2.1 签到功能实现

**输入**: 用户基本信息
**输出**: 签到记录和训练统计更新

**执行内容**:

**A. 数据库设计**
```sql
-- 签到记录表 (check_in_records)
CREATE TABLE check_in_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    check_in_type TINYINT NOT NULL COMMENT '签到类型：1-人脸识别，2-手动签到，3-扫码签到',
    check_in_time TIMESTAMP NOT NULL COMMENT '签到时间',
    check_out_time TIMESTAMP NULL COMMENT '签出时间',
    duration INT COMMENT '训练时长（分钟）',
    device_id VARCHAR(50) COMMENT '设备ID（人脸识别设备）',
    operator_id BIGINT COMMENT '操作员ID（手动签到时）',
    remark VARCHAR(255) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_check_in_time (check_in_time),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='签到记录表';
```

**B. 后端API实现**
```go
// services/checkin_service.go

// 核心方法：
- CheckIn(userID int64, checkInType int8, deviceID string) error
- CheckOut(userID int64) error
- GetTodayCheckInRecord(userID int64) (*CheckInRecord, error)
- UpdateTrainingStats(userID int64) error // 更新训练统计
- CalculateContinuousDays(userID int64) int // 计算连续训练天数

// API路由：
POST   /api/v1/check-in              // 签到
POST   /api/v1/check-out             // 签出
GET    /api/v1/check-in/today/:user_id // 获取今日签到记录
GET    /api/v1/check-in/records      // 签到记录列表
```

**C. 业务逻辑**
```
1. 签到流程：
   - 验证用户是否存在
   - 验证用户状态是否正常
   - 验证是否有有效会员卡
   - 检查今日是否已签到
   - 创建签到记录
   - 更新训练统计（累计天数+1，累计次数+1，本月次数+1，本年次数+1）
   - 计算连续训练天数
   - 返回签到成功

2. 签出流程：
   - 查找今日签到记录
   - 更新签出时间
   - 计算训练时长
   - 返回签出成功

3. 连续训练天数计算逻辑：
   - 查询用户最近的签到记录
   - 从今天往前推，连续有签到记录的天数
   - 如果中断则重置为1
```

**验收标准**:
- 签到接口响应时间 < 500ms
- 同一用户同一天只能签到一次
- 训练统计数据准确无误
- 连续训练天数计算正确
- 支持并发签到（多个设备同时签到）

---

#### 2.2.2 训练数据统计与展示

**输入**: 签到记录数据
**输出**: 训练数据统计报表

**执行内容**:

**A. 后端统计API**
```go
// services/stats_service.go

// 核心方法：
- GetUserTrainingStats(userID int64) (*UserTrainingStats, error)
- GetUserTrainingTrend(userID int64, days int) ([]TrendData, error) // 训练趋势
- GetMonthlyTrainingCalendar(userID int64, year, month int) ([]CalendarData, error) // 月度训练日历
- GetTrainingTimeDistribution(userID int64) (*TimeDistribution, error) // 训练时段分布

// API路由：
GET /api/v1/users/:id/stats           // 用户训练统计
GET /api/v1/users/:id/trend           // 训练趋势（最近30天）
GET /api/v1/users/:id/calendar        // 月度训练日历
GET /api/v1/users/:id/time-distribution // 训练时段分布
```

**B. 前端数据可视化**
```typescript
// src/components/User/TrainingStats.tsx

// 展示内容：
1. 统计卡片
   - 累计训练天数（大数字展示）
   - 累计训练次数（大数字展示）
   - 连续训练天数（大数字展示 + 火焰图标）
   - 本月训练次数（大数字展示）

2. 训练趋势图（折线图）
   - X轴：日期（最近30天）
   - Y轴：训练次数
   - 支持切换时间范围（7天、30天、90天）

3. 月度训练日历（热力图）
   - 类似GitHub贡献图
   - 不同颜色深度表示训练频率
   - 点击日期查看当天详情

4. 训练时段分布（饼图）
   - 早上（6:00-12:00）
   - 下午（12:00-18:00）
   - 晚上（18:00-24:00）
   - 显示各时段占比

5. 训练时长统计（柱状图）
   - X轴：日期
   - Y轴：训练时长（分钟）
   - 显示平均训练时长
```

**验收标准**:
- 统计数据准确无误
- 图表渲染流畅，无卡顿
- 支持时间范围切换
- 数据更新实时生效
- 图表支持导出为图片

---

### 2.3 用户状态管理

#### 2.3.1 用户状态变更功能

**输入**: 用户当前状态
**输出**: 状态变更记录和通知

**执行内容**:

**A. 状态变更API**
```go
// services/user_status_service.go

// 核心方法：
- FreezeUser(userID int64, reason string, operatorID int64) error // 冻结用户
- UnfreezeUser(userID int64, operatorID int64) error // 解冻用户
- AddToBlacklist(userID int64, reason string, operatorID int64) error // 加入黑名单
- RemoveFromBlacklist(userID int64, operatorID int64) error // 移出黑名单
- RecordStatusChange(userID int64, oldStatus, newStatus int8, reason string, operatorID int64) error // 记录状态变更

// API路由：
POST /api/v1/users/:id/freeze        // 冻结用户
POST /api/v1/users/:id/unfreeze      // 解冻用户
POST /api/v1/users/:id/blacklist     // 加入黑名单
DELETE /api/v1/users/:id/blacklist   // 移出黑名单
GET /api/v1/users/:id/status-history // 状态变更历史
```

**B. 状态变更记录表**
```sql
CREATE TABLE user_status_changes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    old_status TINYINT NOT NULL COMMENT '原状态',
    new_status TINYINT NOT NULL COMMENT '新状态',
    reason VARCHAR(255) COMMENT '变更原因',
    operator_id BIGINT NOT NULL COMMENT '操作员ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户状态变更记录表';
```

**C. 前端状态管理界面**
```typescript
// 功能点：
1. 状态变更按钮
   - 冻结/解冻按钮（带确认弹窗）
   - 加入/移出黑名单按钮（带确认弹窗）
   - 需要填写变更原因

2. 状态变更历史
   - 时间线展示
   - 显示操作人、操作时间、变更原因
   - 支持导出

3. 权限控制
   - 只有管理员可以操作
   - 前台人员只能查看
```

**验收标准**:
- 状态变更立即生效
- 状态变更记录完整
- 冻结用户无法签到和预约课程
- 黑名单用户无法登录小程序
- 操作需要二次确认
- 必须填写变更原因

---

### 2.4 用户导入导出

#### 2.4.1 批量导入用户

**输入**: Excel文件
**输出**: 导入结果报告

**执行内容**:

**A. 后端导入API**
```go
// services/user_import_service.go

// 核心方法：
- ImportUsersFromExcel(file multipart.File) (*ImportResult, error)
- ValidateUserData(data *UserData) error // 数据校验
- ParseExcelFile(file multipart.File) ([]*UserData, error) // 解析Excel

// API路由：
POST /api/v1/users/import            // 导入用户
GET  /api/v1/users/import/template   // 下载导入模板
```

**B. Excel模板格式**
```
列名：
- 姓名*（必填）
- 性别（男/女）
- 生日（YYYY-MM-DD）
- 身份证号
- 手机号*（必填）
- 邮箱
- 地址
- 紧急联系人
- 紧急联系电话
- 健康状况
- 训练目标
- 备注
```

**C. 导入逻辑**
```
1. 上传Excel文件
2. 解析Excel数据
3. 逐行校验数据
   - 必填字段检查
   - 格式校验（手机号、身份证、邮箱）
   - 重复性检查（手机号是否已存在）
4. 数据入库
   - 成功的记录直接入库
   - 失败的记录记录错误原因
5. 返回导入结果
   - 成功数量
   - 失败数量
   - 失败详情（行号、错误原因）
```

**D. 前端导入界面**
```typescript
// src/components/User/UserImport.tsx

// 功能点：
1. 下载模板按钮
2. 文件上传组件
   - 支持拖拽上传
   - 文件格式限制：.xlsx, .xls
   - 文件大小限制：最大10MB
3. 导入进度条
4. 导入结果展示
   - 成功数量（绿色）
   - 失败数量（红色）
   - 失败详情表格（行号、姓名、手机号、错误原因）
   - 支持导出失败记录
```

**验收标准**:
- 支持大批量导入（1000+条记录）
- 导入过程有进度提示
- 数据校验准确
- 失败记录可导出重新处理
- 导入过程可中断

---

#### 2.4.2 导出用户数据

**输入**: 筛选条件
**输出**: Excel文件

**执行内容**:

**A. 后端导出API**
```go
// services/user_export_service.go

// 核心方法：
- ExportUsersToExcel(filters *ListUsersRequest) (string, error) // 返回文件路径
- GenerateExcelFile(users []*User) (*bytes.Buffer, error)

// API路由：
POST /api/v1/users/export            // 导出用户（异步）
GET  /api/v1/users/export/:task_id   // 获取导出任务状态
GET  /api/v1/users/download/:file_id // 下载导出文件
```

**B. 导出逻辑**
```
1. 接收筛选条件
2. 查询符合条件的用户数据
3. 生成Excel文件
   - 包含所有字段（敏感信息脱敏）
   - 格式化日期、状态等字段
4. 保存文件到临时目录
5. 返回下载链接
6. 定时清理过期文件（24小时后）
```

**C. 前端导出界面**
```typescript
// 功能点：
1. 导出按钮
   - 导出当前筛选结果
   - 导出全部用户
   - 导出选中用户

2. 导出字段选择
   - 弹窗选择需要导出的字段
   - 支持全选/反选

3. 导出进度提示
   - 数据量大时显示进度条
   - 完成后自动下载

4. 导出记录
   - 显示最近的导出记录
   - 支持重新下载
```

**验收标准**:
- 支持大批量导出（10000+条记录）
- 导出文件格式正确，可正常打开
- 敏感信息脱敏处理
- 导出过程不阻塞其他操作
- 文件下载成功率100%

---

## 三、接口文档

### 3.1 API列表汇总

| 接口路径 | 方法 | 说明 | 权限 |
|---------|------|------|------|
| /api/v1/users | POST | 创建用户 | 管理员、前台 |
| /api/v1/users/:id | GET | 获取用户详情 | 所有角色 |
| /api/v1/users/:id | PUT | 更新用户信息 | 管理员、前台 |
| /api/v1/users/:id | DELETE | 删除用户 | 管理员 |
| /api/v1/users | GET | 用户列表 | 所有角色 |
| /api/v1/users/phone/:phone | GET | 根据手机号查询 | 所有角色 |
| /api/v1/check-in | POST | 签到 | 所有角色 |
| /api/v1/check-out | POST | 签出 | 所有角色 |
| /api/v1/check-in/records | GET | 签到记录列表 | 所有角色 |
| /api/v1/users/:id/stats | GET | 用户训练统计 | 所有角色 |
| /api/v1/users/:id/trend | GET | 训练趋势 | 所有角色 |
| /api/v1/users/:id/freeze | POST | 冻结用户 | 管理员 |
| /api/v1/users/:id/unfreeze | POST | 解冻用户 | 管理员 |
| /api/v1/users/:id/blacklist | POST | 加入黑名单 | 管理员 |
| /api/v1/users/:id/blacklist | DELETE | 移出黑名单 | 管理员 |
| /api/v1/users/import | POST | 导入用户 | 管理员、前台 |
| /api/v1/users/export | POST | 导出用户 | 管理员、前台 |

### 3.2 核心接口详细说明

#### 创建用户
```
POST /api/v1/users

Request Body:
{
  "name": "张三",
  "gender": 1,
  "birthday": "1990-01-01",
  "id_card": "110101199001011234",
  "phone": "13800138000",
  "email": "zhangsan@example.com",
  "address": "北京市朝阳区",
  "emergency_contact": "李四",
  "emergency_phone": "13900139000",
  "health_status": "健康",
  "training_goal": "减脂",
  "source": 1,
  "remark": "备注信息"
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "user_no": "U202401010001",
    "name": "张三",
    ...
  }
}
```

#### 用户列表
```
GET /api/v1/users?page=1&page_size=20&keyword=张三&status=1

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "user_no": "U202401010001",
        "name": "张三",
        "phone": "13800138000",
        "status": 1,
        "training_stats": {
          "total_days": 30,
          "total_times": 45,
          "continuous_days": 5
        },
        "created_at": "2024-01-01T10:00:00Z"
      }
    ],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

## 四、测试用例

### 4.1 单元测试
- 用户创建逻辑测试
- 用户编号生成测试
- 身份证加密解密测试
- 训练统计计算测试
- 连续训练天数计算测试

### 4.2 接口测试
- 所有API接口的正常流程测试
- 参数校验测试
- 权限控制测试
- 并发测试

### 4.3 前端测试
- 表单校验测试
- 列表筛选测试
- 导入导出测试
- 页面交互测试

## 五、上线检查清单

### 5.1 后端检查
- [ ] 数据库表创建完成
- [ ] 所有API接口开发完成
- [ ] 接口文档编写完成
- [ ] 单元测试通过率 > 80%
- [ ] 接口测试全部通过
- [ ] 敏感信息加密处理
- [ ] 日志记录完整
- [ ] 错误处理完善

### 5.2 前端检查
- [ ] 所有页面开发完成
- [ ] 表单校验完整
- [ ] 权限控制正确
- [ ] 页面响应式适配
- [ ] 浏览器兼容性测试
- [ ] 性能优化完成
- [ ] 用户体验流畅

### 5.3 联调检查
- [ ] 前后端接口联调完成
- [ ] 数据流转正常
- [ ] 异常情况处理正确
- [ ] 性能测试通过

## 六、后续优化方向

1. **用户画像**: 基于训练数据生成用户画像，提供个性化推荐
2. **智能提醒**: 根据训练频率自动发送提醒消息
3. **数据分析**: 提供更丰富的数据分析报表
4. **会员等级**: 根据训练数据自动升级会员等级
5. **社交功能**: 添加好友、训练打卡分享等社交功能

---
**任务优先级**: P0（核心功能）  
**预计工期**: 2-3周  
**依赖任务**: 无  
**后续任务**: TASK002（会员卡管理功能依赖用户管理）
