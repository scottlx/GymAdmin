# TASK002 - 会员卡管理功能

## 一、功能概述

### 1.1 功能目标
实现会员卡的全生命周期管理，包括会员卡类型配置、开卡、续费、转卡、冻结、到期提醒等功能。

### 1.2 核心价值
- 规范会员卡管理流程，提高运营效率
- 自动化到期提醒，提升续费转化率
- 支持灵活的会员卡类型配置，满足不同业务场景
- 完整的会员卡操作记录，便于审计和追溯

### 1.3 涉及角色
- **管理员**: 可以配置会员卡类型、管理所有会员卡
- **前台人员**: 可以办理开卡、续费、转卡等业务
- **会员**: 可以查看自己的会员卡信息

## 二、功能详细拆解

### 2.1 会员卡类型管理

#### 2.1.1 数据库设计
**输入**: 业务需求分析
**输出**: 数据库表结构SQL文件

**执行内容**:
```sql
-- 会员卡类型表 (card_types)
CREATE TABLE card_types (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '卡类型ID',
    type_name VARCHAR(50) NOT NULL COMMENT '卡类型名称',
    type_code VARCHAR(20) UNIQUE NOT NULL COMMENT '卡类型编码',
    duration_type TINYINT NOT NULL COMMENT '时长类型：1-天卡，2-月卡，3-季卡，4-年卡，5-次卡',
    duration_value INT NOT NULL COMMENT '时长数值（天数或次数）',
    price DECIMAL(10,2) NOT NULL COMMENT '价格',
    original_price DECIMAL(10,2) COMMENT '原价',
    description TEXT COMMENT '卡类型描述',
    benefits TEXT COMMENT '权益说明（JSON格式）',
    can_freeze TINYINT DEFAULT 1 COMMENT '是否可冻结：0-否，1-是',
    max_freeze_times INT DEFAULT 0 COMMENT '最大冻结次数（0表示不限）',
    max_freeze_days INT DEFAULT 0 COMMENT '最大冻结天数（0表示不限）',
    can_transfer TINYINT DEFAULT 1 COMMENT '是否可转卡：0-否，1-是',
    transfer_fee DECIMAL(10,2) DEFAULT 0 COMMENT '转卡手续费',
    status TINYINT DEFAULT 1 COMMENT '状态：1-启用，2-停用',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_status (status),
    INDEX idx_sort_order (sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员卡类型表';

-- 会员卡表 (membership_cards)
CREATE TABLE membership_cards (
    id BIGINT PRIMARY KEY AUTO_INCREMENT COMMENT '会员卡ID',
    card_no VARCHAR(32) UNIQUE NOT NULL COMMENT '卡号（自动生成）',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    card_type_id BIGINT NOT NULL COMMENT '卡类型ID',
    status TINYINT DEFAULT 1 COMMENT '状态：1-正常，2-已过期，3-已冻结，4-已转出，5-已退卡',
    start_date DATE NOT NULL COMMENT '开始日期',
    end_date DATE NOT NULL COMMENT '结束日期',
    remaining_times INT COMMENT '剩余次数（次卡）',
    total_times INT COMMENT '总次数（次卡）',
    freeze_times INT DEFAULT 0 COMMENT '已冻结次数',
    freeze_days INT DEFAULT 0 COMMENT '已冻结天数',
    is_frozen TINYINT DEFAULT 0 COMMENT '是否冻结中：0-否，1-是',
    frozen_at TIMESTAMP NULL COMMENT '冻结时间',
    source TINYINT DEFAULT 1 COMMENT '来源：1-前台办理，2-小程序购买，3-美团，4-抖音',
    purchase_price DECIMAL(10,2) NOT NULL COMMENT '购买价格',
    operator_id BIGINT COMMENT '操作员ID',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_card_no (card_no),
    INDEX idx_status (status),
    INDEX idx_end_date (end_date),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (card_type_id) REFERENCES card_types(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员卡表';

-- 会员卡操作记录表 (card_operation_logs)
CREATE TABLE card_operation_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    card_id BIGINT NOT NULL COMMENT '会员卡ID',
    operation_type TINYINT NOT NULL COMMENT '操作类型：1-开卡，2-续费，3-冻结，4-解冻，5-转卡，6-退卡',
    old_end_date DATE COMMENT '原结束日期',
    new_end_date DATE COMMENT '新结束日期',
    amount DECIMAL(10,2) COMMENT '金额',
    operator_id BIGINT NOT NULL COMMENT '操作员ID',
    remark VARCHAR(255) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_card_id (card_id),
    INDEX idx_operation_type (operation_type),
    FOREIGN KEY (card_id) REFERENCES membership_cards(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员卡操作记录表';

-- 会员卡冻结记录表 (card_freeze_records)
CREATE TABLE card_freeze_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    card_id BIGINT NOT NULL COMMENT '会员卡ID',
    freeze_start_date DATE NOT NULL COMMENT '冻结开始日期',
    freeze_end_date DATE COMMENT '冻结结束日期',
    freeze_days INT NOT NULL COMMENT '冻结天数',
    reason VARCHAR(255) COMMENT '冻结原因',
    operator_id BIGINT NOT NULL COMMENT '操作员ID',
    status TINYINT DEFAULT 1 COMMENT '状态：1-冻结中，2-已解冻',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_card_id (card_id),
    INDEX idx_status (status),
    FOREIGN KEY (card_id) REFERENCES membership_cards(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员卡冻结记录表';

-- 会员卡转让记录表 (card_transfer_records)
CREATE TABLE card_transfer_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    card_id BIGINT NOT NULL COMMENT '会员卡ID',
    from_user_id BIGINT NOT NULL COMMENT '转出用户ID',
    to_user_id BIGINT NOT NULL COMMENT '转入用户ID',
    transfer_fee DECIMAL(10,2) DEFAULT 0 COMMENT '转卡手续费',
    operator_id BIGINT NOT NULL COMMENT '操作员ID',
    remark VARCHAR(255) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_card_id (card_id),
    INDEX idx_from_user_id (from_user_id),
    INDEX idx_to_user_id (to_user_id),
    FOREIGN KEY (card_id) REFERENCES membership_cards(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员卡转让记录表';
```

