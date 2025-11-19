# TASK005 - 私教课预约功能

## 一、功能概述

### 1.1 功能目标
实现会员与私教之间的课程预约管理系统，包括课程创建、预约、取消、签到、评价等完整流程，支持时间冲突检测和可视化排课。

### 1.2 核心价值
- 提升私教课程管理效率，减少人工协调成本
- 避免课程时间冲突，优化教练和场地资源利用
- 提供可视化排课界面，直观展示课程安排
- 记录课程数据，为教练考核和会员服务提供依据

### 1.3 涉及角色
- **管理员**: 可以管理所有课程、查看所有预约记录、设置课程规则
- **教练**: 可以设置自己的可约时间、查看自己的课程安排、确认/取消课程
- **会员**: 可以浏览教练信息、预约课程、取消预约、课后评价
- **前台人员**: 可以代会员预约课程、处理预约纠纷

## 二、功能详细拆解

### 2.1 私教课程类型管理

#### 2.1.1 数据库设计
**输入**: 业务需求分析
**输出**: 数据库表结构SQL文件

**执行内容**:
```sql
-- 课程类型表 (course_types)
CREATE TABLE course_types (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '课程类型ID',
    name VARCHAR(50) NOT NULL COMMENT '课程名称',
    description TEXT COMMENT '课程描述',
    duration INT NOT NULL COMMENT '课程时长（分钟）',
    price DECIMAL(10,2) NOT NULL COMMENT '单次价格',
    icon_url VARCHAR(255) COMMENT '图标URL',
    sort_order INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态：1-启用，2-禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_sort_order (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程类型表';

-- 教练课程关联表 (coach_courses)
CREATE TABLE coach_courses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    coach_id BIGINT NOT NULL COMMENT '教练ID',
    course_type_id BIGINT NOT NULL COMMENT '课程类型ID',
    is_active TINYINT DEFAULT 1 COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_coach_course (coach_id, course_type_id),
    FOREIGN KEY (coach_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (course_type_id) REFERENCES course_types(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教练课程关联表';

-- 教练可约时间表 (coach_available_times)
CREATE TABLE coach_available_times (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    coach_id BIGINT NOT NULL COMMENT '教练ID',
    day_of_week TINYINT NOT NULL COMMENT '星期几：1-7（周一到周日）',
    start_time TIME NOT NULL COMMENT '开始时间',
    end_time TIME NOT NULL COMMENT '结束时间',
    is_active TINYINT DEFAULT 1 COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_coach_id (coach_id),
    FOREIGN KEY (coach_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教练可约时间表';

-- 教练请假记录表 (coach_leaves)
CREATE TABLE coach_leaves (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    coach_id BIGINT NOT NULL COMMENT '教练ID',
    start_time DATETIME NOT NULL COMMENT '请假开始时间',
    end_time DATETIME NOT NULL COMMENT '请假结束时间',
    reason VARCHAR(255) COMMENT '请假原因',
    status TINYINT DEFAULT 1 COMMENT '状态：1-待审核，2-已批准，3-已拒绝',
    approver_id BIGINT COMMENT '审批人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_coach_id (coach_id),
    INDEX idx_time_range (start_time, end_time),
    FOREIGN KEY (coach_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='教练请假记录表';

-- 课程预约表 (course_bookings)
CREATE TABLE course_bookings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '预约ID',
    booking_no VARCHAR(32) UNIQUE NOT NULL COMMENT '预约编号',
    user_id BIGINT NOT NULL COMMENT '会员ID',
    coach_id BIGINT NOT NULL COMMENT '教练ID',
    course_type_id BIGINT NOT NULL COMMENT '课程类型ID',
    membership_card_id BIGINT COMMENT '会员卡ID（扣课时）',
    booking_date DATE NOT NULL COMMENT '预约日期',
    start_time TIME NOT NULL COMMENT '开始时间',
    end_time TIME NOT NULL COMMENT '结束时间',
    duration INT NOT NULL COMMENT '课程时长（分钟）',
    status TINYINT DEFAULT 1 COMMENT '状态：1-待确认，2-已确认，3-已完成，4-已取消，5-已过期',
    cancel_reason VARCHAR(255) COMMENT '取消原因',
    cancel_by BIGINT COMMENT '取消人ID',
    cancel_time DATETIME COMMENT '取消时间',
    check_in_time DATETIME COMMENT '签到时间',
    check_out_time DATETIME COMMENT '签出时间',
    actual_duration INT COMMENT '实际时长（分钟）',
    rating TINYINT COMMENT '评分：1-5',
    comment TEXT COMMENT '评价内容',
    comment_time DATETIME COMMENT '评价时间',
    remark TEXT COMMENT '备注',
    created_by BIGINT COMMENT '创建人ID（前台代约时使用）',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_coach_id (coach_id),
    INDEX idx_booking_date (booking_date),
    INDEX idx_status (status),
    INDEX idx_booking_no (booking_no),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (coach_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (course_type_id) REFERENCES course_types(id),
    FOREIGN KEY (membership_card_id) REFERENCES membership_cards(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='课程预约表';
```

