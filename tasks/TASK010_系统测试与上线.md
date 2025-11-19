# TASK010 - 系统测试与上线

## 一、功能概述

### 1.1 功能目标
完成健身房管理系统的全面测试、性能优化、安全加固和上线部署，确保系统稳定可靠地投入生产环境使用。

### 1.2 核心价值
- 确保系统功能完整性和正确性
- 保障系统性能满足业务需求
- 加固系统安全，防范潜在风险
- 建立完善的运维监控体系
- 提供平滑的上线和回滚方案

### 1.3 涉及角色
- **测试人员**: 执行各类测试，提交测试报告
- **开发人员**: 修复bug，优化性能
- **运维人员**: 部署系统，配置监控
- **项目经理**: 协调资源，把控进度

## 二、功能详细拆解

### 2.1 功能测试

#### 2.1.1 单元测试
**输入**: 各模块代码
**输出**: 单元测试报告

**执行内容**:

**A. 后端单元测试 (Go)**
```go
// tests/unit/user_service_test.go
package tests

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "gym-admin/internal/service"
)

func TestCreateUser(t *testing.T) {
    // 测试用例1：正常创建用户
    t.Run("创建用户成功", func(t *testing.T) {
        req := &service.CreateUserRequest{
            Name:  "测试用户",
            Phone: "13800138000",
        }
        
        user, err := service.CreateUser(req)
        
        assert.NoError(t, err)
        assert.NotNil(t, user)
        assert.Equal(t, "测试用户", user.Name)
        assert.NotEmpty(t, user.UserNo)
    })
    
    // 测试用例2：手机号重复
    t.Run("手机号重复", func(t *testing.T) {
        req := &service.CreateUserRequest{
            Name:  "测试用户2",
            Phone: "13800138000", // 重复的手机号
        }
        
        user, err := service.CreateUser(req)
        
        assert.Error(t, err)
        assert.Nil(t, user)
        assert.Contains(t, err.Error(), "手机号已存在")
    })
    
    // 测试用例3：参数校验
    t.Run("参数校验失败", func(t *testing.T) {
        req := &service.CreateUserRequest{
            Name:  "", // 姓名为空
            Phone: "13800138000",
        }
        
        user, err := service.CreateUser(req)
        
        assert.Error(t, err)
        assert.Nil(t, user)
    })
}

func TestGenerateUserNo(t *testing.T) {
    // 测试用户编号生成
    userNo1 := service.GenerateUserNo()
    userNo2 := service.GenerateUserNo()
    
    assert.NotEmpty(t, userNo1)
    assert.NotEmpty(t, userNo2)
    assert.NotEqual(t, userNo1, userNo2) // 确保唯一性
    assert.Regexp(t, "^U\\d{12}$", userNo1) // 格式校验
}

// 运行测试
// go test -v ./tests/unit/...
// go test -cover ./tests/unit/...
```

**B. 前端单元测试 (React + Jest)**
```typescript
// src/components/User/__tests__/UserForm.test.tsx
import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import UserForm from '../UserForm'

describe('UserForm', () => {
  test('渲染表单', () => {
    render(<UserForm />)
    
    expect(screen.getByLabelText('姓名')).toBeInTheDocument()
    expect(screen.getByLabelText('手机号')).toBeInTheDocument()
    expect(screen.getByRole('button', { name: '提交' })).toBeInTheDocument()
  })
  
  test('表单验证 - 必填项', async () => {
    render(<UserForm />)
    
    const submitButton = screen.getByRole('button', { name: '提交' })
    fireEvent.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('请输入姓名')).toBeInTheDocument()
      expect(screen.getByText('请输入手机号')).toBeInTheDocument()
    })
  })
  
  test('表单验证 - 手机号格式', async () => {
    render(<UserForm />)
    
    const phoneInput = screen.getByLabelText('手机号')
    await userEvent.type(phoneInput, '123')
    
    const submitButton = screen.getByRole('button', { name: '提交' })
    fireEvent.click(submitButton)
    
    await waitFor(() => {
      expect(screen.getByText('请输入正确的手机号')).toBeInTheDocument()
    })
  })
  
  test('提交表单成功', async () => {
    const onSubmit = jest.fn()
    render(<UserForm onSubmit={onSubmit} />)
    
    await userEvent.type(screen.getByLabelText('姓名'), '张三')
    await userEvent.type(screen.getByLabelText('手机号'), '13800138000')
    
    const submitButton = screen.getByRole('button', { name: '提交' })
    fireEvent.click(submitButton)
    
    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith({
        name: '张三',
        phone: '13800138000'
      })
    })
  })
})

// 运行测试
// npm test
// npm test -- --coverage
```