**验收标准**:
- 表结构设计合理，支持各种会员卡业务场景
- 包含必要的索引优化查询性能
- 外键约束保证数据一致性
- 操作记录完整，便于审计

---

#### 2.1.2 会员卡类型配置功能

**输入**: 数据库表结构
**输出**: 会员卡类型管理API和前端页面

**执行内容**:

**A. 后端API实现**
```go
// services/card_type_service.go

// 核心方法：
- CreateCardType(req *CreateCardTypeRequest) (*CardType, error)
- UpdateCardType(id int64, req *UpdateCardTypeRequest) error
- DeleteCardType(id int64) error
- GetCardTypeByID(id int64) (*CardType, error)
- ListCardTypes(req *ListCardTypesRequest) ([]*CardType, int64, error)
- EnableCardType(id int64) error
- DisableCardType(id int64) error

// API路由：
POST   /api/v1/card-types           // 创建卡类型
PUT    /api/v1/card-types/:id       // 更新卡类型
DELETE /api/v1/card-types/:id       // 删除卡类型
GET    /api/v1/card-types/:id       // 获取卡类型详情
GET    /api/v1/card-types           // 卡类型列表
POST   /api/v1/card-types/:id/enable  // 启用卡类型
POST   /api/v1/card-types/:id/disable // 停用卡类型
```

**B. 前端页面实现**
```typescript
// src/pages/CardType/CardTypeList.tsx

// 功能点：
1. 卡类型列表展示
   - 表格列：类型名称、类型编码、时长、价格、状态、操作
   - 支持拖拽排序
   - 状态标识（启用/停用）

2. 操作功能
   - 新增卡类型：打开表单弹窗
   - 编辑：打开编辑弹窗
   - 启用/停用：切换状态
   - 删除：二次确认（检查是否有关联会员卡）

3. 卡类型表单
   - 基本信息：类型名称、类型编码、时长类型、时长数值
   - 价格信息：价格、原价
   - 权益配置：描述、权益说明
   - 规则配置：是否可冻结、最大冻结次数/天数、是否可转卡、转卡手续费
```

**验收标准**:
- 卡类型配置灵活，满足各种业务场景
- 表单校验完整，防止配置错误
- 启用/停用状态切换正常
- 删除前检查是否有关联数据
- 支持拖拽排序，调整展示顺序

---

### 2.2 会员卡办理与续费

#### 2.2.1 开卡功能

**输入**: 用户信息、卡类型信息
**输出**: 会员卡记录

**执行内容**:

**A. 后端API实现**
```go
// services/membership_card_service.go

// 核心方法：
- CreateCard(req *CreateCardRequest) (*MembershipCard, error)
- GenerateCardNo() string // 生成卡号
- CalculateEndDate(startDate time.Time, cardType *CardType) time.Time
- ValidateUserCanCreateCard(userID int64) error // 验证用户是否可以开卡

// API路由：
POST /api/v1/membership-cards        // 开卡
```

**B. 开卡业务逻辑**
```
1. 验证用户是否存在且状态正常
2. 验证卡类型是否存在且已启用
3. 生成唯一卡号
4. 计算开始日期和结束日期
   - 时间卡：结束日期 = 开始日期 + 时长
   - 次卡：结束日期 = 开始日期 + 1年（默认有效期）
5. 创建会员卡记录
6. 创建操作日志
7. 创建用户训练统计记录（如果不存在）
8. 返回会员卡信息
```

**C. 前端开卡界面**
```typescript
// src/components/Card/CreateCardForm.tsx

// 表单字段：
1. 用户选择
   - 搜索用户（姓名、手机号）
   - 显示用户基本信息

2. 卡类型选择
   - 卡片式展示所有启用的卡类型
   - 显示价格、时长、权益

3. 开卡信息
   - 开始日期（默认今天，可修改）
   - 购买价格（默认卡类型价格，可修改）
   - 备注

4. 支付信息（可选）
   - 支付方式：现金、微信、支付宝、银行卡
   - 支付金额
```

**验收标准**:
- 卡号自动生成且唯一
- 结束日期计算准确
- 开卡成功后立即生效
- 操作日志记录完整
- 支持批量开卡（可选）

---

#### 2.2.2 续费功能

**输入**: 会员卡信息、续费时长
**输出**: 更新后的会员卡记录

**执行内容**:

**A. 后端API实现**
```go
// services/membership_card_service.go

// 核心方法：
- RenewCard(cardID int64, req *RenewCardRequest) error
- CalculateRenewEndDate(card *MembershipCard, renewType *CardType) time.Time

// API路由：
POST /api/v1/membership-cards/:id/renew  // 续费
```

**B. 续费业务逻辑**
```
1. 验证会员卡是否存在
2. 验证会员卡状态（已过期、已转出、已退卡不能续费）
3. 验证续费卡类型
4. 计算新的结束日期
   - 如果当前未过期：新结束日期 = 当前结束日期 + 续费时长
   - 如果已过期：新结束日期 = 今天 + 续费时长
5. 更新会员卡信息
6. 如果卡状态是已过期，更新为正常
7. 创建操作日志
8. 返回更新后的会员卡信息
```

**C. 前端续费界面**
```typescript
// src/components/Card/RenewCardForm.tsx

// 展示内容：
1. 当前会员卡信息
   - 卡号、卡类型、当前结束日期、剩余天数/次数

2. 续费选项
   - 选择续费卡类型（可以与原卡类型不同）
   - 显示续费后的新结束日期
   - 续费价格（可修改）

3. 支付信息
   - 支付方式
   - 支付金额

4. 备注
```

**验收标准**:
- 续费后结束日期计算准确
- 支持不同卡类型续费
- 已过期卡续费后状态更新为正常
- 操作日志记录完整
- 续费成功后立即生效

---

### 2.3 会员卡冻结与解冻

#### 2.3.1 冻结功能

**输入**: 会员卡信息、冻结天数、冻结原因
**输出**: 冻结记录、更新后的会员卡

**执行内容**:

**A. 后端API实现**
```go
// services/card_freeze_service.go

// 核心方法：
- FreezeCard(cardID int64, req *FreezeCardRequest) error
- ValidateCanFreeze(card *MembershipCard, freezeDays int) error
- CalculateNewEndDate(card *MembershipCard, freezeDays int) time.Time

// API路由：
POST /api/v1/membership-cards/:id/freeze    // 冻结
POST /api/v1/membership-cards/:id/unfreeze  // 解冻
GET  /api/v1/membership-cards/:id/freeze-records // 冻结记录
```

**B. 冻结业务逻辑**
```
1. 验证会员卡是否存在且状态正常
2. 验证卡类型是否允许冻结
3. 验证冻结次数是否超限
4. 验证冻结天数是否超限
5. 创建冻结记录
6. 更新会员卡状态为已冻结
7. 更新会员卡的冻结次数和冻结天数
8. 延长会员卡结束日期（结束日期 + 冻结天数）
9. 创建操作日志
10. 返回冻结成功
```

