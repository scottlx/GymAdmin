# TASK004 - 美团抖音券核销功能

## 一、功能概述

### 1.1 功能目标
实现与美团、抖音等第三方平台的券核销对接,支持会员通过平台购买的券进行会员卡办理、续费、私教课购买等操作。

### 1.2 核心价值
- 打通线上线下渠道,扩大获客来源
- 提升用户购买便利性
- 增加营收渠道
- 自动化核销流程,减少人工成本

### 1.3 涉及角色
- **管理员**: 配置平台对接参数、查看核销数据
- **前台人员**: 执行券核销操作
- **会员**: 在平台购买券、到店核销

## 二、功能详细拆解

### 2.1 数据库设计

```sql
-- 第三方平台配置表
CREATE TABLE third_party_platforms (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    platform_name VARCHAR(50) NOT NULL COMMENT '平台名称:meituan,douyin',
    platform_code VARCHAR(20) UNIQUE NOT NULL COMMENT '平台编码',
    app_id VARCHAR(100) COMMENT '应用ID',
    app_secret VARCHAR(255) COMMENT '应用密钥',
    api_url VARCHAR(255) COMMENT 'API地址',
    status TINYINT DEFAULT 1 COMMENT '状态:1-启用,2-停用',
    config JSON COMMENT '其他配置参数',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='第三方平台配置表';

-- 券信息表
CREATE TABLE vouchers (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    voucher_no VARCHAR(50) UNIQUE NOT NULL COMMENT '券号',
    platform_id BIGINT NOT NULL COMMENT '平台ID',
    platform_order_no VARCHAR(100) COMMENT '平台订单号',
    voucher_type TINYINT NOT NULL COMMENT '券类型:1-会员卡,2-私教课,3-团课',
    product_id BIGINT COMMENT '关联产品ID',
    product_name VARCHAR(100) COMMENT '产品名称',
    original_price DECIMAL(10,2) COMMENT '原价',
    paid_price DECIMAL(10,2) COMMENT '实付价',
    user_phone VARCHAR(11) COMMENT '购买用户手机号',
    user_name VARCHAR(50) COMMENT '购买用户姓名',
    status TINYINT DEFAULT 1 COMMENT '状态:1-未核销,2-已核销,3-已过期,4-已退款',
    expire_date DATE COMMENT '过期日期',
    verified_at TIMESTAMP NULL COMMENT '核销时间',
    verified_by BIGINT COMMENT '核销人ID',
    user_id BIGINT COMMENT '核销后关联的用户ID',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_voucher_no (voucher_no),
    INDEX idx_platform_order_no (platform_order_no),
    INDEX idx_user_phone (user_phone),
    INDEX idx_status (status),
    FOREIGN KEY (platform_id) REFERENCES third_party_platforms(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='券信息表';

-- 券核销记录表
CREATE TABLE voucher_verification_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    voucher_id BIGINT NOT NULL COMMENT '券ID',
    operation_type TINYINT NOT NULL COMMENT '操作类型:1-核销,2-撤销核销',
    operator_id BIGINT NOT NULL COMMENT '操作员ID',
    result_type TINYINT COMMENT '核销结果类型:1-开卡,2-续费,3-购买课程',
    result_id BIGINT COMMENT '结果ID(会员卡ID或课程ID)',
    remark VARCHAR(255) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_voucher_id (voucher_id),
    FOREIGN KEY (voucher_id) REFERENCES vouchers(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='券核销记录表';
```

### 2.2 美团平台对接

#### 2.2.1 美团API集成
**输入**: 美团开放平台凭证
**输出**: 完整的美团API对接

**执行内容**:
- 注册美团开放平台账号
- 创建应用获取AppID和AppSecret
- 实现OAuth2.0授权流程
- 实现券查询接口
- 实现券核销接口
- 实现券撤销接口
- 实现订单查询接口

**API路由**:
```
POST /api/v1/meituan/auth          // 美团授权
GET  /api/v1/meituan/voucher/:code // 查询券信息
POST /api/v1/meituan/verify        // 核销券
POST /api/v1/meituan/cancel        // 撤销核销
GET  /api/v1/meituan/orders        // 订单列表
```

**核心方法**:
```go
// services/meituan_service.go
- GetAccessToken() (string, error)
- QueryVoucher(code string) (*VoucherInfo, error)
- VerifyVoucher(code string, operatorID int64) error
- CancelVerification(code string, operatorID int64) error
- SyncOrders(startDate, endDate string) error
```

