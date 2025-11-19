# TASK006 - 教练管理功能

## 一、功能概述

### 1.1 功能目标
实现健身房教练的完整管理体系，包括教练注册、资质认证、课程绑定、排班管理、业绩统计等核心功能，为私教课预约提供基础支持。

### 1.2 核心价值
- 建立完善的教练档案和资质管理体系
- 实现教练与课程的灵活绑定和管理
- 提供教练业绩统计和考核依据
- 支持教练排班和时间管理
- 为会员选择教练提供信息支持

### 1.3 涉及角色
- **管理员**: 可以管理所有教练信息、审核资质、分配课程、查看业绩
- **教练**: 可以维护自己的信息、设置排班、查看课程安排和业绩
- **前台人员**: 可以查看教练信息、协助排班
- **会员**: 可以查看教练公开信息、评价教练

## 二、功能详细拆解

### 2.1 教练基本信息管理

#### 2.1.1 数据库设计
**输入**: 业务需求分析
**输出**: 数据库表结构SQL文件

**执行内容**:
```sql
-- 教练表 (coaches)
CREATE TABLE coaches (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '教练ID',
    user_id BIGINT UNIQUE NOT NULL COMMENT '关联用户ID',
    coach_no VARCHAR(32) UNIQUE NOT NULL COMMENT '教练编号',
    real_name VARCHAR(50) NOT NULL COMMENT '真实姓名',
    nickname VARCHAR(50) COMMENT '昵称/艺名',
    gender TINYINT COMMENT '性别：1-男，2-女',
    birthday DATE COMMENT '生日',
    phone VARCHAR(11) UNIQUE NOT NULL COMMENT '手机号',
    email VARCHAR(100) COMMENT '邮箱',
    avatar_url VARCHAR(255) COMMENT '头像URL',
    id_card VARCHAR(18) COMMENT '身份证号（加密）',
    
    -- 专业信息
    specialties TEXT COMMENT '擅长领域（JSON数组）',
    certifications TEXT COMMENT '资质证书（JSON数组）',
    experience_years INT DEFAULT 0 COMMENT '从业年限',
    introduction TEXT COMMENT '个人简介',
    achievements TEXT COMMENT '个人成就',
    
    -- 工作信息
    employment_type TINYINT DEFAULT 1 COMMENT '雇佣类型：1-全职，2-兼职',
    entry_date DATE COMMENT '入职日期',
    contract_start_date DATE COMMENT '合同开始日期',
    contract_end_date DATE COMMENT '合同结束日期',
    base_salary DECIMAL(10,2) COMMENT '基本工资',
    commission_rate DECIMAL(5,2) COMMENT '提成比例',
    
    -- 状态信息
    status TINYINT DEFAULT 1 COMMENT '状态：1-在职，2-休假，3-离职',
    rating DECIMAL(3,2) DEFAULT 5.00 COMMENT '评分（1-5）',
    total_ratings INT DEFAULT 0 COMMENT '总评价数',
    total_courses INT DEFAULT 0 COMMENT '总授课数',
    
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL COMMENT '软删除时间',
    
    INDEX idx_user_id (user_id),
    INDEX idx_coach_no (coach_no),
    INDEX idx_phone (phone),
    INDEX idx_status (status),
    INDEX idx_rating (rating),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教练表';

-- 教练资质证书表 (coach_certifications)
CREATE TABLE coach_certifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    coach_id BIGINT NOT NULL COMMENT '教练ID',
    cert_name VARCHAR(100) NOT NULL COMMENT '证书名称',
    cert_no VARCHAR(100) COMMENT '证书编号',
    issuing_authority VARCHAR(100) COMMENT '颁发机构',
    issue_date DATE COMMENT '颁发日期',
    expiry_date DATE COMMENT '过期日期',
    cert_image_url VARCHAR(255) COMMENT '证书图片URL',
    status TINYINT DEFAULT 1 COMMENT '状态：1-待审核，2-已认证，3-已过期，4-已拒绝',
    verify_time DATETIME COMMENT '审核时间',
    verifier_id BIGINT COMMENT '审核人ID',
    verify_remark VARCHAR(255) COMMENT '审核备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_coach_id (coach_id),
    INDEX idx_status (status),
    FOREIGN KEY (coach_id) REFERENCES coaches(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教练资质证书表';

-- 教练业绩统计表 (coach_performance)
CREATE TABLE coach_performance (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    coach_id BIGINT NOT NULL COMMENT '教练ID',
    stat_date DATE NOT NULL COMMENT '统计日期',
    stat_type TINYINT NOT NULL COMMENT '统计类型：1-日，2-周，3-月，4-年',
    
    -- 课程统计
    total_bookings INT DEFAULT 0 COMMENT '总预约数',
    completed_bookings INT DEFAULT 0 COMMENT '已完成数',
    cancelled_bookings INT DEFAULT 0 COMMENT '已取消数',
    completion_rate DECIMAL(5,2) COMMENT '完成率',
    
    -- 评价统计
    total_ratings INT DEFAULT 0 COMMENT '总评价数',
    avg_rating DECIMAL(3,2) COMMENT '平均评分',
    five_star_count INT DEFAULT 0 COMMENT '5星评价数',
    four_star_count INT DEFAULT 0 COMMENT '4星评价数',
    three_star_count INT DEFAULT 0 COMMENT '3星评价数',
    two_star_count INT DEFAULT 0 COMMENT '2星评价数',
    one_star_count INT DEFAULT 0 COMMENT '1星评价数',
    
    -- 收入统计
    total_income DECIMAL(10,2) DEFAULT 0 COMMENT '总收入',
    commission_income DECIMAL(10,2) DEFAULT 0 COMMENT '提成收入',
    
    -- 时长统计
    total_hours DECIMAL(10,2) DEFAULT 0 COMMENT '总授课时长（小时）',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_coach_date_type (coach_id, stat_date, stat_type),
    INDEX idx_coach_id (coach_id),
    INDEX idx_stat_date (stat_date),
    FOREIGN KEY (coach_id) REFERENCES coaches(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教练业绩统计表';

-- 教练工作日志表 (coach_work_logs)
CREATE TABLE coach_work_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    coach_id BIGINT NOT NULL COMMENT '教练ID',
    log_date DATE NOT NULL COMMENT '日志日期',
    work_hours DECIMAL(5,2) COMMENT '工作时长（小时）',
    course_count INT DEFAULT 0 COMMENT '授课数量',
    content TEXT COMMENT '工作内容',
    mood TINYINT COMMENT '心情：1-很好，2-好，3-一般，4-差，5-很差',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_coach_id (coach_id),
    INDEX idx_log_date (log_date),
    FOREIGN KEY (coach_id) REFERENCES coaches(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教练工作日志表';
```