**验收标准**:
- 表结构符合第三范式
- 包含必要的索引优化查询性能
- 支持时间冲突检测
- 支持软删除和状态管理

---

#### 2.1.2 后端API开发 - 课程类型管理
**输入**: 数据库表结构
**输出**: 课程类型管理API接口

**执行内容**:

**A. 定义数据模型 (models/course.go)**
```go
type CourseType struct {
    ID          int64     `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"size:50;not null"`
    Description string    `json:"description" gorm:"type:text"`
    Duration    int       `json:"duration" gorm:"not null"`
    Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
    IconURL     string    `json:"icon_url" gorm:"size:255"`
    SortOrder   int       `json:"sort_order" gorm:"default:0"`
    Status      int8      `json:"status" gorm:"default:1"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CoachCourse struct {
    ID           int64     `json:"id" gorm:"primaryKey"`
    CoachID      int64     `json:"coach_id" gorm:"not null"`
    CourseTypeID int64     `json:"course_type_id" gorm:"not null"`
    IsActive     int8      `json:"is_active" gorm:"default:1"`
    CreatedAt    time.Time `json:"created_at"`
}

type CoachAvailableTime struct {
    ID        int64     `json:"id" gorm:"primaryKey"`
    CoachID   int64     `json:"coach_id" gorm:"not null"`
    DayOfWeek int8      `json:"day_of_week" gorm:"not null"`
    StartTime string    `json:"start_time" gorm:"type:time;not null"`
    EndTime   string    `json:"end_time" gorm:"type:time;not null"`
    IsActive  int8      `json:"is_active" gorm:"default:1"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CourseBooking struct {
    ID               int64      `json:"id" gorm:"primaryKey"`
    BookingNo        string     `json:"booking_no" gorm:"uniqueIndex;size:32"`
    UserID           int64      `json:"user_id" gorm:"not null"`
    CoachID          int64      `json:"coach_id" gorm:"not null"`
    CourseTypeID     int64      `json:"course_type_id" gorm:"not null"`
    MembershipCardID *int64     `json:"membership_card_id"`
    BookingDate      time.Time  `json:"booking_date" gorm:"type:date;not null"`
    StartTime        string     `json:"start_time" gorm:"type:time;not null"`
    EndTime          string     `json:"end_time" gorm:"type:time;not null"`
    Duration         int        `json:"duration" gorm:"not null"`
    Status           int8       `json:"status" gorm:"default:1"`
    CancelReason     string     `json:"cancel_reason" gorm:"size:255"`
    CancelBy         *int64     `json:"cancel_by"`
    CancelTime       *time.Time `json:"cancel_time"`
    CheckInTime      *time.Time `json:"check_in_time"`
    CheckOutTime     *time.Time `json:"check_out_time"`
    ActualDuration   *int       `json:"actual_duration"`
    Rating           *int8      `json:"rating"`
    Comment          string     `json:"comment" gorm:"type:text"`
    CommentTime      *time.Time `json:"comment_time"`
    Remark           string     `json:"remark" gorm:"type:text"`
    CreatedBy        *int64     `json:"created_by"`
    CreatedAt        time.Time  `json:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at"`
}
```

**B. 实现业务逻辑层 (services/course_service.go)**
```go
// 核心方法列表
- CreateCourseType(req *CreateCourseTypeRequest) (*CourseType, error)
- UpdateCourseType(id int64, req *UpdateCourseTypeRequest) error
- DeleteCourseType(id int64) error
- GetCourseTypeByID(id int64) (*CourseType, error)
- ListCourseTypes(status int8) ([]*CourseType, error)
- AssignCourseToCoach(coachID, courseTypeID int64) error
- RemoveCourseFromCoach(coachID, courseTypeID int64) error
- GetCoachCourses(coachID int64) ([]*CourseType, error)
```

**C. 实现API接口层 (controllers/course_controller.go)**
```go
// API路由定义
POST   /api/v1/course-types              // 创建课程类型
PUT    /api/v1/course-types/:id          // 更新课程类型
DELETE /api/v1/course-types/:id          // 删除课程类型
GET    /api/v1/course-types/:id          // 获取课程类型详情
GET    /api/v1/course-types              // 课程类型列表
POST   /api/v1/coaches/:id/courses       // 为教练分配课程
DELETE /api/v1/coaches/:id/courses/:course_id // 移除教练课程
GET    /api/v1/coaches/:id/courses       // 获取教练课程列表
```

**验收标准**:
- 所有API接口通过测试
- 课程类型支持排序
- 教练课程关联正确
- 接口响应时间 < 200ms

---

### 2.2 教练可约时间管理

#### 2.2.1 教练设置可约时间
**输入**: 教练ID和时间段
**输出**: 可约时间配置

**执行内容**:

**A. 后端API实现**
```go
// services/coach_schedule_service.go