**验收标准**:
- API调用成功率>99%
- 响应时间<2秒
- 支持自动重试机制
- 异常情况有完善的错误处理

---

#### 2.2.2 美团券核销流程
**业务逻辑**:
```
1. 前台扫描/输入券码
2. 调用美团API查询券信息
3. 验证券状态(未核销、未过期)
4. 验证券类型和产品信息
5. 查询或创建用户
   - 根据手机号查询用户
   - 不存在则创建新用户
6. 根据券类型执行对应操作
   - 会员卡券:创建会员卡
   - 私教课券:创建课程订单
   - 团课券:创建团课预约
7. 调用美团API核销券
8. 更新本地券状态
9. 创建核销记录
10. 返回核销成功
```

**异常处理**:
- 券不存在:提示券码错误
- 券已核销:提示重复核销
- 券已过期:提示券已过期
- API调用失败:记录日志,人工处理

---

### 2.3 抖音平台对接

#### 2.3.1 抖音API集成
**输入**: 抖音开放平台凭证
**输出**: 完整的抖音API对接

**执行内容**:
- 注册抖音开放平台账号
- 创建小程序/应用
- 实现抖音登录授权
- 实现券查询接口
- 实现券核销接口
- 实现订单同步接口

**API路由**:
```
POST /api/v1/douyin/auth           // 抖音授权
GET  /api/v1/douyin/voucher/:code  // 查询券信息
POST /api/v1/douyin/verify         // 核销券
POST /api/v1/douyin/cancel         // 撤销核销
GET  /api/v1/douyin/orders         // 订单列表
```

**核心方法**:
```go
// services/douyin_service.go
- GetAccessToken() (string, error)
- QueryVoucher(code string) (*VoucherInfo, error)
- VerifyVoucher(code string, operatorID int64) error
- CancelVerification(code string, operatorID int64) error
- SyncOrders(startDate, endDate string) error
```

**验收标准**:
- API调用成功率>99%
- 响应时间<2秒
- 支持自动重试机制
- 异常情况有完善的错误处理

---

### 2.4 券管理功能

#### 2.4.1 券列表与查询
**输入**: 查询条件
**输出**: 券列表数据

**执行内容**:

**后端API**:
```go
// API路由
GET /api/v1/vouchers              // 券列表
GET /api/v1/vouchers/:id          // 券详情
GET /api/v1/vouchers/no/:code     // 根据券号查询
GET /api/v1/vouchers/stats        // 券统计

// 查询条件
type ListVouchersRequest struct {
    Page         int    `form:"page"`
    PageSize     int    `form:"page_size"`
    Keyword      string `form:"keyword"`      // 券号、手机号
    PlatformID   int64  `form:"platform_id"`  // 平台
    VoucherType  int8   `form:"voucher_type"` // 券类型
    Status       int8   `form:"status"`       // 状态
    StartDate    string `form:"start_date"`   // 开始日期
    EndDate      string `form:"end_date"`     // 结束日期
}
```

**前端页面**:
```typescript
// src/pages/Voucher/VoucherList.tsx

// 功能点:
1. 券列表展示
   - 表格列:券号、平台、券类型、产品名称、购买用户、状态、过期日期、操作
   - 状态标识:未核销(蓝色)、已核销(绿色)、已过期(灰色)、已退款(红色)

2. 搜索与筛选
   - 关键词搜索(券号、手机号)
   - 平台筛选
   - 券类型筛选
   - 状态筛选
   - 日期范围筛选

3. 操作按钮
   - 核销:打开核销弹窗
   - 查看详情:显示券详细信息
   - 撤销核销:二次确认后撤销

4. 批量操作
   - 批量导出
   - 批量同步
```

**验收标准**:
- 列表查询性能良好
- 筛选条件灵活
- 操作按钮权限控制正确

---

#### 2.4.2 券核销界面
**输入**: 券码
**输出**: 核销结果

**执行内容**:

