# TASK003 - 人脸识别入场功能

## 一、功能概述

### 1.1 功能目标
实现基于人脸识别的智能入场系统,与市面主流摄像头品牌联动,提供快速、便捷、安全的入场验证方式。

### 1.2 核心价值
- 提升入场效率,减少排队等待时间
- 降低前台人工成本
- 提供更好的会员体验
- 自动记录签到数据,便于统计分析

### 1.3 涉及角色
- **管理员**: 配置人脸识别设备、管理人脸数据
- **前台人员**: 协助会员录入人脸、处理识别异常
- **会员**: 录入人脸、使用人脸识别入场

## 二、功能详细拆解

### 2.1 人脸数据管理

#### 2.1.1 数据库设计
```sql
-- 人脸信息表
CREATE TABLE face_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    face_id VARCHAR(100) UNIQUE NOT NULL COMMENT '人脸ID(第三方SDK返回)',
    face_image_url VARCHAR(255) COMMENT '人脸照片URL',
    face_feature TEXT COMMENT '人脸特征数据(加密存储)',
    quality_score DECIMAL(5,2) COMMENT '人脸质量分数',
    device_id VARCHAR(50) COMMENT '录入设备ID',
    status TINYINT DEFAULT 1 COMMENT '状态:1-正常,2-已删除',
    operator_id BIGINT COMMENT '操作员ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_face_id (face_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='人脸信息表';

-- 人脸识别记录表
CREATE TABLE face_recognition_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT COMMENT '识别到的用户ID',
    device_id VARCHAR(50) NOT NULL COMMENT '设备ID',
    recognition_time TIMESTAMP NOT NULL COMMENT '识别时间',
    confidence DECIMAL(5,2) COMMENT '置信度',
    is_success TINYINT NOT NULL COMMENT '是否成功:0-失败,1-成功',
    fail_reason VARCHAR(255) COMMENT '失败原因',
    image_url VARCHAR(255) COMMENT '识别时的照片URL',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_device_id (device_id),
    INDEX idx_recognition_time (recognition_time)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='人脸识别记录表';

-- 设备管理表
CREATE TABLE face_devices (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    device_no VARCHAR(50) UNIQUE NOT NULL COMMENT '设备编号',
    device_name VARCHAR(100) NOT NULL COMMENT '设备名称',
    device_type VARCHAR(50) NOT NULL COMMENT '设备类型/品牌',
    location VARCHAR(100) COMMENT '设备位置',
    ip_address VARCHAR(50) COMMENT 'IP地址',
    api_key VARCHAR(255) COMMENT 'API密钥',
    status TINYINT DEFAULT 1 COMMENT '状态:1-在线,2-离线,3-故障',
    last_heartbeat TIMESTAMP COMMENT '最后心跳时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_device_no (device_no),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='人脸识别设备表';
```

#### 2.1.2 人脸录入功能
**输入**: 用户信息、人脸照片
**输出**: 人脸特征数据

**执行内容**:
- 调用第三方人脸识别SDK提取特征
- 质量检测(光线、角度、清晰度)
- 活体检测(防止照片攻击)
- 特征数据加密存储
- 支持多张人脸录入(提高识别率)

**API路由**:
```
POST /api/v1/face/enroll          // 录入人脸
DELETE /api/v1/face/:id           // 删除人脸
GET /api/v1/users/:user_id/faces  // 获取用户人脸列表
PUT /api/v1/face/:id              // 更新人脸
```

**前端界面**:
- 人脸采集组件(调用摄像头)
- 实时预览和质量提示
- 多角度采集引导
- 采集结果展示

**验收标准**:
- 人脸质量分数>80分才能录入
- 支持活体检测
- 录入成功率>95%
- 单次录入时间<10秒

---

### 2.2 摄像头设备对接

#### 2.2.1 设备SDK集成
**支持的主流品牌**:
- 海康威视
- 大华
- 宇视
- 华为
- 旷视

**集成方式**:
- 统一设备接口抽象层
- 适配器模式支持多品牌
- 设备注册与认证
- 心跳检测机制

**API设计**:
```go
// 设备接口抽象
type FaceDevice interface {
    Connect() error
    Disconnect() error
    EnrollFace(userID int64, image []byte) (*FaceData, error)
    RecognizeFace(image []byte) (*RecognitionResult, error)
    DeleteFace(faceID string) error
    GetDeviceStatus() (*DeviceStatus, error)
}

// 具体实现
- HikVisionDevice  // 海康威视
- DahuaDevice      // 大华
- UnivDevice       // 宇视
```