// 核心方法：
- SetAvailableTime(coachID int64, req *SetAvailableTimeRequest) error
- GetAvailableTimes(coachID int64) ([]*CoachAvailableTime, error)
- UpdateAvailableTime(id int64, req *UpdateAvailableTimeRequest) error
- DeleteAvailableTime(id int64) error
- CheckTimeConflict(coachID int64, dayOfWeek int8, startTime, endTime string) bool
- GetAvailableSlots(coachID int64, date time.Time) ([]*TimeSlot, error) // 获取某天的可约时间段

// API路由：
POST   /api/v1/coaches/:id/available-times     // 设置可约时间
GET    /api/v1/coaches/:id/available-times     // 获取可约时间列表
PUT    /api/v1/coaches/:id/available-times/:time_id // 更新可约时间
DELETE /api/v1/coaches/:id/available-times/:time_id // 删除可约时间
GET    /api/v1/coaches/:id/available-slots     // 获取某天的可约时间段（供会员预约使用）
```

**B. 请求/响应结构**
```go
type SetAvailableTimeRequest struct {
    DayOfWeek int8   `json:"day_of_week" binding:"required,min=1,max=7"`
    StartTime string `json:"start_time" binding:"required"` // HH:MM格式
    EndTime   string `json:"end_time" binding:"required"`   // HH:MM格式
}

type TimeSlot struct {
    StartTime string `json:"start_time"`
    EndTime   string `json:"end_time"`
    IsBooked  bool   `json:"is_booked"`
    BookingID *int64 `json:"booking_id"`
}
```

**C. 业务逻辑**
```
1. 设置可约时间流程：
   - 验证时间格式
   - 验证开始时间 < 结束时间
   - 检查时间段是否冲突（同一天的时间段不能重叠）
   - 保存可约时间配置
   - 返回成功

2. 获取可约时间段流程：
   - 根据日期获取星期几
   - 查询教练该星期的可约时间配置
   - 查询该日期已预约的课程
   - 计算可用时间段（排除已预约和请假时间）
   - 按时间顺序返回
```

**验收标准**:
- 时间段冲突检测准确
- 可约时间配置灵活
- 支持批量设置（一次设置多天）
- 时间段计算准确

---

#### 2.2.2 教练请假管理
**输入**: 请假时间段和原因
**输出**: 请假记录

**执行内容**:

**A. 后端API实现**
```go
// services/coach_leave_service.go

// 核心方法：
- ApplyLeave(coachID int64, req *ApplyLeaveRequest) error
- ApproveLeave(leaveID int64, approverID int64) error
- RejectLeave(leaveID int64, approverID int64, reason string) error
- CancelLeave(leaveID int64) error
- GetCoachLeaves(coachID int64, status int8) ([]*CoachLeave, error)
- CheckLeaveConflict(coachID int64, startTime, endTime time.Time) ([]*CourseBooking, error) // 检查请假时间段内是否有已预约课程

