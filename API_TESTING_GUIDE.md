# API 测试指南

本文档提供了新实现功能的API测试示例。

## 1. 签到功能 API

### 1.1 用户签到
```bash
POST /api/v1/checkins
Content-Type: application/json

{
  "user_id": 1,
  "check_in_type": 2,
  "device_id": "device_001",
  "card_id": 1
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "check-in successful",
  "data": {
    "id": 1,
    "user_id": 1,
    "card_id": 1,
    "check_in_type": 2,
    "check_in_time": "2024-01-15T10:30:00Z",
    "device_id": "device_001"
  }
}
```

### 1.2 获取今日签到记录
```bash
GET /api/v1/users/1/checkin/today
```

### 1.3 查询签到记录列表
```bash
GET /api/v1/checkins?page=1&page_size=10&user_id=1&start_date=2024-01-01&end_date=2024-01-31
```

## 2. 会员卡类型管理 API

### 2.1 创建会员卡类型
```bash
POST /api/v1/card-types
Content-Type: application/json

{
  "type_name": "月卡",
  "type_code": "MONTH_CARD",
  "duration_type": 2,
  "duration_value": 30,
  "price": 299.00,
  "original_price": 399.00,
  "description": "30天有效期会员卡",
  "benefits": "{\"free_courses\": 2, \"discount\": 0.9}",
  "can_freeze": 1,
  "max_freeze_times": 2,
  "max_freeze_days": 7,
  "can_transfer": 1,
  "transfer_fee": 50.00,
  "sort_order": 1
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "card type created successfully",
  "data": {
    "id": 1,
    "type_name": "月卡",
    "type_code": "MONTH_CARD",
    "duration_type": 2,
    "duration_value": 30,
    "price": 299.00,
    "original_price": 399.00,
    "status": 1,
    "created_at": "2024-01-15T10:00:00Z"
  }
}
```

### 2.2 获取会员卡类型列表
```bash
GET /api/v1/card-types?status=1
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "type_name": "月卡",
      "type_code": "MONTH_CARD",
      "duration_type": 2,
      "duration_value": 30,
      "price": 299.00,
      "status": 1
    },
    {
      "id": 2,
      "type_name": "季卡",
      "type_code": "SEASON_CARD",
      "duration_type": 3,
      "duration_value": 90,
      "price": 799.00,
      "status": 1
    }
  ]
}
```

### 2.3 获取会员卡类型详情
```bash
GET /api/v1/card-types/1
```

### 2.4 更新会员卡类型
```bash
PUT /api/v1/card-types/1
Content-Type: application/json

{
  "price": 279.00,
  "description": "30天有效期会员卡（优惠价）"
}
```

### 2.5 启用/停用会员卡类型
```bash
# 启用
POST /api/v1/card-types/1/enable

# 停用
POST /api/v1/card-types/1/disable
```

### 2.6 更新排序
```bash
POST /api/v1/card-types/sort-order
Content-Type: application/json

{
  "orders": {
    "1": 1,
    "2": 2,
    "3": 3
  }
}
```

### 2.7 删除会员卡类型
```bash
DELETE /api/v1/card-types/1
```

**注意:** 如果该类型正在被使用，将返回错误。

## 3. 测试流程建议

### 3.1 会员卡类型管理测试流程
1. 创建多个会员卡类型（月卡、季卡、年卡、次卡）
2. 查询会员卡类型列表，验证创建成功
3. 更新某个会员卡类型的价格
4. 调整会员卡类型的排序
5. 停用某个会员卡类型
6. 尝试删除正在使用的会员卡类型（应该失败）

### 3.2 签到功能测试流程
1. 创建用户和会员卡
2. 执行签到操作
3. 查询今日签到记录
4. 尝试重复签到（应该失败）
5. 查询签到记录列表
6. 验证训练统计数据是否更新

## 4. 错误处理

### 4.1 常见错误码
- `400` - 请求参数错误
- `401` - 未授权
- `404` - 资源不存在
- `500` - 服务器内部错误

### 4.2 错误响应示例
```json
{
  "error": "type code already exists"
}
```

## 5. 数据类型说明