**验收标准**:
- 单元测试覆盖率 > 80%
- 所有测试用例通过
- 关键业务逻辑测试完整
- 边界条件测试充分

---

#### 2.1.2 集成测试
**输入**: 完整的系统模块
**输出**: 集成测试报告

**执行内容**:

**A. API集成测试**
```go
// tests/integration/api_test.go
package tests

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestUserAPI(t *testing.T) {
    // 初始化测试环境
    router := setupTestRouter()
    
    t.Run("创建用户API", func(t *testing.T) {
        reqBody := map[string]interface{}{
            "name":  "测试用户",
            "phone": "13800138000",
        }
        
        body, _ := json.Marshal(reqBody)
        req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+getTestToken())
        
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusOK, w.Code)
        
        var resp map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &resp)
        
        assert.Equal(t, float64(200), resp["code"])
        assert.NotNil(t, resp["data"])
    })
    
    t.Run("获取用户列表API", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/v1/users?page=1&page_size=10", nil)
        req.Header.Set("Authorization", "Bearer "+getTestToken())
        
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusOK, w.Code)
        
        var resp map[string]interface{}
        json.Unmarshal(w.Body.Bytes(), &resp)
        
        assert.Equal(t, float64(200), resp["code"])
        data := resp["data"].(map[string]interface{})
        assert.NotNil(t, data["list"])
        assert.NotNil(t, data["total"])
    })
}

func TestBookingFlow(t *testing.T) {
    // 测试完整的预约流程
    router := setupTestRouter()
    token := getTestToken()
    
    // 1. 获取课程类型
    t.Run("获取课程类型", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/v1/course-types", nil)
        req.Header.Set("Authorization", "Bearer "+token)
        
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusOK, w.Code)
    })
    
    // 2. 获取教练列表
    t.Run("获取教练列表", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/v1/coaches", nil)
        req.Header.Set("Authorization", "Bearer "+token)
        
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusOK, w.Code)
    })
    
    // 3. 创建预约
    t.Run("创建预约", func(t *testing.T) {
        reqBody := map[string]interface{}{
            "user_id":        1,
            "coach_id":       1,
            "course_type_id": 1,
            "booking_date":   "2024-01-15",
            "start_time":     "10:00",
        }
        
        body, _ := json.Marshal(reqBody)
        req := httptest.NewRequest("POST", "/api/v1/bookings", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Authorization", "Bearer "+token)
        
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)
        
        assert.Equal(t, http.StatusOK, w.Code)
    })
}
```

**B. 端到端测试 (E2E)**
```typescript
// e2e/user-management.spec.ts
import { test, expect } from '@playwright/test'

test.describe('用户管理', () => {
  test.beforeEach(async ({ page }) => {
    // 登录
    await page.goto('/login')
    await page.fill('input[name="username"]', 'admin')
    await page.fill('input[name="password"]', 'admin123')
    await page.click('button[type="submit"]')
    await page.waitForURL('/dashboard')
  })
  
  test('创建用户流程', async ({ page }) => {
    // 进入用户管理页面
    await page.goto('/users')
    await expect(page.locator('h1')).toContainText('用户管理')
    
    // 点击新增按钮
    await page.click('button:has-text("新增用户")')
    
    // 填写表单
    await page.fill('input[name="name"]', '测试用户')
    await page.fill('input[name="phone"]', '13800138000')
    await page.fill('input[name="email"]', 'test@example.com')
    
    // 提交表单
    await page.click('button:has-text("提交")')
    
    // 验证成功提示
    await expect(page.locator('.ant-message-success')).toBeVisible()
    
    // 验证用户出现在列表中
    await expect(page.locator('table')).toContainText('测试用户')
    await expect(page.locator('table')).toContainText('13800138000')
  })
  
  test('搜索用户', async ({ page }) => {
    await page.goto('/users')
    
    // 输入搜索关键词
    await page.fill('input[placeholder="搜索用户"]', '张三')
    await page.click('button:has-text("搜索")')
    
    // 验证搜索结果
    await expect(page.locator('table tbody tr')).toHaveCount(1)
    await expect(page.locator('table')).toContainText('张三')
  })
})

test.describe('课程预约', () => {
  test('完整预约流程', async ({ page }) => {
    await page.goto('/booking/create')
    
    // 步骤1：选择课程
    await page.click('.course-item:first-child')
    await page.click('button:has-text("下一步")')
    
    // 步骤2：选择教练
    await page.click('.coach-item:first-child')
    await page.click('button:has-text("下一步")')
    
    // 步骤3：选择时间
    await page.click('.date-item:nth-child(2)')
    await page.click('.time-slot:has-text("10:00")')
    await page.click('button:has-text("下一步")')
    
    // 步骤4：确认预约
    await page.fill('textarea[name="remark"]', '测试预约')
    await page.click('button:has-text("确认预约")')
    
    // 验证成功
    await expect(page.locator('.ant-message-success')).toBeVisible()
  })
})

// 运行测试
// npx playwright test
// npx playwright test --headed
```