**C. 解冻业务逻辑**
```
1. 验证会员卡是否存在且状态为已冻结
2. 查找当前冻结记录
3. 计算实际冻结天数
4. 更新冻结记录状态为已解冻
5. 更新会员卡状态为正常
6. 如果实际冻结天数 < 申请冻结天数，调整结束日期
7. 创建操作日志
8. 返回解冻成功
```

**D. 前端冻结界面**
```typescript
// src/components/Card/FreezeCardForm.tsx

// 表单字段：
1. 当前会员卡信息展示
   - 卡号、卡类型、结束日期
   - 已冻结次数/最大冻结次数
   - 已冻结天数/最大冻结天数

2. 冻结信息
   - 冻结天数（输入框，校验不超过最大限制）
   - 冻结原因（下拉选择 + 自定义输入）
   - 预计解冻日期（自动计算）
   - 冻结后的新结束日期（自动计算）

3. 冻结记录列表
   - 历史冻结记录
   - 显示冻结时间、解冻时间、冻结天数、原因
```

**验收标准**:
- 冻结次数和天数限制校验准确
- 冻结后结束日期自动延长
- 解冻后根据实际冻结天数调整结束日期
- 冻结期间不能签到
- 操作日志记录完整

---

### 2.4 会员卡转让

#### 2.4.1 转卡功能

**输入**: 会员卡信息、转入用户信息
**输出**: 转让记录、更新后的会员卡

**执行内容**:

**A. 后端API实现**
```go
// services/card_transfer_service.go

// 核心方法：
- TransferCard(cardID int64, req *TransferCardRequest) error
- ValidateCanTransfer(card *MembershipCard) error
- ValidateTargetUser(userID int64) error

// API路由：
POST /api/v1/membership-cards/:id/transfer  // 转卡
GET  /api/v1/membership-cards/:id/transfer-records // 转让记录
```

**B. 转卡业务逻辑**
```
1. 验证会员卡是否存在且状态正常
2. 验证卡类型是否允许转卡
3. 验证转入用户是否存在且状态正常
4. 验证转入用户不是当前持卡人
5. 计算转卡手续费
6. 创建转让记录
7. 更新会员卡的用户ID为转入用户
8. 创建操作日志
9. 发送通知给转入和转出用户
10. 返回转卡成功
```

**C. 前端转卡界面**
```typescript
// src/components/Card/TransferCardForm.tsx

// 表单字段：
1. 当前会员卡信息
   - 卡号、卡类型、当前持卡人、结束日期、剩余天数/次数

2. 转入用户选择
   - 搜索用户（姓名、手机号）
   - 显示用户基本信息
   - 验证用户状态

3. 转卡信息
   - 转卡手续费（自动计算，可修改）
   - 转卡原因
   - 备注

4. 确认信息
   - 显示转出用户和转入用户信息
   - 二次确认弹窗
```

**验收标准**:
- 转卡权限校验准确
- 转卡后会员卡归属立即变更
- 转卡手续费计算正确
- 转让记录完整
- 转出和转入用户都收到通知
- 操作日志记录完整

---

### 2.5 会员卡到期管理

#### 2.5.1 到期提醒功能

**输入**: 会员卡数据
**输出**: 到期提醒消息

**执行内容**:

**A. 后端定时任务实现**
```go
// services/card_expiry_service.go

// 核心方法：
- CheckExpiringCards() error // 检查即将到期的卡
- SendExpiryNotification(card *MembershipCard, daysLeft int) error
- UpdateExpiredCards() error // 更新已过期的卡状态

// 定时任务：
- 每天凌晨1点执行：检查即将到期的卡（7天、3天、1天）
- 每天凌晨2点执行：更新已过期的卡状态
```

**B. 提醒规则配置**
```go
// 提醒时间点配置
type ExpiryReminderConfig struct {
    Days []int // 提前多少天提醒，如：[7, 3, 1]
    Channels []string // 提醒渠道：sms, wechat, app_push
}

// 提醒内容模板
- 7天提醒：您的会员卡将在7天后到期，请及时续费
- 3天提醒：您的会员卡将在3天后到期，续费享优惠
- 1天提醒：您的会员卡明天到期，请尽快续费
- 当天提醒：您的会员卡今天到期，立即续费
```