**验收标准**:
- 表结构符合第三范式
- 包含必要的索引优化查询性能
- 敏感字段加密存储
- 支持软删除机制
- 支持业绩统计和分析

---

#### 2.1.2 后端API开发 - 教练CRUD
**输入**: 数据库表结构
**输出**: 完整的教练管理API接口

**执行内容**:

**A. 定义数据模型 (models/coach.go)**
```go
type Coach struct {
    ID                 int64      `json:"id" gorm:"primaryKey"`
    UserID             int64      `json:"user_id" gorm:"uniqueIndex;not null"`
    CoachNo            string     `json:"coach_no" gorm:"uniqueIndex;size:32"`
    RealName           string     `json:"real_name" gorm:"size:50;not null"`
    Nickname           string     `json:"nickname" gorm:"size:50"`
    Gender             int8       `json:"gender"`
    Birthday           *time.Time `json:"birthday"`
    Phone              string     `json:"phone" gorm:"uniqueIndex;size:11;not null"`
    Email              string     `json:"email" gorm:"size:100"`
    AvatarURL          string     `json:"avatar_url" gorm:"size:255"`
    IDCard             string     `json:"id_card" gorm:"size:18"`
    Specialties        string     `json:"specialties" gorm:"type:text"`
    Certifications     string     `json:"certifications" gorm:"type:text"`
    ExperienceYears    int        `json:"experience_years" gorm:"default:0"`
    Introduction       string     `json:"introduction" gorm:"type:text"`
    Achievements       string     `json:"achievements" gorm:"type:text"`
    EmploymentType     int8       `json:"employment_type" gorm:"default:1"`
    EntryDate          *time.Time `json:"entry_date"`
    ContractStartDate  *time.Time `json:"contract_start_date"`
    ContractEndDate    *time.Time `json:"contract_end_date"`
    BaseSalary         float64    `json:"base_salary" gorm:"type:decimal(10,2)"`
    CommissionRate     float64    `json:"commission_rate" gorm:"type:decimal(5,2)"`
    Status             int8       `json:"status" gorm:"default:1"`
    Rating             float64    `json:"rating" gorm:"type:decimal(3,2);default:5.00"`
    TotalRatings       int        `json:"total_ratings" gorm:"default:0"`
    TotalCourses       int        `json:"total_courses" gorm:"default:0"`
    Remark             string     `json:"remark" gorm:"type:text"`
    CreatedAt          time.Time  `json:"created_at"`
    UpdatedAt          time.Time  `json:"updated_at"`
    DeletedAt          *time.Time `json:"deleted_at" gorm:"index"`
}

type CoachCertification struct {
    ID               int64      `json:"id" gorm:"primaryKey"`
    CoachID          int64      `json:"coach_id" gorm:"not null"`
    CertName         string     `json:"cert_name" gorm:"size:100;not null"`
    CertNo           string     `json:"cert_no" gorm:"size:100"`
    IssuingAuthority string     `json:"issuing_authority" gorm:"size:100"`
    IssueDate        *time.Time `json:"issue_date"`
    ExpiryDate       *time.Time `json:"expiry_date"`
    CertImageURL     string     `json:"cert_image_url" gorm:"size:255"`
    Status           int8       `json:"status" gorm:"default:1"`
    VerifyTime       *time.Time `json:"verify_time"`
    VerifierID       *int64     `json:"verifier_id"`
    VerifyRemark     string     `json:"verify_remark" gorm:"size:255"`
    CreatedAt        time.Time  `json:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at"`
}