### 5.1 会员卡时长类型 (duration_type)
- `1` - 天卡
- `2` - 月卡
- `3` - 季卡
- `4` - 年卡
- `5` - 次卡

### 5.2 签到类型 (check_in_type)
- `1` - 人脸识别
- `2` - 刷卡
- `3` - 手动签到

### 5.3 状态 (status)
- `1` - 启用/正常
- `2` - 停用/已过期

## 6. 使用 curl 测试示例

### 创建会员卡类型
```bash
curl -X POST http://localhost:8080/api/v1/card-types \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "type_name": "月卡",
    "type_code": "MONTH_CARD",
    "duration_type": 2,
    "duration_value": 30,
    "price": 299.00
  }'
```

### 用户签到
```bash
curl -X POST http://localhost:8080/api/v1/checkins \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "user_id": 1,
    "check_in_type": 2,
    "device_id": "device_001"
  }'
```

## 7. 会员卡操作功能 API

### 7.1 会员卡续费
```bash
POST /api/v1/cards/:id/renew
Content-Type: application/json

{
  "months": 3,
  "amount": 299.00,
  "operator_id": 1,
  "remark": "续费3个月"
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "Card renewed successfully",
  "data": null
}
```

**业务规则:**
- 如果会员卡已过期，从今天开始计算新的到期日期
- 如果会员卡未过期，从当前到期日期延长
- 已转出或已退卡的会员卡不能续费

### 7.2 冻结会员卡
```bash
POST /api/v1/cards/:id/freeze
Content-Type: application/json

{
  "freeze_days": 7,
  "operator_id": 1,
  "remark": "用户请假"
}
```

**业务规则:**
- 只有正常状态的会员卡可以冻结
- 检查会员卡类型是否允许冻结
- 检查冻结次数是否超过限制
- 检查冻结天数是否超过限制
- 冻结后会员卡到期日期自动延长相应天数

### 7.3 解冻会员卡
```bash
POST /api/v1/cards/:id/unfreeze
Content-Type: application/json

{
  "operator_id": 1,
  "remark": "用户销假"
}
```

### 7.4 转卡
```bash
POST /api/v1/cards/:id/transfer
Content-Type: application/json

{
  "to_user_id": 2,
  "transfer_fee": 50.00,
  "operator_id": 1,
  "remark": "转给朋友"
}
```

**业务规则:**
- 只有正常状态的会员卡可以转卡
- 已冻结的会员卡不能转卡
- 检查会员卡类型是否允许转卡
- 检查目标用户是否存在
- 转卡后会员卡状态变为"已转出"

### 7.5 查询操作历史
```bash
GET /api/v1/cards/:id/operations
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "card_id": 1,
      "operation_type": 1,
      "operator_id": 1,
      "amount": 299.00,
      "old_end_date": "2024-01-31",
      "new_end_date": "2024-04-30",
      "remark": "续费3个月",
      "created_at": "2024-01-15T10:00:00Z"
    },
    {
      "id": 2,
      "card_id": 1,
      "operation_type": 2,
      "operator_id": 1,
      "freeze_days": 7,
      "old_end_date": "2024-04-30",
      "new_end_date": "2024-05-07",
      "remark": "用户请假",
      "created_at": "2024-02-01T10:00:00Z"
    }
  ]
}
```

**操作类型说明:**
- `1` - 续费
- `2` - 冻结
- `3` - 解冻
- `4` - 转卡
- `5` - 退卡

## 8. 通知系统 API

### 8.1 获取通知列表
```bash
GET /api/v1/notifications?user_id=1&page=1&page_size=20&is_read=0
```

**查询参数:**
- `user_id` (必填) - 用户ID
- `page` - 页码，默认1
- `page_size` - 每页数量，默认20
- `is_read` - 是否已读: 0-未读, 1-已读

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "user_id": 1,
        "type": 1,
        "title": "会员卡即将到期提醒",
        "content": "您的会员卡（卡号：C202401010001）将在7天后到期，到期日期为2024-02-01，请及时续费。",
        "related_id": 1,
        "is_read": 0,
        "read_at": null,
        "created_at": "2024-01-25T09:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 20
  }
}
```

### 8.2 标记通知为已读
```bash
POST /api/v1/notifications/:id/read
```

### 8.3 全部标记为已读
```bash
POST /api/v1/notifications/read-all
Content-Type: application/json