// API路由：
POST   /api/v1/coaches/:id/leaves           // 申请请假
PUT    /api/v1/leaves/:id/approve           // 批准请假
PUT    /api/v1/leaves/:id/reject            // 拒绝请假
DELETE /api/v1/leaves/:id                   // 取消请假
GET    /api/v1/coaches/:id/leaves           // 获取请假记录
GET    /api/v1/leaves/:id/conflicts         // 检查请假冲突
```

**B. 业务逻辑**
```
1. 申请请假流程：
   - 验证请假时间段
   - 检查是否有已预约的课程
   - 如有冲突，提示需要先处理预约
   - 创建请假记录（待审核状态）
   - 发送通知给管理员
   - 返回成功

2. 批准请假流程：
   - 更新请假状态为已批准
   - 自动取消请假时间段内的待确认预约
   - 发送通知给受影响的会员
   - 返回成功

3. 拒绝请假流程：
   - 更新请假状态为已拒绝
   - 记录拒绝原因
   - 发送通知给教练
   - 返回成功
```

**验收标准**:
- 请假冲突检测准确
- 自动处理受影响的预约
- 通知发送及时
- 支持请假记录查询和统计

---

### 2.3 课程预约功能

#### 2.3.1 会员预约课程
**输入**: 会员ID、教练ID、课程类型、预约时间
**输出**: 预约记录

**执行内容**:

**A. 后端API实现**
```go
// services/booking_service.go

// 核心方法：
- CreateBooking(req *CreateBookingRequest) (*CourseBooking, error)
- CancelBooking(bookingID int64, cancelBy int64, reason string) error
- ConfirmBooking(bookingID int64, coachID int64) error
- CheckInBooking(bookingID int64) error
- CheckOutBooking(bookingID int64) error
- RateBooking(bookingID int64, req *RateBookingRequest) error
- GetBookingByID(id int64) (*CourseBooking, error)
- ListUserBookings(userID int64, req *ListBookingsRequest) ([]*CourseBooking, int64, error)
- ListCoachBookings(coachID int64, req *ListBookingsRequest) ([]*CourseBooking, int64, error)
- GenerateBookingNo() string
- CheckBookingConflict(coachID int64, date time.Time, startTime, endTime string) bool
- DeductMembershipCardCourse(membershipCardID int64) error // 扣除会员卡课时

// API路由：
POST   /api/v1/bookings                      // 创建预约
PUT    /api/v1/bookings/:id/cancel           // 取消预约
PUT    /api/v1/bookings/:id/confirm          // 确认预约（教练）
PUT    /api/v1/bookings/:id/check-in         // 签到
PUT    /api/v1/bookings/:id/check-out        // 签出
PUT    /api/v1/bookings/:id/rate             // 评价
GET    /api/v1/bookings/:id                  // 获取预约详情
GET    /api/v1/users/:id/bookings            // 获取会员预约列表
GET    /api/v1/coaches/:id/bookings          // 获取教练预约列表
GET    /api/v1/bookings/calendar             // 获取日历视图数据
```

**B. 请求/响应结构**
```go
type CreateBookingRequest struct {
    UserID           int64  `json:"user_id" binding:"required"`
    CoachID          int64  `json:"coach_id" binding:"required"`
    CourseTypeID     int64  `json:"course_type_id" binding:"required"`
    MembershipCardID *int64 `json:"membership_card_id"` // 可选，使用会员卡课时
    BookingDate      string `json:"booking_date" binding:"required"` // YYYY-MM-DD
    StartTime        string `json:"start_time" binding:"required"`   // HH:MM
    Remark           string `json:"remark"`
}

type RateBookingRequest struct {
    Rating  int8   `json:"rating" binding:"required,min=1,max=5"`
    Comment string `json:"comment" binding:"required,max=500"`
}