**C. 前端到期管理界面**
```typescript
// src/pages/Card/ExpiryManagement.tsx

// 功能点：
1. 即将到期列表
   - 筛选：7天内、3天内、今天到期
   - 表格列：用户姓名、手机号、卡号、卡类型、到期日期、剩余天数
   - 操作：续费、发送提醒

2. 已过期列表
   - 筛选：过期时间范围
   - 表格列：用户姓名、手机号、卡号、卡类型、过期日期、过期天数
   - 操作：续费、联系用户

3. 统计数据
   - 今日到期数量
   - 7天内到期数量
   - 本月到期数量
   - 已过期未续费数量

4. 批量操作
   - 批量发送提醒
   - 导出到期名单
```

**验收标准**:
- 定时任务准时执行
- 提醒消息准确发送
- 已过期卡状态自动更新
- 提醒记录可追溯
- 支持手动发送提醒

---

#### 2.5.2 自动续费功能（可选）

**输入**: 用户授权、会员卡信息
**输出**: 自动续费记录

**执行内容**:

**A. 自动续费配置**
```sql
-- 自动续费配置表
CREATE TABLE auto_renew_configs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    card_type_id BIGINT NOT NULL COMMENT '卡类型ID',
    is_enabled TINYINT DEFAULT 1 COMMENT '是否启用',
    payment_method TINYINT NOT NULL COMMENT '支付方式：1-微信，2-支付宝',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_card_type (user_id, card_type_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='自动续费配置表';
```

**B. 自动续费逻辑**
```
1. 定时任务：每天检查即将到期的卡（提前1天）
2. 查询是否开启自动续费
3. 调用支付接口扣款
4. 扣款成功后自动续费
5. 发送续费成功通知
6. 扣款失败则发送提醒通知
```

**验收标准**:
- 自动续费配置灵活
- 扣款失败有重试机制
- 续费成功后发送通知
- 用户可随时取消自动续费

---

### 2.6 会员卡查询与统计

#### 2.6.1 会员卡列表与查询

**输入**: 查询条件
**输出**: 会员卡列表

**执行内容**:

**A. 后端API实现**
```go
// services/membership_card_service.go

// 核心方法：
- ListCards(req *ListCardsRequest) ([]*MembershipCard, int64, error)
- GetCardByID(id int64) (*MembershipCard, error)
- GetCardByCardNo(cardNo string) (*MembershipCard, error)
- GetUserCards(userID int64) ([]*MembershipCard, error)
- GetCardOperationLogs(cardID int64) ([]*CardOperationLog, error)

// API路由：
GET /api/v1/membership-cards              // 会员卡列表
GET /api/v1/membership-cards/:id          // 会员卡详情
GET /api/v1/membership-cards/no/:card_no  // 根据卡号查询
GET /api/v1/users/:user_id/cards          // 用户的会员卡列表
GET /api/v1/membership-cards/:id/logs     // 操作日志
```

**B. 查询条件**
```go
type ListCardsRequest struct {
    Page       int    `form:"page"`
    PageSize   int    `form:"page_size"`
    Keyword    string `form:"keyword"`     // 卡号、用户姓名、手机号
    CardTypeID int64  `form:"card_type_id"` // 卡类型
    Status     int8   `form:"status"`      // 状态
    Source     int8   `form:"source"`      // 来源
    StartDate  string `form:"start_date"`  // 开始日期范围
    EndDate    string `form:"end_date"`    // 结束日期范围
    ExpiryDays int    `form:"expiry_days"` // 即将到期天数
}
```

**C. 前端列表页面**
```typescript
// src/pages/Card/CardList.tsx

// 功能点：
1. 会员卡列表
   - 表格列：卡号、用户姓名、手机号、卡类型、状态、开始日期、结束日期、剩余天数/次数、操作
   - 状态标识：正常（绿色）、已过期（红色）、已冻结（橙色）
   - 支持排序

2. 搜索与筛选
   - 关键词搜索
   - 卡类型筛选
   - 状态筛选
   - 来源筛选
   - 日期范围筛选
   - 即将到期筛选

3. 操作按钮
   - 查看详情
   - 续费
   - 冻结/解冻
   - 转卡
   - 退卡

4. 批量操作
   - 批量导出
   - 批量发送到期提醒
```

**验收标准**:
- 列表查询性能良好（< 500ms）
- 筛选条件组合灵活
- 分页加载流畅
- 操作按钮权限控制正确

---

#### 2.6.2 会员卡统计报表