**验收标准**:
- 所有API接口测试通过
- 业务流程测试完整
- E2E测试覆盖核心功能
- 测试环境数据隔离

---

### 2.2 性能测试

#### 2.2.1 压力测试
**输入**: 系统接口
**输出**: 性能测试报告

**执行内容**:

**A. 使用JMeter进行压力测试**
```xml
<!-- jmeter-test-plan.jmx -->
<?xml version="1.0" encoding="UTF-8"?>
<jmeterTestPlan version="1.2">
  <hashTree>
    <TestPlan guiclass="TestPlanGui" testclass="TestPlan" testname="健身房管理系统压力测试">
      <elementProp name="TestPlan.user_defined_variables" elementType="Arguments">
        <collectionProp name="Arguments.arguments">
          <elementProp name="BASE_URL" elementType="Argument">
            <stringProp name="Argument.name">BASE_URL</stringProp>
            <stringProp name="Argument.value">http://localhost:8080</stringProp>
          </elementProp>
        </collectionProp>
      </elementProp>
    </TestPlan>
    
    <hashTree>
      <!-- 线程组：模拟100个并发用户 -->
      <ThreadGroup guiclass="ThreadGroupGui" testclass="ThreadGroup" testname="用户并发测试">
        <stringProp name="ThreadGroup.num_threads">100</stringProp>
        <stringProp name="ThreadGroup.ramp_time">10</stringProp>
        <stringProp name="ThreadGroup.duration">60</stringProp>
      </ThreadGroup>
      
      <hashTree>
        <!-- HTTP请求：获取用户列表 -->
        <HTTPSamplerProxy guiclass="HttpTestSampleGui" testclass="HTTPSamplerProxy" testname="获取用户列表">
          <stringProp name="HTTPSampler.domain">${BASE_URL}</stringProp>
          <stringProp name="HTTPSampler.path">/api/v1/users</stringProp>
          <stringProp name="HTTPSampler.method">GET</stringProp>
        </HTTPSamplerProxy>
        
        <!-- 断言：响应时间 < 200ms -->
        <DurationAssertion guiclass="DurationAssertionGui" testclass="DurationAssertion">
          <stringProp name="DurationAssertion.duration">200</stringProp>
        </DurationAssertion>
      </hashTree>
    </hashTree>
  </hashTree>
</jmeterTestPlan>
```

**B. 使用k6进行负载测试**
```javascript
// k6-load-test.js
import http from 'k6/http'
import { check, sleep } from 'k6'

export const options = {
  stages: [
    { duration: '1m', target: 50 },   // 1分钟内增加到50个用户
    { duration: '3m', target: 100 },  // 3分钟内增加到100个用户
    { duration: '2m', target: 200 },  // 2分钟内增加到200个用户
    { duration: '2m', target: 0 },    // 2分钟内降到0
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95%的请求响应时间 < 500ms
    http_req_failed: ['rate<0.01'],   // 错误率 < 1%
  },
}

const BASE_URL = 'http://localhost:8080'
const TOKEN = 'your_test_token'

export default function () {
  // 测试场景1：获取用户列表
  const usersRes = http.get(`${BASE_URL}/api/v1/users?page=1&page_size=20`, {
    headers: { Authorization: `Bearer ${TOKEN}` },
  })
  
  check(usersRes, {
    '状态码为200': (r) => r.status === 200,
    '响应时间 < 200ms': (r) => r.timings.duration < 200,
    '返回数据正确': (r) => JSON.parse(r.body).code === 200,
  })
  
  sleep(1)
  
  // 测试场景2：获取仪表板统计
  const dashboardRes = http.get(`${BASE_URL}/api/v1/stats/dashboard`, {
    headers: { Authorization: `Bearer ${TOKEN}` },
  })
  
  check(dashboardRes, {
    '状态码为200': (r) => r.status === 200,
    '响应时间 < 500ms': (r) => r.timings.duration < 500,
  })
  
  sleep(1)
}

// 运行测试
// k6 run k6-load-test.js
```