type ListBookingsRequest struct {
    Page      int    `form:"page" binding:"required,min=1"`
    PageSize  int    `form:"page_size" binding:"required,min=1,max=100"`
    Status    int8   `form:"status"`
    StartDate string `form:"start_date"`
    EndDate   string `form:"end_date"`
}
```

**C. 业务逻辑**
```
1. 创建预约流程：
   - 验证会员是否存在且状态正常
   - 验证教练是否存在且状态正常
   - 验证课程类型是否存在
   - 验证预约时间是否在教练可约时间内
   - 检查时间冲突（教练该时间段是否已被预约）
   - 检查教练是否请假
   - 如使用会员卡，验证会员卡是否有效且有剩余课时
   - 生成预约编号
   - 创建预约记录（待确认状态）
   - 如使用会员卡，扣除课时
   - 发送通知给教练
   - 返回预约信息

2. 取消预约流程：
   - 验证预约是否存在
   - 验证预约状态（只能取消待确认和已确认的预约）
   - 检查取消时间限制（如：开课前2小时不能取消）
   - 更新预约状态为已取消
   - 记录取消原因和取消人
   - 如已扣除课时，退还课时
   - 发送通知给相关人员
   - 返回成功

3. 确认预约流程（教练）：
   - 验证预约是否存在
   - 验证操作人是否为该课程教练
   - 验证预约状态为待确认
   - 更新预约状态为已确认
   - 发送通知给会员
   - 返回成功

4. 签到流程：
   - 验证预约是否存在
   - 验证预约状态为已确认
   - 验证当前时间是否在课程时间范围内（允许提前15分钟签到）
   - 记录签到时间
   - 返回成功

5. 签出流程：
   - 验证预约是否存在
   - 验证是否已签到
   - 记录签出时间
   - 计算实际课程时长
   - 更新预约状态为已完成
   - 返回成功

6. 评价流程：
   - 验证预约是否存在
   - 验证预约状态为已完成
   - 验证是否已评价
   - 保存评分和评价内容
   - 记录评价时间
   - 更新教练评分统计
   - 返回成功
```

**验收标准**:
- 预约冲突检测准确
- 会员卡课时扣除和退还正确
- 取消时间限制生效
- 通知发送及时
- 预约状态流转正确
- 接口响应时间 < 300ms

---

#### 2.3.2 前端页面开发 - 预约课程
**输入**: 后端API接口文档
**输出**: 预约课程页面

**执行内容**:

**A. 教练选择页面 (src/pages/Booking/CoachSelection.tsx)**
```typescript
// 页面布局：
1. 课程类型筛选
   - 横向滚动的课程类型卡片
   - 显示课程名称、时长、价格

2. 教练列表
   - 教练卡片：头像、姓名、评分、擅长课程、简介
   - 筛选条件：性别、评分、擅长课程
   - 排序：评分、预约量

3. 教练详情弹窗
   - 教练基本信息
   - 擅长课程
   - 评分和评价列表
   - 可约时间展示
   - 预约按钮
```

**B. 时间选择页面 (src/pages/Booking/TimeSelection.tsx)**
```typescript
// 页面布局：
1. 日期选择器
   - 横向滚动的日期列表（显示未来7天）
   - 标注今天、明天
   - 不可选择已过期日期

2. 时间段选择
   - 网格布局展示可约时间段
   - 已预约时间段置灰不可选
   - 显示时间段状态（可约、已约、休息）

3. 课程信息确认
   - 教练信息
   - 课程类型
   - 预约时间
   - 使用会员卡选择（如有）
   - 备注输入
   - 确认预约按钮
```

**C. 预约列表页面 (src/pages/Booking/BookingList.tsx)**
```typescript
// 页面布局：
1. Tab切换
   - 待确认
   - 已确认
   - 已完成
   - 已取消

2. 预约卡片
   - 预约编号
   - 教练信息（头像、姓名）
   - 课程类型
   - 预约时间
   - 状态标签
   - 操作按钮（取消、签到、评价）

3. 预约详情弹窗
   - 完整预约信息
   - 操作按钮
   - 取消原因（如已取消）
   - 评价内容（如已评价）
```

**D. 日历视图页面 (src/pages/Booking/BookingCalendar.tsx)**
```typescript
// 页面布局：
1. 日历组件
   - 月视图
   - 日期上显示预约数量标记
   - 点击日期查看当天预约

2. 日程列表
   - 时间轴展示
   - 预约卡片（教练、课程、时间）
   - 支持拖拽调整时间（管理员）

3. 筛选条件
   - 教练筛选
   - 课程类型筛选
   - 状态筛选