{
  "user_id": 1
}
```

### 8.4 获取未读数量
```bash
GET /api/v1/notifications/unread-count?user_id=1
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "count": 3
  }
}
```

### 8.5 删除通知
```bash
DELETE /api/v1/notifications/:id
```

## 9. 管理员功能 API

### 9.1 手动触发到期检查
```bash
POST /api/v1/admin/trigger-expiry-check
```

**功能说明:**
- 立即执行会员卡到期检查
- 更新过期会员卡状态
- 发送到期提醒通知

### 9.2 获取定时任务状态
```bash
GET /api/v1/admin/scheduler-status
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "Running - Next check at 2024-01-26 09:00:00"
  }
}
```

### 9.3 获取即将到期会员卡
```bash
GET /api/v1/admin/expiring-cards?days=7
```

**查询参数:**
- `days` - 天数，默认7天

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "days": 7,
    "count": 5,
    "cards": [
      {
        "id": 1,
        "card_no": "C202401010001",
        "user_id": 1,
        "end_date": "2024-02-01",
        "status": 1
      }
    ]
  }
}
```

### 9.4 获取已过期会员卡
```bash
GET /api/v1/admin/expired-cards
```

## 10. 训练统计 API

### 10.1 获取用户基础统计
```bash
GET /api/v1/users/:user_id/stats
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "total_days": 45,
    "total_times": 52,
    "continuous_days": 7,
    "last_check_in_date": "2024-01-25",
    "month_times": 12,
    "year_times": 45
  }
}
```

### 10.2 获取详细统计数据
```bash
GET /api/v1/users/:user_id/stats/detailed
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "basic_stats": {
      "total_days": 45,
      "total_times": 52,
      "continuous_days": 7
    },
    "this_month_days": 12,
    "this_year_days": 45,
    "last_7_days": {
      "total_times": 7,
      "unique_days": 7
    },
    "last_30_days": {
      "total_times": 28,
      "unique_days": 25
    },
    "monthly_trend": [
      {"month": "2023-08", "days": 8},
      {"month": "2023-09", "days": 10},
      {"month": "2023-10", "days": 12},
      {"month": "2023-11", "days": 9},
      {"month": "2023-12", "days": 11},
      {"month": "2024-01", "days": 12}
    ],
    "avg_per_week": 3.5
  }
}
```

### 10.3 获取签到日历
```bash
GET /api/v1/users/:user_id/stats/calendar?year=2024&month=1
```

**查询参数:**
- `year` - 年份，默认当前年
- `month` - 月份（1-12），默认当前月

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "year": 2024,
    "month": 1,
    "total_days": 12,
    "total_times": 15,
    "calendar": [
      {
        "date": "2024-01-01",
        "count": 1,
        "check_ins": [
          {
            "id": 1,
            "user_id": 1,
            "check_in_type": 1,
            "check_in_time": "2024-01-01T10:30:00Z"
          }
        ]
      },
      {
        "date": "2024-01-02",
        "count": 2,
        "check_ins": [...]
      }
    ]
  }
}
```

### 10.4 重新计算统计
```bash
POST /api/v1/users/:user_id/stats/recalculate
```

**功能说明:**
- 从签到历史重新计算所有统计数据
- 用于数据修复或统计异常时
- 会覆盖现有统计数据

## 11. 用户状态管理 API

### 11.1 冻结用户
```bash
POST /api/v1/users/:id/freeze
Content-Type: application/json

{
  "reason": "长期未使用",
  "operator_id": 1
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "User frozen successfully"
}
```

### 11.2 解冻用户
```bash
POST /api/v1/users/:id/unfreeze
Content-Type: application/json

{
  "reason": "用户申请恢复",
  "operator_id": 1
}
```

### 11.3 加入黑名单
```bash
POST /api/v1/users/:id/blacklist
Content-Type: application/json

{
  "reason": "违反规定",
  "operator_id": 1
}
```

### 11.4 移出黑名单
```bash
DELETE /api/v1/users/:id/blacklist
Content-Type: application/json