**前端界面**:
```typescript
// src/components/Voucher/VerifyVoucher.tsx

// 功能点:
1. 券码输入
   - 手动输入券码
   - 扫码枪扫描
   - 二维码扫描

2. 券信息展示
   - 平台来源
   - 券类型
   - 产品名称
   - 原价/实付价
   - 购买用户信息
   - 过期日期

3. 用户信息
   - 自动匹配用户(根据手机号)
   - 显示用户基本信息
   - 不存在则提示创建新用户

4. 核销确认
   - 显示核销后的结果(会员卡信息/课程信息)
   - 确认核销按钮
   - 取消按钮

5. 核销结果
   - 成功提示
   - 打印小票(可选)
   - 发送短信通知(可选)
```

**验收标准**:
- 扫码识别准确
- 核销流程流畅
- 异常提示清晰
- 核销成功率>99%

---

### 2.5 订单同步

#### 2.5.1 定时同步任务
**输入**: 平台订单数据
**输出**: 本地券数据

**执行内容**:

**定时任务**:
```go
// services/voucher_sync_service.go

// 核心方法:
- SyncMeituanOrders() error  // 同步美团订单
- SyncDouyinOrders() error   // 同步抖音订单
- ProcessNewOrders() error   // 处理新订单
- UpdateOrderStatus() error  // 更新订单状态

// 定时任务配置:
- 每15分钟执行一次订单同步
- 每小时执行一次状态更新
- 每天凌晨执行一次数据对账
```

**同步逻辑**:
```
1. 调用平台API获取订单列表
2. 过滤已同步的订单
3. 解析订单信息
4. 创建券记录
5. 更新同步状态
6. 记录同步日志
7. 异常订单告警
```

**验收标准**:
- 订单同步及时(延迟<15分钟)
- 同步准确率100%
- 异常订单有告警
- 支持手动触发同步

---

### 2.6 数据统计与报表

#### 2.6.1 券核销统计
**输入**: 统计时间范围
**输出**: 统计报表

**执行内容**:

**统计指标**:
```
1. 总体统计
   - 总券数量
   - 已核销数量
   - 未核销数量
   - 已过期数量
   - 核销率

2. 平台分布
   - 各平台券数量
   - 各平台核销率
   - 各平台销售额

3. 券类型分布
   - 各类型券数量
   - 各类型核销率

4. 核销趋势
   - 每日核销数量
   - 每日核销金额
```

**前端统计页面**:
```typescript
// src/pages/Voucher/VoucherStats.tsx

// 展示内容:
1. 统计卡片
   - 总券数、已核销数、未核销数、核销率

2. 平台分布图(饼图)
   - 各平台券数量占比

3. 核销趋势图(折线图)
   - X轴:日期
   - Y轴:核销数量/金额

4. 数据导出
   - 导出统计报表
```

**验收标准**:
- 统计数据准确
- 图表展示清晰
- 支持时间范围切换
- 报表可导出

---

## 三、技术方案

### 3.1 API对接方案
- 使用HTTP客户端库(Resty)
- 实现统一的API调用封装
- 支持自动重试和超时控制
- 完善的日志记录

### 3.2 数据安全
- API密钥加密存储
- HTTPS传输
- 签名验证
- 防重放攻击

### 3.3 性能优化
- API调用结果缓存
- 异步处理非关键流程
- 批量同步优化
- 数据库索引优化

---

## 四、接口文档

### 4.1 核心接口

#### 查询券信息
```
GET /api/v1/vouchers/no/:code

Response:
{
  "code": 200,
  "data": {
    "voucher_no": "MT123456",
    "platform_name": "美团",
    "voucher_type": 1,
    "product_name": "月卡",
    "original_price": 299.00,
    "paid_price": 199.00,
    "user_phone": "13800138000",
    "status": 1,
    "expire_date": "2024-12-31"
  }
}
```

#### 核销券
```
POST /api/v1/vouchers/verify

Request:
{
  "voucher_code": "MT123456",
  "operator_id": 1
}

Response:
{
  "code": 200,
  "data": {
    "user_id": 1,
    "result_type": 1,
    "result_id": 100,
    "message": "核销成功,已为用户开通月卡"
  }
}
```

---

## 五、测试用例

### 5.1 功能测试
- 券查询测试
- 券核销测试
- 撤销核销测试
- 订单同步测试

### 5.2 异常测试
- 券不存在测试
- 券已核销测试
- 券已过期测试
- API调用失败测试

### 5.3 性能测试
- 并发核销测试
- 大批量同步测试

---

**任务优先级**: P1  
**预计工期**: 2-3周  
**依赖任务**: TASK001, TASK002  
**后续任务**: TASK005