**输入**: 统计时间范围
**输出**: 统计报表数据

**执行内容**:

**A. 后端统计API**
```go
// services/card_stats_service.go

// 核心方法：
- GetCardStatsSummary() (*CardStatsSummary, error) // 总体统计
- GetCardTypeDistribution() ([]CardTypeStats, error) // 卡类型分布
- GetCardSalesTrend(days int) ([]SalesData, error) // 销售趋势
- GetCardRenewalRate(month string) (*RenewalRate, error) // 续费率统计

// API路由：
GET /api/v1/cards/stats/summary       // 总体统计
GET /api/v1/cards/stats/distribution  // 卡类型分布
GET /api/v1/cards/stats/sales-trend   // 销售趋势
GET /api/v1/cards/stats/renewal-rate  // 续费率
```

**B. 统计指标**
```
1. 总体统计
   - 总会员卡数量
   - 正常卡数量
   - 已过期卡数量
   - 已冻结卡数量
   - 本月新增数量
   - 本月续费数量

2. 卡类型分布
   - 各卡类型的数量和占比
   - 各卡类型的销售额和占比

3. 销售趋势
   - 每日/每月新增会员卡数量
   - 每日/每月销售额

4. 续费率统计
   - 到期会员卡数量
   - 续费会员卡数量
   - 续费率 = 续费数量 / 到期数量
```

**C. 前端统计页面**
```typescript
// src/pages/Card/CardStats.tsx

// 展示内容：
1. 统计卡片
   - 总会员卡数、正常卡数、已过期卡数、已冻结卡数
   - 本月新增、本月续费、本月销售额

2. 卡类型分布图（饼图）
   - 各卡类型数量占比
   - 点击查看详情

3. 销售趋势图（折线图）
   - X轴：日期
   - Y轴：销售数量/销售额
   - 支持切换时间范围

4. 续费率统计（柱状图）
   - X轴：月份
   - Y轴：续费率
   - 显示到期数量和续费数量

5. 数据导出
   - 导出统计报表
   - 支持Excel和PDF格式
```

**验收标准**:
- 统计数据准确无误
- 图表展示清晰直观
- 支持时间范围切换
- 数据更新及时
- 报表可导出

---

## 三、接口文档

### 3.1 API列表汇总

| 接口路径 | 方法 | 说明 | 权限 |
|---------|------|------|------|
| /api/v1/card-types | POST | 创建卡类型 | 管理员 |
| /api/v1/card-types/:id | PUT | 更新卡类型 | 管理员 |
| /api/v1/card-types/:id | DELETE | 删除卡类型 | 管理员 |
| /api/v1/card-types | GET | 卡类型列表 | 所有角色 |
| /api/v1/membership-cards | POST | 开卡 | 管理员、前台 |
| /api/v1/membership-cards/:id/renew | POST | 续费 | 管理员、前台 |
| /api/v1/membership-cards/:id/freeze | POST | 冻结 | 管理员、前台 |
| /api/v1/membership-cards/:id/unfreeze | POST | 解冻 | 管理员、前台 |
| /api/v1/membership-cards/:id/transfer | POST | 转卡 | 管理员、前台 |
| /api/v1/membership-cards | GET | 会员卡列表 | 所有角色 |
| /api/v1/membership-cards/:id | GET | 会员卡详情 | 所有角色 |
| /api/v1/cards/stats/summary | GET | 统计数据 | 管理员 |

## 四、测试用例

### 4.1 功能测试
- 开卡流程测试
- 续费流程测试
- 冻结/解冻流程测试
- 转卡流程测试
- 到期提醒测试

### 4.2 边界测试
- 冻结次数/天数限制测试
- 结束日期计算准确性测试
- 并发开卡测试
- 卡号唯一性测试

### 4.3 性能测试
- 列表查询性能测试
- 统计报表性能测试
- 定时任务性能测试

## 五、上线检查清单

- [ ] 数据库表创建完成
- [ ] 所有API接口开发完成
- [ ] 前端页面开发完成
- [ ] 定时任务配置完成
- [ ] 到期提醒功能测试通过
- [ ] 操作日志记录完整
- [ ] 权限控制正确
- [ ] 性能测试通过

---
**任务优先级**: P0（核心功能）  
**预计工期**: 3-4周  
**依赖任务**: TASK001（用户管理功能）  
**后续任务**: TASK003（人脸识别入场功能）