```

**验收标准**:
- 页面交互流畅
- 时间选择直观易用
- 预约状态实时更新
- 支持快速预约
- 日历视图清晰展示课程安排

---

### 2.4 课程评价与统计

#### 2.4.1 课程评价功能
**输入**: 预约ID、评分、评价内容
**输出**: 评价记录

**执行内容**:

**A. 评价展示页面**
```typescript
// src/components/Booking/RatingForm.tsx

// 功能点：
1. 评分组件
   - 5星评分
   - 支持半星
   - 实时预览

2. 评价内容
   - 文本域输入
   - 字数限制（500字）
   - 字数统计

3. 标签选择（可选）
   - 预设标签：专业、耐心、准时、效果好等
   - 支持多选

4. 提交按钮
   - 验证评分和内容
   - 提交后不可修改
```

**B. 评价列表展示**
```typescript
// src/components/Coach/RatingList.tsx

// 功能点：
1. 评价统计
   - 平均评分（大数字展示）
   - 各星级数量分布（柱状图）
   - 总评价数

2. 评价列表
   - 会员头像、昵名
   - 评分星级
   - 评价内容
   - 评价时间
   - 课程类型

3. 筛选和排序
   - 按评分筛选
   - 按时间排序
   - 按课程类型筛选
```

**验收标准**:
- 评价提交成功
- 评分统计准确
- 评价展示完整
- 支持分页加载

---

#### 2.4.2 课程统计报表
**输入**: 时间范围、教练ID
**输出**: 统计报表

**执行内容**:

**A. 后端统计API**
```go
// services/booking_stats_service.go

// 核心方法：
- GetCoachBookingStats(coachID int64, startDate, endDate time.Time) (*CoachStats, error)
- GetCourseTypeStats(startDate, endDate time.Time) ([]*CourseTypeStats, error)
- GetBookingTrend(startDate, endDate time.Time) ([]*TrendData, error)
- GetPeakHoursStats(startDate, endDate time.Time) ([]*PeakHourData, error)
- GetCancellationRate(startDate, endDate time.Time) (float64, error)

// API路由：
GET /api/v1/stats/coaches/:id/bookings    // 教练预约统计
GET /api/v1/stats/course-types            // 课程类型统计
GET /api/v1/stats/bookings/trend          // 预约趋势
GET /api/v1/stats/bookings/peak-hours     // 高峰时段统计
GET /api/v1/stats/bookings/cancellation   // 取消率统计
```

**B. 前端统计页面**
```typescript
// src/pages/Stats/BookingStats.tsx

// 展示内容：
1. 统计卡片
   - 总预约数
   - 已完成数
   - 取消率
   - 平均评分

2. 预约趋势图（折线图）
   - X轴：日期
   - Y轴：预约数量
   - 多条线：总预约、已完成、已取消

3. 课程类型分布（饼图）
   - 各课程类型预约占比
   - 点击查看详情

4. 高峰时段分析（热力图）
   - X轴：时间段
   - Y轴：星期几
   - 颜色深度表示预约密度

5. 教练排行榜
   - 预约量排行
   - 评分排行
   - 完成率排行

6. 导出功能
   - 导出统计报表（Excel）
   - 导出图表（图片）