**C. 数据库性能优化**
```sql
-- 分析慢查询
SHOW VARIABLES LIKE 'slow_query_log';
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;

-- 查看慢查询日志
SELECT * FROM mysql.slow_log ORDER BY query_time DESC LIMIT 10;

-- 分析查询执行计划
EXPLAIN SELECT * FROM users WHERE phone = '13800138000';

-- 添加必要的索引
CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_check_in_records_user_time ON check_in_records(user_id, check_in_time);
CREATE INDEX idx_course_bookings_coach_date ON course_bookings(coach_id, booking_date);

-- 优化统计查询
CREATE INDEX idx_orders_created_status ON orders(created_at, status);
CREATE INDEX idx_membership_cards_expire ON membership_cards(expire_date, status);
```

**验收标准**:
- 并发100用户时，接口响应时间 < 500ms
- 并发200用户时，系统稳定运行
- 错误率 < 1%
- 数据库查询优化完成
- 慢查询已优化

---

### 2.3 安全测试

#### 2.3.1 安全加固
**输入**: 系统代码和配置
**输出**: 安全加固报告

**执行内容**:

**A. SQL注入防护**
```go
// 使用参数化查询，防止SQL注入
// ❌ 错误示例
query := fmt.Sprintf("SELECT * FROM users WHERE phone = '%s'", phone)
db.Raw(query).Scan(&user)

// ✅ 正确示例
db.Where("phone = ?", phone).First(&user)
```

**B. XSS防护**
```typescript
// 前端输入过滤
import DOMPurify from 'dompurify'

// 清理用户输入
const sanitizeInput = (input: string) => {
  return DOMPurify.sanitize(input, {
    ALLOWED_TAGS: [],
    ALLOWED_ATTR: []
  })
}

// 使用
const userInput = sanitizeInput(formData.name)
```

**C. CSRF防护**
```go
// 添加CSRF中间件
func CSRFMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 验证CSRF Token
        token := c.GetHeader("X-CSRF-Token")
        if token == "" {
            c.JSON(403, gin.H{"error": "CSRF token missing"})
            c.Abort()
            return
        }
        
        // 验证token有效性
        if !validateCSRFToken(token) {
            c.JSON(403, gin.H{"error": "Invalid CSRF token"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

**D. 敏感数据加密**
```go
// 身份证号加密存储
import "crypto/aes"

func EncryptIDCard(idCard string) (string, error) {
    key := []byte(config.EncryptionKey)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    
    // 加密逻辑
    encrypted := encrypt(block, []byte(idCard))
    return base64.StdEncoding.EncodeToString(encrypted), nil
}