{
  "reason": "申诉成功",
  "operator_id": 1
}
```

### 11.5 获取状态变更日志
```bash
GET /api/v1/users/:id/status-logs?page=1&page_size=10
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "user_id": 1,
        "old_status": 1,
        "new_status": 2,
        "reason": "长期未使用",
        "operator_id": 1,
        "created_at": "2024-01-25T10:00:00Z"
      }
    ],
    "total": 5,
    "page": 1,
    "page_size": 10
  }
}
```

### 11.6 获取用户状态统计
```bash
GET /api/v1/users/status/summary
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "active": 150,
    "frozen": 10,
    "blacklist": 5,
    "total": 165
  }
}
```

### 11.7 批量冻结用户
```bash
POST /api/v1/users/batch/freeze
Content-Type: application/json

{
  "user_ids": [1, 2, 3],
  "reason": "批量冻结操作",
  "operator_id": 1
}
```

**响应示例:**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "success_count": 2,
    "failed_count": 1,
    "errors": [
      "User 2: user is already in this status"
    ]
  }
}
```

### 11.8 批量解冻用户
```bash
POST /api/v1/users/batch/unfreeze
Content-Type: application/json

{
  "user_ids": [1, 2, 3],
  "reason": "批量解冻操作",
  "operator_id": 1
}
```

## 12. 测试流程建议 (更新)

### 8.1 会员卡完整流程测试
1. 创建会员卡类型（设置冻结和转卡规则）
2. 为用户创建会员卡
3. 测试续费功能
4. 测试冻结功能（验证规则限制）
5. 测试解冻功能
6. 测试转卡功能
7. 查询操作历史记录
8. 验证会员卡状态变化

### 12.2 边界条件测试
1. 尝试续费已转出的会员卡（应该失败）
2. 尝试冻结已冻结的会员卡（应该失败）
3. 超过冻结次数限制（应该失败）
4. 超过冻结天数限制（应该失败）
5. 转卡到不存在的用户（应该失败）
6. 转卡不允许转卡的类型（应该失败）

### 12.3 到期提醒系统测试
1. 创建即将到期的会员卡（设置到期日期为7天后）
2. 手动触发到期检查
3. 查询用户通知列表（验证收到提醒）
4. 测试通知标记为已读
5. 测试获取未读数量
6. 查询即将到期会员卡列表
7. 创建已过期的会员卡（设置到期日期为昨天）
8. 手动触发到期检查
9. 验证会员卡状态自动更新为"已过期"
10. 查询已过期会员卡列表

### 12.4 定时任务测试
1. 启动服务器，验证定时任务自动启动
2. 查询定时任务状态
3. 手动触发到期检查
4. 查看日志输出，验证任务执行情况

### 12.5 训练统计测试
1. 创建用户并进行多次签到
2. 查询基础统计数据，验证总天数、总次数
3. 连续多天签到，验证连续天数计算
4. 跨月签到，验证月度统计重置
5. 查询详细统计，验证近7天、30天数据
6. 查询签到日历，验证按月分组
7. 查询月度趋势，验证最近6个月数据
8. 手动重新计算统计，验证数据一致性
9. 测试边界情况：
   - 同一天多次签到（只计1天）
   - 跨年签到（年度统计重置）
   - 中断签到后再签到（连续天数重置）

### 12.6 用户状态管理测试
1. 创建测试用户
2. 冻结用户，验证状态变更
3. 尝试重复冻结，验证错误提示
4. 解冻用户，验证状态恢复
5. 加入黑名单，验证状态变更
6. 移出黑名单，验证状态恢复
7. 查询状态日志，验证记录完整性
8. 查询状态统计，验证数据准确性
9. 批量冻结多个用户，验证批量操作
10. 批量解冻，验证部分成功场景
11. 测试边界情况：
    - 冻结状态的用户能否签到
    - 黑名单用户能否签到
    - 操作人ID追踪
    - 状态变更原因记录

## 13. Postman 集合

建议创建 Postman 集合来管理这些 API 测试，包括：
- 环境变量配置（base_url, token）
- 预请求脚本（自动添加 token）
- 测试脚本（验证响应）
- 定时任务测试场景
- 统计数据验证脚本
- 状态管理测试场景

---

**更新时间:** 2024年
**版本:** v5.0 (新增用户状态管理功能)