type CoachPerformance struct {
    ID                 int64     `json:"id" gorm:"primaryKey"`
    CoachID            int64     `json:"coach_id" gorm:"not null"`
    StatDate           time.Time `json:"stat_date" gorm:"type:date;not null"`
    StatType           int8      `json:"stat_type" gorm:"not null"`
    TotalBookings      int       `json:"total_bookings" gorm:"default:0"`
    CompletedBookings  int       `json:"completed_bookings" gorm:"default:0"`
    CancelledBookings  int       `json:"cancelled_bookings" gorm:"default:0"`
    CompletionRate     float64   `json:"completion_rate" gorm:"type:decimal(5,2)"`
    TotalRatings       int       `json:"total_ratings" gorm:"default:0"`
    AvgRating          float64   `json:"avg_rating" gorm:"type:decimal(3,2)"`
    FiveStarCount      int       `json:"five_star_count" gorm:"default:0"`
    FourStarCount      int       `json:"four_star_count" gorm:"default:0"`
    ThreeStarCount     int       `json:"three_star_count" gorm:"default:0"`
    TwoStarCount       int       `json:"two_star_count" gorm:"default:0"`
    OneStarCount       int       `json:"one_star_count" gorm:"default:0"`
    TotalIncome        float64   `json:"total_income" gorm:"type:decimal(10,2);default:0"`
    CommissionIncome   float64   `json:"commission_income" gorm:"type:decimal(10,2);default:0"`
    TotalHours         float64   `json:"total_hours" gorm:"type:decimal(10,2);default:0"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
}
```

**B. 实现业务逻辑层 (services/coach_service.go)**
```go
// 核心方法列表
- CreateCoach(req *CreateCoachRequest) (*Coach, error)
- GetCoachByID(id int64) (*Coach, error)
- GetCoachByUserID(userID int64) (*Coach, error)
- UpdateCoach(id int64, req *UpdateCoachRequest) error
- DeleteCoach(id int64) error
- ListCoaches(req *ListCoachesRequest) ([]*Coach, int64, error)
- GenerateCoachNo() string
- UpdateCoachRating(coachID int64, rating int8) error
- IncrementCoachCourseCount(coachID int64) error
- GetCoachWithDetails(id int64) (*CoachDetailVO, error) // 包含课程、评价等信息
```

**C. 实现API接口层 (controllers/coach_controller.go)**
```go
// API路由定义
POST   /api/v1/coaches              // 创建教练
GET    /api/v1/coaches/:id          // 获取教练详情
PUT    /api/v1/coaches/:id          // 更新教练信息
DELETE /api/v1/coaches/:id          // 删除教练
GET    /api/v1/coaches              // 教练列表
GET    /api/v1/coaches/:id/details  // 获取教练详细信息（含课程、评价）
PUT    /api/v1/coaches/:id/status   // 更新教练状态
```

**D. 请求/响应结构定义**
```go
type CreateCoachRequest struct {
    UserID          int64    `json:"user_id" binding:"required"`
    RealName        string   `json:"real_name" binding:"required,max=50"`
    Nickname        string   `json:"nickname" binding:"omitempty,max=50"`
    Gender          int8     `json:"gender" binding:"omitempty,oneof=1 2"`
    Birthday        string   `json:"birthday" binding:"omitempty"`
    Phone           string   `json:"phone" binding:"required,len=11"`
    Email           string   `json:"email" binding:"omitempty,email"`
    IDCard          string   `json:"id_card" binding:"omitempty,len=18"`
    Specialties     []string `json:"specialties"`
    ExperienceYears int      `json:"experience_years"`
    Introduction    string   `json:"introduction"`
    EmploymentType  int8     `json:"employment_type" binding:"omitempty,oneof=1 2"`
    EntryDate       string   `json:"entry_date"`
    BaseSalary      float64  `json:"base_salary"`
    CommissionRate  float64  `json:"commission_rate"`
}

type ListCoachesRequest struct {
    Page       int     `form:"page" binding:"required,min=1"`
    PageSize   int     `form:"page_size" binding:"required,min=1,max=100"`
    Keyword    string  `form:"keyword"`
    Status     int8    `form:"status"`
    Specialty  string  `form:"specialty"`
    MinRating  float64 `form:"min_rating"`
    SortBy     string  `form:"sort_by"` // rating, courses, experience
    SortOrder  string  `form:"sort_order"` // asc, desc
}
```

**验收标准**:
- 所有API接口通过测试
- 教练编号自动生成且唯一
- 支持按多条件筛选和排序
- 评分更新准确
- 接口响应时间 < 200ms

---

### 2.2 教练资质认证管理

#### 2.2.1 资质证书上传与审核
**输入**: 证书信息和图片
**输出**: 证书记录和审核状态

**执行内容**:

**A. 后端API实现**
```go
// services/coach_certification_service.go

// 核心方法：
- UploadCertification(coachID int64, req *UploadCertificationRequest) (*CoachCertification, error)
- UpdateCertification(id int64, req *UpdateCertificationRequest) error
- DeleteCertification(id int64) error
- GetCoachCertifications(coachID int64) ([]*CoachCertification, error)
- VerifyCertification(id int64, verifierID int64, status int8, remark string) error
- CheckExpiredCertifications() error // 定时任务：检查过期证书
- GetPendingCertifications() ([]*CoachCertification, error) // 获取待审核证书列表

// API路由：
POST   /api/v1/coaches/:id/certifications        // 上传证书
PUT    /api/v1/certifications/:id                // 更新证书
DELETE /api/v1/certifications/:id                // 删除证书
GET    /api/v1/coaches/:id/certifications        // 获取教练证书列表
PUT    /api/v1/certifications/:id/verify         // 审核证书
GET    /api/v1/certifications/pending            // 获取待审核证书列表
```

**B. 请求/响应结构**
```go
type UploadCertificationRequest struct {
    CertName         string `json:"cert_name" binding:"required,max=100"`
    CertNo           string `json:"cert_no" binding:"omitempty,max=100"`
    IssuingAuthority string `json:"issuing_authority" binding:"omitempty,max=100"`
    IssueDate        string `json:"issue_date" binding:"omitempty"`
    ExpiryDate       string `json:"expiry_date" binding:"omitempty"`
    CertImageURL     string `json:"cert_image_url" binding:"required"`
}

type VerifyCertificationRequest struct {
    Status int8   `json:"status" binding:"required,oneof=2 4"` // 2-已认证，4-已拒绝
    Remark string `json:"remark" binding:"omitempty,max=255"`
}
```

**C. 业务逻辑**
```
1. 上传证书流程：
   - 验证教练是否存在
   - 上传证书图片到OSS
   - 创建证书记录（待审核状态）
   - 发送通知给管理员
   - 返回证书信息

2. 审核证书流程：
   - 验证证书是否存在
   - 验证审核人权限
   - 更新证书状态
   - 记录审核人和审核时间
   - 发送通知给教练
   - 返回成功

3. 检查过期证书（定时任务）：
   - 查询所有已认证的证书
   - 检查过期日期
   - 更新过期证书状态
   - 发送过期提醒给教练和管理员
```

**验收标准**:
- 证书图片上传成功
- 审核流程完整
- 过期检查准确
- 通知发送及时

---

#### 2.2.2 前端页面开发 - 资质管理
**输入**: 后端API接口文档
**输出**: 资质管理页面

**执行内容**:

**A. 证书上传页面 (src/pages/Coach/CertificationUpload.tsx)**
```typescript
// 页面布局：
1. 证书信息表单
   - 证书名称（必填）
   - 证书编号
   - 颁发机构
   - 颁发日期
   - 过期日期
   - 证书图片上传

2. 图片上传组件
   - 支持拖拽上传
   - 图片预览
   - 支持裁剪
   - 格式限制：jpg、png
   - 大小限制：最大5MB

3. 提交按钮
   - 验证表单
   - 提交审核
```

**B. 证书列表页面 (src/pages/Coach/CertificationList.tsx)**
```typescript
// 页面布局：
1. 证书卡片列表
   - 证书图片缩略图
   - 证书名称
   - 颁发机构
   - 颁发日期
   - 过期日期
   - 状态标签（待审核、已认证、已过期、已拒绝）
   - 操作按钮（查看、编辑、删除）

2. 状态筛选
   - 全部
   - 待审核
   - 已认证
   - 已过期
   - 已拒绝

3. 证书详情弹窗
   - 大图展示
   - 完整信息
   - 审核记录（如有）
```

**C. 证书审核页面（管理员）(src/pages/Admin/CertificationReview.tsx)**
```typescript
// 页面布局：
1. 待审核列表
   - 教练信息
   - 证书信息
   - 提交时间
   - 操作按钮（审核）

2. 审核弹窗
   - 证书大图
   - 证书详细信息
   - 审核选项（通过、拒绝）
   - 审核备注
   - 提交按钮
```

**验收标准**:
- 证书上传流畅
- 图片预览清晰
- 审核流程完整
- 状态更新实时

---

### 2.3 教练业绩统计

#### 2.3.1 业绩数据统计
**输入**: 教练ID、时间范围
**输出**: 业绩统计数据

**执行内容**:

**A. 后端统计API**
```go
// services/coach_performance_service.go

// 核心方法：
- CalculateDailyPerformance(coachID int64, date time.Time) error
- CalculateMonthlyPerformance(coachID int64, year, month int) error
- GetCoachPerformance(coachID int64, statType int8, startDate, endDate time.Time) ([]*CoachPerformance, error)
- GetCoachRanking(statType int8, date time.Time, limit int) ([]*CoachRankingVO, error)
- GetCoachIncome(coachID int64, startDate, endDate time.Time) (*IncomeStatVO, error)
- ExportPerformanceReport(coachID int64, startDate, endDate time.Time) (string, error)

// API路由：
GET /api/v1/coaches/:id/performance        // 获取教练业绩
GET /api/v1/coaches/:id/income             // 获取教练收入统计
GET /api/v1/coaches/ranking                // 获取教练排行榜
POST /api/v1/coaches/:id/performance/export // 导出业绩报表
```

**B. 定时任务**
```go
// 每日凌晨1点执行
func DailyPerformanceTask() {
    // 统计所有教练昨日业绩
    // 更新教练总评分和总课程数
    // 生成日报
}

// 每月1号凌晨2点执行
func MonthlyPerformanceTask() {
    // 统计所有教练上月业绩
    // 生成月报
    // 发送业绩通知
}
```

**C. 业务逻辑**
```
1. 日业绩统计：
   - 统计当日预约数、完成数、取消数
   - 计算完成率
   - 统计当日评价数和平均评分
   - 计算当日收入和提成
   - 统计当日授课时长
   - 保存统计数据

2. 月业绩统计：
   - 汇总当月所有日统计数据
   - 计算月度指标
   - 生成月度报表
   - 保存统计数据

3. 排行榜计算：
   - 按预约量排行
   - 按评分排行
   - 按收入排行
   - 按完成率排行
```

**验收标准**:
- 统计数据准确
- 定时任务稳定执行
- 排行榜实时更新
- 报表导出正常

---

#### 2.3.2 前端页面开发 - 业绩展示
**输入**: 后端API接口文档
**输出**: 业绩统计页面

**执行内容**:

**A. 教练业绩页面 (src/pages/Coach/Performance.tsx)**
```typescript
// 页面布局：
1. 时间范围选择
   - 今日、本周、本月、自定义

2. 统计卡片
   - 总预约数
   - 已完成数
   - 完成率
   - 平均评分
   - 总收入
   - 提成收入
   - 授课时长

3. 业绩趋势图（折线图）
   - X轴：日期
   - Y轴：预约数/收入
   - 支持切换指标

4. 评分分布图（柱状图）
   - 各星级评价数量

5. 课程类型分布（饼图）
   - 各课程类型授课占比

6. 详细数据表格
   - 日期
   - 预约数
   - 完成数
   - 评分
   - 收入
   - 时长
```

**B. 教练排行榜页面 (src/pages/Coach/Ranking.tsx)**
```typescript
// 页面布局：
1. 排行榜类型切换
   - 预约量排行
   - 评分排行
   - 收入排行
   - 完成率排行

2. 时间范围选择
   - 本周、本月、本年

3. 排行榜列表
   - 排名
   - 教练头像和姓名
   - 对应指标数值
   - 趋势图标（上升/下降）

4. 我的排名（高亮显示）
```

**C. 收入统计页面 (src/pages/Coach/Income.tsx)**
```typescript
// 页面布局：
1. 收入概览
   - 总收入
   - 基本工资
   - 提成收入
   - 奖金

2. 收入趋势图（折线图）
   - 按月展示收入变化

3. 收入明细表格
   - 日期
   - 课程数
   - 基本工资
   - 提成收入
   - 奖金
   - 总收入

4. 导出功能
   - 导出收入报表
```

**验收标准**:
- 数据展示准确
- 图表渲染流畅
- 支持时间范围切换
- 排行榜实时更新
- 报表导出成功

---

### 2.4 教练工作日志

#### 2.4.1 工作日志记录
**输入**: 日志内容
**输出**: 日志记录

**执行内容**:

**A. 后端API实现**
```go
// services/coach_work_log_service.go

// 核心方法：
- CreateWorkLog(coachID int64, req *CreateWorkLogRequest) (*CoachWorkLog, error)
- UpdateWorkLog(id int64, req *UpdateWorkLogRequest) error
- DeleteWorkLog(id int64) error
- GetWorkLog(id int64) (*CoachWorkLog, error)
- ListWorkLogs(coachID int64, startDate, endDate time.Time) ([]*CoachWorkLog, error)
- GetWorkLogByDate(coachID int64, date time.Time) (*CoachWorkLog, error)

// API路由：
POST   /api/v1/coaches/:id/work-logs     // 创建工作日志
PUT    /api/v1/work-logs/:id             // 更新工作日志
DELETE /api/v1/work-logs/:id             // 删除工作日志
GET    /api/v1/work-logs/:id             // 获取工作日志
GET    /api/v1/coaches/:id/work-logs     // 获取工作日志列表
```

**B. 前端页面开发**
```typescript
// src/pages/Coach/WorkLog.tsx

// 页面布局：
1. 日历视图
   - 月视图
   - 有日志的日期标记
   - 点击日期查看/编辑日志

2. 日志编辑表单
   - 日期（自动填充）
   - 工作时长
   - 授课数量（自动统计）
   - 工作内容（文本域）
   - 心情选择（表情图标）
   - 保存按钮

3. 日志列表
   - 按日期倒序
   - 显示日期、工作时长、授课数、心情
   - 操作按钮（查看、编辑、删除）

4. 统计卡片
   - 本月工作天数
   - 本月总工作时长
   - 本月总授课数
   - 平均心情指数
```

**验收标准**:
- 日志保存成功
- 日历视图清晰
- 统计数据准确
- 支持快速编辑

---

## 三、接口文档

### 3.1 API列表汇总

| 接口路径 | 方法 | 说明 | 权限 |
|---------|------|------|------|
| /api/v1/coaches | POST | 创建教练 | 管理员 |
| /api/v1/coaches/:id | GET | 获取教练详情 | 所有角色 |
| /api/v1/coaches/:id | PUT | 更新教练信息 | 管理员、教练本人 |
| /api/v1/coaches/:id | DELETE | 删除教练 | 管理员 |
| /api/v1/coaches | GET | 教练列表 | 所有角色 |
| /api/v1/coaches/:id/details | GET | 获取教练详细信息 | 所有角色 |
| /api/v1/coaches/:id/status | PUT | 更新教练状态 | 管理员 |
| /api/v1/coaches/:id/certifications | POST | 上传证书 | 教练本人 |
| /api/v1/certifications/:id | PUT | 更新证书 | 教练本人 |
| /api/v1/certifications/:id | DELETE | 删除证书 | 教练本人、管理员 |
| /api/v1/coaches/:id/certifications | GET | 获取证书列表 | 所有角色 |
| /api/v1/certifications/:id/verify | PUT | 审核证书 | 管理员 |
| /api/v1/certifications/pending | GET | 获取待审核证书 | 管理员 |
| /api/v1/coaches/:id/performance | GET | 获取业绩统计 | 教练本人、管理员 |
| /api/v1/coaches/:id/income | GET | 获取收入统计 | 教练本人、管理员 |
| /api/v1/coaches/ranking | GET | 获取排行榜 | 所有角色 |
| /api/v1/coaches/:id/work-logs | POST | 创建工作日志 | 教练本人 |
| /api/v1/coaches/:id/work-logs | GET | 获取工作日志列表 | 教练本人、管理员 |

### 3.2 核心接口详细说明

#### 创建教练
```
POST /api/v1/coaches

Request Body:
{
  "user_id": 1,
  "real_name": "张教练",
  "nickname": "小张",
  "gender": 1,
  "birthday": "1990-01-01",
  "phone": "13800138000",
  "email": "coach@example.com",
  "specialties": ["减脂", "增肌", "康复训练"],
  "experience_years": 5,
  "introduction": "专业健身教练，擅长减脂增肌",
  "employment_type": 1,
  "entry_date": "2024-01-01",
  "base_salary": 5000.00,
  "commission_rate": 30.00
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "coach_no": "C202401010001",
    "real_name": "张教练",
    "rating": 5.00,
    "status": 1,
    ...
  }
}
```

#### 获取教练业绩
```
GET /api/v1/coaches/1/performance?stat_type=3&start_date=2024-01-01&end_date=2024-01-31

Response:
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "stat_date": "2024-01-01",
      "total_bookings": 10,
      "completed_bookings": 9,
      "cancelled_bookings": 1,
      "completion_rate": 90.00,
      "avg_rating": 4.8,
      "total_income": 3000.00,
      "commission_income": 900.00,
      "total_hours": 10.5
    }
  ]
}
```

## 四、测试用例

### 4.1 单元测试
- 教练编号生成测试
- 评分更新测试
- 业绩统计计算测试
- 证书过期检查测试

### 4.2 接口测试
- 所有API接口的正常流程测试
- 权限控制测试
- 数据校验测试
- 并发测试

### 4.3 前端测试
- 教练信息展示测试
- 证书上传测试
- 业绩图表测试
- 工作日志测试

## 五、上线检查清单

### 5.1 后端检查
- [ ] 数据库表创建完成
- [ ] 所有API接口开发完成
- [ ] 业绩统计逻辑正确
- [ ] 定时任务配置完成
- [ ] 单元测试通过率 > 80%
- [ ] 接口测试全部通过

### 5.2 前端检查
- [ ] 所有页面开发完成
- [ ] 教练信息展示完整
- [ ] 证书管理功能正常
- [ ] 业绩图表展示正确
- [ ] 页面响应式适配

### 5.3 联调检查
- [ ] 前后端接口联调完成
- [ ] 业绩统计准确
- [ ] 证书审核流程正常
- [ ] 性能测试通过

## 六、后续优化方向

1. **智能排班**: 基于历史数据和会员偏好自动推荐排班
2. **教练培训**: 添加教练培训课程和考核系统
3. **教练社区**: 建立教练交流社区，分享经验
4. **AI助手**: 提供AI助手帮助教练制定训练计划
5. **视频教学**: 支持教练录制教学视频
6. **会员匹配**: 基于会员目标和教练专长智能匹配

---
**任务优先级**: P0（核心功能）  
**预计工期**: 2-3周  
**依赖任务**: TASK001（用户管理）  
**后续任务**: TASK005（私教课预约功能）