func DecryptIDCard(encrypted string) (string, error) {
    key := []byte(config.EncryptionKey)
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }
    
    data, _ := base64.StdEncoding.DecodeString(encrypted)
    decrypted := decrypt(block, data)
    return string(decrypted), nil
}
```

**E. API限流**
```go
// 使用Redis实现限流
func RateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        key := fmt.Sprintf("rate_limit:%s", ip)
        
        // 获取当前请求次数
        count, _ := cache.Client.Incr(ctx, key).Result()
        
        if count == 1 {
            // 设置过期时间（1分钟）
            cache.Client.Expire(ctx, key, time.Minute)
        }
        
        // 限制每分钟100次请求
        if count > 100 {
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

**验收标准**:
- SQL注入测试通过
- XSS攻击测试通过
- CSRF防护生效
- 敏感数据加密存储
- API限流正常工作
- 安全扫描无高危漏洞

---

### 2.4 部署上线

#### 2.4.1 生产环境部署
**输入**: 测试通过的代码
**输出**: 生产环境运行的系统

**执行内容**:

**A. 部署脚本 (deploy.sh)**
```bash
#!/bin/bash

# 健身房管理系统部署脚本

set -e

echo "========== 开始部署 =========="

# 1. 拉取最新代码
echo "1. 拉取最新代码..."
git pull origin main

# 2. 备份数据库
echo "2. 备份数据库..."
BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"
mysqldump -u root -p$DB_PASSWORD gym_admin > /backup/$BACKUP_FILE
echo "数据库备份完成: $BACKUP_FILE"

# 3. 构建后端
echo "3. 构建后端..."
cd backend
go build -o bin/server cmd/server/main.go
echo "后端构建完成"

# 4. 构建前端
echo "4. 构建前端..."
cd ../frontend
npm install
npm run build
echo "前端构建完成"

# 5. 停止旧服务
echo "5. 停止旧服务..."
pm2 stop gym-admin-backend || true
nginx -s stop || true

# 6. 部署后端
echo "6. 部署后端..."
cd ../backend
pm2 start bin/server --name gym-admin-backend

# 7. 部署前端
echo "7. 部署前端..."
cd ../frontend
rm -rf /var/www/gym-admin/*
cp -r dist/* /var/www/gym-admin/
nginx

# 8. 数据库迁移
echo "8. 执行数据库迁移..."
cd ../backend
./bin/server migrate

# 9. 健康检查
echo "9. 健康检查..."
sleep 5
HEALTH_CHECK=$(curl -s http://localhost:8080/health)
if [[ $HEALTH_CHECK == *"ok"* ]]; then
    echo "✅ 健康检查通过"
else
    echo "❌ 健康检查失败"
    exit 1
fi

echo "========== 部署完成 =========="
```

**B. Nginx配置**
```nginx
# /etc/nginx/sites-available/gym-admin
server {
    listen 80;
    server_name gym-admin.example.com;
    
    # 重定向到HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name gym-admin.example.com;
    
    # SSL证书配置
    ssl_certificate /etc/nginx/ssl/gym-admin.crt;
    ssl_certificate_key /etc/nginx/ssl/gym-admin.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    
    # 前端静态文件
    location / {
        root /var/www/gym-admin;
        index index.html;
        try_files $uri $uri/ /index.html;
        
        # 缓存配置
        expires 7d;
        add_header Cache-Control "public, immutable";
    }
    
    # API代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时配置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
    
    # 日志配置
    access_log /var/log/nginx/gym-admin-access.log;
    error_log /var/log/nginx/gym-admin-error.log;
}
```

**C. PM2配置**
```json
// ecosystem.config.js
module.exports = {
  apps: [{
    name: 'gym-admin-backend',
    script: './bin/server',
    instances: 4,
    exec_mode: 'cluster',
    env: {
      NODE_ENV: 'production',
      PORT: 8080
    },
    error_file: './logs/err.log',
    out_file: './logs/out.log',
    log_date_format: 'YYYY-MM-DD HH:mm:ss',
    merge_logs: true,
    max_memory_restart: '1G',
    autorestart: true,
    watch: false
  }]
}
```

**D. 监控配置 (Prometheus + Grafana)**
```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'gym-admin-backend'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    
  - job_name: 'mysql'
    static_configs:
      - targets: ['localhost:9104']
      
  - job_name: 'redis'
    static_configs:
      - targets: ['localhost:9121']
      
  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']

alerting:
  alertmanagers:
    - static_configs:
        - targets: ['localhost:9093']

rule_files:
  - 'alerts.yml'
```

**E. 告警规则**
```yaml
# alerts.yml
groups:
  - name: gym-admin-alerts
    rules:
      # API响应时间告警
      - alert: HighAPILatency
        expr: http_request_duration_seconds{quantile="0.95"} > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "API响应时间过高"
          description: "95%的请求响应时间超过1秒"
      
      # 错误率告警
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "错误率过高"
          description: "5xx错误率超过5%"
      
      # 数据库连接数告警
      - alert: HighDBConnections
        expr: mysql_global_status_threads_connected > 100
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "数据库连接数过高"
          description: "MySQL连接数超过100"
      
      # 磁盘空间告警
      - alert: LowDiskSpace
        expr: (node_filesystem_avail_bytes / node_filesystem_size_bytes) < 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "磁盘空间不足"
          description: "可用磁盘空间低于10%"
```

**验收标准**:
- 部署脚本执行成功
- 服务正常启动
- 健康检查通过
- 监控告警配置完成
- 日志收集正常

---

### 2.5 上线检查清单

#### 2.5.1 上线前检查
**执行内容**:

**完整检查清单**:

```markdown
## 功能检查
- [ ] 用户管理功能正常
- [ ] 会员卡管理功能正常
- [ ] 人脸识别功能正常
- [ ] 美团抖音券核销功能正常
- [ ] 私教课预约功能正常
- [ ] 教练管理功能正常
- [ ] 数据统计功能正常
- [ ] 微信小程序功能正常

## 性能检查
- [ ] 接口响应时间 < 500ms
- [ ] 并发100用户测试通过
- [ ] 数据库查询优化完成
- [ ] 缓存策略配置正确
- [ ] 静态资源CDN配置

## 安全检查
- [ ] SQL注入防护
- [ ] XSS防护
- [ ] CSRF防护
- [ ] 敏感数据加密
- [ ] API限流配置
- [ ] HTTPS配置
- [ ] 防火墙规则

## 数据检查
- [ ] 数据库备份策略
- [ ] 数据迁移脚本测试
- [ ] 数据一致性验证
- [ ] 历史数据清理

## 运维检查
- [ ] 监控告警配置
- [ ] 日志收集配置
- [ ] 自动重启配置
- [ ] 负载均衡配置
- [ ] 域名解析配置
- [ ] SSL证书配置

## 文档检查
- [ ] API文档完整
- [ ] 部署文档完整
- [ ] 运维手册完整
- [ ] 用户手册完整
- [ ] 应急预案完整

## 回滚准备
- [ ] 数据库备份完成
- [ ] 代码版本标记
- [ ] 回滚脚本准备
- [ ] 回滚流程测试
```

---

## 三、上线流程

### 3.1 上线步骤

1. **上线前准备** (T-1天)
   - 代码冻结
   - 完成所有测试
   - 准备部署脚本
   - 备份生产数据

2. **灰度发布** (T-Day 00:00-02:00)
   - 部署到10%服务器
   - 观察监控指标
   - 验证核心功能
   - 收集用户反馈

3. **全量发布** (T-Day 02:00-04:00)
   - 部署到所有服务器
   - 全面监控系统
   - 准备应急响应

4. **上线后观察** (T-Day 04:00-24:00)
   - 持续监控系统
   - 及时处理问题
   - 收集用户反馈

### 3.2 回滚方案

```bash
#!/bin/bash
# rollback.sh - 回滚脚本

echo "========== 开始回滚 =========="

# 1. 停止当前服务
pm2 stop gym-admin-backend
nginx -s stop

# 2. 恢复代码
git checkout $PREVIOUS_VERSION

# 3. 恢复数据库
mysql -u root -p$DB_PASSWORD gym_admin < /backup/$BACKUP_FILE

# 4. 重启服务
pm2 start gym-admin-backend
nginx

echo "========== 回滚完成 =========="
```

---

## 四、测试报告模板

### 4.1 测试报告

```markdown
# 健身房管理系统测试报告

## 一、测试概述
- 测试时间: 2024-01-01 ~ 2024-01-15
- 测试环境: 测试环境
- 测试人员: 测试团队
- 测试版本: v1.0.0

## 二、测试结果汇总
| 测试类型 | 用例总数 | 通过数 | 失败数 | 通过率 |
|---------|---------|--------|--------|--------|
| 功能测试 | 150 | 148 | 2 | 98.7% |
| 性能测试 | 20 | 20 | 0 | 100% |
| 安全测试 | 30 | 30 | 0 | 100% |
| 兼容性测试 | 15 | 15 | 0 | 100% |

## 三、发现的问题
1. 【高】会员卡到期提醒功能偶现失败
2. 【中】用户列表导出Excel格式问题

## 四、性能指标
- 接口平均响应时间: 180ms
- 并发100用户: 稳定
- 数据库查询优化: 完成
- 缓存命中率: 85%

## 五、建议
1. 修复高优先级bug后再上线
2. 增加监控告警规则
3. 完善应急预案
```

---

## 五、上线检查清单

- [ ] 所有测试通过
- [ ] 性能指标达标
- [ ] 安全加固完成
- [ ] 部署脚本准备
- [ ] 监控告警配置
- [ ] 数据备份完成
- [ ] 回滚方案准备
- [ ] 文档完善
- [ ] 团队培训完成
- [ ] 应急预案准备

---

**任务优先级**: P0（上线必须）  
**预计工期**: 2-3周  
**依赖任务**: TASK001-TASK009（所有功能）  
**后续任务**: 运维维护