**验收标准**:
- 至少支持2个主流品牌
- 设备连接成功率>99%
- 设备状态实时监控
- 支持设备热插拔

---

### 2.3 人脸识别入场

#### 2.3.1 识别流程
**业务逻辑**:
```
1. 摄像头捕获人脸
2. 调用识别接口
3. 匹配人脸库
4. 验证会员卡状态
   - 是否有有效会员卡
   - 会员卡是否过期
   - 会员卡是否冻结
5. 验证今日是否已签到
6. 创建签到记录
7. 更新训练统计
8. 开闸放行
9. 显示欢迎信息
```

**API路由**:
```
POST /api/v1/face/recognize      // 人脸识别
POST /api/v1/face/check-in       // 识别后签到
GET /api/v1/face/recognition-logs // 识别记录
```

**识别参数配置**:
- 置信度阈值(默认80%)
- 识别超时时间(默认3秒)
- 失败重试次数(默认3次)
- 陌生人告警开关

**验收标准**:
- 识别准确率>98%
- 识别速度<1秒
- 误识率<1%
- 拒识率<2%

---

### 2.4 异常处理

#### 2.4.1 识别失败处理
**失败场景**:
- 人脸未录入
- 识别置信度低
- 会员卡已过期
- 会员卡已冻结
- 用户状态异常
- 设备故障

**处理策略**:
- 语音/屏幕提示失败原因
- 引导到前台人工处理
- 记录失败日志
- 异常告警通知

#### 2.4.2 设备故障处理
**监控指标**:
- 设备在线状态
- 识别成功率
- 响应时间
- 错误率

**故障处理**:
- 自动切换备用设备
- 降级到手动签到
- 故障告警通知
- 设备重启机制

---

### 2.5 管理后台

#### 2.5.1 设备管理
**功能点**:
- 设备列表展示
- 设备添加/编辑/删除
- 设备状态监控
- 设备参数配置
- 设备日志查看

#### 2.5.2 人脸库管理
**功能点**:
- 人脸列表查询
- 批量导入人脸
- 人脸质量检测
- 人脸库同步
- 人脸数据备份

#### 2.5.3 识别记录查询
**功能点**:
- 识别记录列表
- 多条件筛选
- 识别统计报表
- 异常记录告警
- 数据导出

---

## 三、技术方案

### 3.1 人脸识别SDK选型
**推荐方案**:
- 百度AI人脸识别
- 腾讯云人脸识别
- 阿里云人脸识别
- 旷视Face++

**选型标准**:
- 识别准确率
- 响应速度
- 价格成本
- 技术支持

### 3.2 数据安全
- 人脸特征数据加密存储
- 传输过程HTTPS加密
- 访问权限控制
- 数据脱敏处理
- 符合《个人信息保护法》

### 3.3 性能优化
- 人脸特征数据缓存
- 识别结果缓存
- 数据库索引优化
- 异步处理非关键流程

---

## 四、接口文档

### 4.1 核心接口

#### 人脸录入
```
POST /api/v1/face/enroll

Request:
{
  "user_id": 1,
  "image": "base64_encoded_image",
  "device_id": "device_001"
}

Response:
{
  "code": 200,
  "data": {
    "face_id": "face_123456",
    "quality_score": 95.5,
    "face_image_url": "https://..."
  }
}
```

#### 人脸识别
```
POST /api/v1/face/recognize

Request:
{
  "image": "base64_encoded_image",
  "device_id": "device_001"
}

Response:
{
  "code": 200,
  "data": {
    "user_id": 1,
    "user_name": "张三",
    "confidence": 98.5,
    "card_status": "normal",
    "can_check_in": true
  }
}
```

---

## 五、测试用例

### 5.1 功能测试
- 人脸录入流程测试
- 人脸识别流程测试
- 签到流程测试
- 异常场景测试

### 5.2 性能测试
- 识别速度测试
- 并发识别测试
- 大库容量测试(10000+人脸)

### 5.3 安全测试
- 照片攻击测试
- 视频攻击测试
- 数据加密测试

---

**任务优先级**: P1  
**预计工期**: 3-4周  
**依赖任务**: TASK001, TASK002  
**后续任务**: TASK004