```

**验收标准**:
- 统计数据准确
- 图表渲染流畅
- 支持时间范围选择
- 支持导出功能

---

## 三、接口文档

### 3.1 API列表汇总

| 接口路径 | 方法 | 说明 | 权限 |
|---------|------|------|------|
| /api/v1/course-types | POST | 创建课程类型 | 管理员 |
| /api/v1/course-types/:id | PUT | 更新课程类型 | 管理员 |
| /api/v1/course-types/:id | DELETE | 删除课程类型 | 管理员 |
| /api/v1/course-types | GET | 课程类型列表 | 所有角色 |
| /api/v1/coaches/:id/courses | POST | 为教练分配课程 | 管理员 |
| /api/v1/coaches/:id/courses | GET | 获取教练课程列表 | 所有角色 |
| /api/v1/coaches/:id/available-times | POST | 设置可约时间 | 教练、管理员 |
| /api/v1/coaches/:id/available-times | GET | 获取可约时间 | 所有角色 |
| /api/v1/coaches/:id/available-slots | GET | 获取可约时间段 | 所有角色 |
| /api/v1/coaches/:id/leaves | POST | 申请请假 | 教练 |
| /api/v1/leaves/:id/approve | PUT | 批准请假 | 管理员 |
| /api/v1/bookings | POST | 创建预约 | 会员、前台 |
| /api/v1/bookings/:id/cancel | PUT | 取消预约 | 会员、教练、管理员 |
| /api/v1/bookings/:id/confirm | PUT | 确认预约 | 教练 |
| /api/v1/bookings/:id/check-in | PUT | 签到 | 教练、前台 |
| /api/v1/bookings/:id/check-out | PUT | 签出 | 教练、前台 |
| /api/v1/bookings/:id/rate | PUT | 评价 | 会员 |
| /api/v1/users/:id/bookings | GET | 获取会员预约列表 | 会员、管理员 |
| /api/v1/coaches/:id/bookings | GET | 获取教练预约列表 | 教练、管理员 |
| /api/v1/bookings/calendar | GET | 获取日历视图数据 | 所有角色 |

### 3.2 核心接口详细说明

#### 创建预约
```
POST /api/v1/bookings

Request Body:
{
  "user_id": 1,
  "coach_id": 2,
  "course_type_id": 1,
  "membership_card_id": 1,
  "booking_date": "2024-01-15",
  "start_time": "10:00",
  "remark": "第一次上课"
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "booking_no": "BK202401150001",
    "user_id": 1,
    "coach_id": 2,
    "course_type_id": 1,
    "booking_date": "2024-01-15",
    "start_time": "10:00",
    "end_time": "11:00",
    "duration": 60,
    "status": 1,
    "created_at": "2024-01-10T10:00:00Z"
  }
}
```

#### 获取可约时间段
```
GET /api/v1/coaches/2/available-slots?date=2024-01-15

Response:
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "start_time": "09:00",
      "end_time": "10:00",
      "is_booked": false,
      "booking_id": null
    },
    {
      "start_time": "10:00",
      "end_time": "11:00",
      "is_booked": true,
      "booking_id": 1
    },
    {
      "start_time": "11:00",
      "end_time": "12:00",
      "is_booked": false,
      "booking_id": null
    }
  ]
}
```

## 四、测试用例

### 4.1 单元测试
- 预约编号生成测试
- 时间冲突检测测试
- 会员卡课时扣除测试
- 连续预约限制测试
- 取消时间限制测试

### 4.2 接口测试
- 所有API接口的正常流程测试
- 预约冲突场景测试
- 并发预约测试
- 权限控制测试
- 异常情况处理测试

### 4.3 前端测试
- 时间选择交互测试
- 预约流程测试
- 日历视图测试
- 评价功能测试

## 五、上线检查清单

### 5.1 后端检查
- [ ] 数据库表创建完成
- [ ] 所有API接口开发完成
- [ ] 时间冲突检测逻辑正确
- [ ] 会员卡课时扣除逻辑正确
- [ ] 通知发送功能正常
- [ ] 单元测试通过率 > 80%
- [ ] 接口测试全部通过
- [ ] 并发测试通过

### 5.2 前端检查
- [ ] 所有页面开发完成
- [ ] 预约流程流畅
- [ ] 日历视图展示正确
- [ ] 时间选择交互友好
- [ ] 评价功能正常
- [ ] 页面响应式适配
- [ ] 浏览器兼容性测试

### 5.3 联调检查
- [ ] 前后端接口联调完成
- [ ] 预约流程端到端测试通过
- [ ] 通知推送正常
- [ ] 性能测试通过

## 六、后续优化方向

1. **智能推荐**: 根据会员训练目标和历史数据推荐合适的教练和课程
2. **自动排课**: 基于教练和会员的时间偏好自动推荐最佳预约时间
3. **课程包管理**: 支持课程包购买和使用
4. **候补预约**: 当时间段已满时支持候补，有取消时自动通知
5. **视频课程**: 支持线上视频课程预约和上课
6. **课程回放**: 录制课程视频供会员回看

---
**任务优先级**: P0（核心功能）  
**预计工期**: 3-4周  
**依赖任务**: TASK001（用户管理）、TASK006（教练管理）  
**后续任务**: TASK009（数据统计与报表）
