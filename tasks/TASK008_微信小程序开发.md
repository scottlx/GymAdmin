# TASK008 - 微信小程序开发

## 一、功能概述

### 1.1 功能目标
开发健身房管理系统的微信小程序端，为会员提供便捷的移动端服务，包括会员信息查看、课程预约、签到打卡、教练选择、消息通知等核心功能。

### 1.2 核心价值
- 提供便捷的移动端入口，提升会员体验
- 实现线上线下服务闭环
- 降低前台工作量，提高运营效率
- 增强会员粘性和活跃度
- 支持社交分享，助力营销推广

### 1.3 涉及角色
- **会员**: 查看个人信息、预约课程、签到打卡、查看训练记录
- **教练**: 查看课程安排、确认预约、查看学员信息
- **管理员**: 通过小程序发布通知、查看运营数据

## 二、功能详细拆解

### 2.1 小程序基础架构搭建

#### 2.1.1 项目初始化
**输入**: 微信小程序开发者账号
**输出**: 小程序项目框架

**执行内容**:

**A. 项目目录结构**
```
gym-admin-miniprogram/
├── pages/                          # 页面目录
│   ├── index/                      # 首页
│   │   ├── index.js
│   │   ├── index.json
│   │   ├── index.wxml
│   │   └── index.wxss
│   ├── login/                      # 登录页
│   ├── user/                       # 个人中心
│   │   ├── profile/                # 个人信息
│   │   ├── card/                   # 我的会员卡
│   │   ├── records/                # 训练记录
│   │   └── settings/               # 设置
│   ├── booking/                    # 预约相关
│   │   ├── list/                   # 预约列表
│   │   ├── create/                 # 创建预约
│   │   ├── detail/                 # 预约详情
│   │   └── calendar/               # 日历视图
│   ├── coach/                      # 教练相关
│   │   ├── list/                   # 教练列表
│   │   ├── detail/                 # 教练详情
│   │   └── schedule/               # 教练排班
│   └── checkin/                    # 签到
│       ├── scan/                   # 扫码签到
│       └── face/                   # 人脸签到
├── components/                     # 组件目录
│   ├── card/                       # 卡片组件
│   ├── calendar/                   # 日历组件
│   ├── rating/                     # 评分组件
│   └── empty/                      # 空状态组件
├── utils/                          # 工具函数
│   ├── request.js                  # 网络请求封装
│   ├── auth.js                     # 认证相关
│   ├── storage.js                  # 本地存储
│   └── util.js                     # 通用工具
├── api/                            # API接口
│   ├── user.js
│   ├── booking.js
│   ├── coach.js
│   └── checkin.js
├── config/                         # 配置文件
│   └── config.js
├── images/                         # 图片资源
├── styles/                         # 全局样式
│   └── common.wxss
├── app.js                          # 小程序逻辑
├── app.json                        # 小程序配置
├── app.wxss                        # 小程序样式
├── project.config.json             # 项目配置
└── sitemap.json                    # 站点地图
```

**B. 全局配置 (app.json)**
```json
{
  "pages": [
    "pages/index/index",
    "pages/login/login",
    "pages/user/profile/profile",
    "pages/user/card/card",
    "pages/user/records/records",
    "pages/booking/list/list",
    "pages/booking/create/create",
    "pages/booking/detail/detail",
    "pages/coach/list/list",
    "pages/coach/detail/detail",
    "pages/checkin/scan/scan"
  ],
  "window": {
    "backgroundTextStyle": "light",
    "navigationBarBackgroundColor": "#fff",
    "navigationBarTitleText": "健身房管理",
    "navigationBarTextStyle": "black",
    "backgroundColor": "#f5f5f5"
  },
  "tabBar": {
    "color": "#999999",
    "selectedColor": "#1890ff",
    "backgroundColor": "#ffffff",
    "borderStyle": "black",
    "list": [
      {
        "pagePath": "pages/index/index",
        "text": "首页",
        "iconPath": "images/tab/home.png",
        "selectedIconPath": "images/tab/home-active.png"
      },
      {
        "pagePath": "pages/booking/list/list",
        "text": "预约",
        "iconPath": "images/tab/booking.png",
        "selectedIconPath": "images/tab/booking-active.png"
      },
      {
        "pagePath": "pages/checkin/scan/scan",
        "text": "签到",
        "iconPath": "images/tab/checkin.png",
        "selectedIconPath": "images/tab/checkin-active.png"
      },
      {
        "pagePath": "pages/user/profile/profile",
        "text": "我的",
        "iconPath": "images/tab/user.png",
        "selectedIconPath": "images/tab/user-active.png"
      }
    ]
  },
  "permission": {
    "scope.userLocation": {
      "desc": "你的位置信息将用于签到定位"
    },
    "scope.camera": {
      "desc": "需要使用你的摄像头进行扫码签到"
    }
  },
  "requiredBackgroundModes": ["audio"],
  "usingComponents": {}
}
```

**C. 网络请求封装 (utils/request.js)**
```javascript
const config = require('../config/config.js')
const auth = require('./auth.js')

// 请求拦截器
function request(options) {
  return new Promise((resolve, reject) => {
    // 显示加载提示
    if (options.showLoading !== false) {
      wx.showLoading({
        title: '加载中...',
        mask: true
      })
    }

    // 获取token
    const token = auth.getToken()

    wx.request({
      url: config.apiBaseUrl + options.url,
      method: options.method || 'GET',
      data: options.data || {},
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : ''
      },
      success: (res) => {
        wx.hideLoading()
        
        if (res.statusCode === 200) {
          if (res.data.code === 200) {
            resolve(res.data.data)
          } else if (res.data.code === 401) {
            // token过期，跳转登录
            auth.clearToken()
            wx.redirectTo({
              url: '/pages/login/login'
            })
            reject(new Error('请先登录'))
          } else {
            wx.showToast({
              title: res.data.message || '请求失败',
              icon: 'none'
            })
            reject(new Error(res.data.message))
          }
        } else {
          wx.showToast({
            title: '网络错误',
            icon: 'none'
          })
          reject(new Error('网络错误'))
        }
      },
      fail: (err) => {
        wx.hideLoading()
        wx.showToast({
          title: '网络请求失败',
          icon: 'none'
        })
        reject(err)
      }
    })
  })
}

module.exports = {
  get: (url, data, options = {}) => {
    return request({
      url,
      method: 'GET',
      data,
      ...options
    })
  },
  post: (url, data, options = {}) => {
    return request({
      url,
      method: 'POST',
      data,
      ...options
    })
  },
  put: (url, data, options = {}) => {
    return request({
      url,
      method: 'PUT',
      data,
      ...options
    })
  },
  delete: (url, data, options = {}) => {
    return request({
      url,
      method: 'DELETE',
      data,
      ...options
    })
  }
}
```

**验收标准**:
- 项目结构清晰
- 网络请求封装完善
- 支持token认证
- 错误处理完整

---

### 2.2 用户登录与授权

#### 2.2.1 微信登录实现
**输入**: 微信授权
**输出**: 用户登录状态

**执行内容**:

**A. 登录页面 (pages/login/login.wxml)**
```xml
<view class="login-container">
  <view class="logo">
    <image src="/images/logo.png" mode="aspectFit"></image>
  </view>
  
  <view class="title">欢迎来到健身房</view>
  <view class="subtitle">开启你的健康之旅</view>
  
  <view class="login-methods">
    <!-- 微信授权登录 -->
    <button 
      class="login-btn wechat" 
      open-type="getUserInfo" 
      bindgetuserinfo="onGetUserInfo"
      wx:if="{{!hasUserInfo}}"
    >
      <image src="/images/wechat-icon.png"></image>
      <text>微信授权登录</text>
    </button>
    
    <!-- 手机号登录 -->
    <button 
      class="login-btn phone" 
      bindtap="onPhoneLogin"
      wx:if="{{hasUserInfo}}"
    >
      <image src="/images/phone-icon.png"></image>
      <text>手机号登录</text>
    </button>
  </view>
  
  <view class="agreement">
    <checkbox-group bindchange="onAgreeChange">
      <label>
        <checkbox value="agree" checked="{{agreed}}"/>
        <text>我已阅读并同意</text>
        <text class="link" bindtap="onShowAgreement">《用户协议》</text>
        <text>和</text>
        <text class="link" bindtap="onShowPrivacy">《隐私政策》</text>
      </label>
    </checkbox-group>
  </view>
</view>
```

**B. 登录逻辑 (pages/login/login.js)**
```javascript
const auth = require('../../utils/auth.js')
const userApi = require('../../api/user.js')

Page({
  data: {
    hasUserInfo: false,
    agreed: false
  },

  onLoad() {
    // 检查是否已登录
    if (auth.isLoggedIn()) {
      wx.switchTab({
        url: '/pages/index/index'
      })
    }
  },

  // 获取用户信息
  onGetUserInfo(e) {
    if (e.detail.userInfo) {
      this.setData({
        hasUserInfo: true,
        userInfo: e.detail.userInfo
      })
    } else {
      wx.showToast({
        title: '需要授权才能登录',
        icon: 'none'
      })
    }
  },

  // 手机号登录
  onPhoneLogin() {
    if (!this.data.agreed) {
      wx.showToast({
        title: '请先同意用户协议',
        icon: 'none'
      })
      return
    }

    // 微信登录
    wx.login({
      success: (res) => {
        if (res.code) {
          // 发送code到后端
          this.loginWithCode(res.code)
        } else {
          wx.showToast({
            title: '登录失败',
            icon: 'none'
          })
        }
      }
    })
  },

  // 使用code登录
  async loginWithCode(code) {
    try {
      wx.showLoading({
        title: '登录中...'
      })

      const result = await userApi.login({
        code: code,
        userInfo: this.data.userInfo
      })

      // 保存token
      auth.setToken(result.token)
      
      // 保存用户信息
      wx.setStorageSync('userInfo', result.userInfo)

      wx.hideLoading()
      wx.showToast({
        title: '登录成功',
        icon: 'success'
      })

      // 跳转到首页
      setTimeout(() => {
        wx.switchTab({
          url: '/pages/index/index'
        })
      }, 1500)

    } catch (err) {
      wx.hideLoading()
      wx.showToast({
        title: err.message || '登录失败',
        icon: 'none'
      })
    }
  },

  // 同意协议
  onAgreeChange(e) {
    this.setData({
      agreed: e.detail.value.length > 0
    })
  },

  // 显示用户协议
  onShowAgreement() {
    wx.navigateTo({
      url: '/pages/agreement/agreement'
    })
  },

  // 显示隐私政策
  onShowPrivacy() {
    wx.navigateTo({
      url: '/pages/privacy/privacy'
    })
  }
})
```

**C. 认证工具 (utils/auth.js)**
```javascript
const TOKEN_KEY = 'auth_token'

module.exports = {
  // 保存token
  setToken(token) {
    wx.setStorageSync(TOKEN_KEY, token)
  },

  // 获取token
  getToken() {
    return wx.getStorageSync(TOKEN_KEY)
  },

  // 清除token
  clearToken() {
    wx.removeStorageSync(TOKEN_KEY)
  },

  // 检查是否登录
  isLoggedIn() {
    return !!this.getToken()
  },

  // 登出
  logout() {
    this.clearToken()
    wx.removeStorageSync('userInfo')
    wx.redirectTo({
      url: '/pages/login/login'
    })
  }
}
```

**验收标准**:
- 微信授权登录成功
- Token保存和读取正确
- 登录状态持久化
- 未登录自动跳转登录页

---

### 2.3 首页开发

#### 2.3.1 首页布局与功能
**输入**: 用户登录状态
**输出**: 首页展示

**执行内容**:

**A. 首页布局 (pages/index/index.wxml)**
```xml
<view class="container">
  <!-- 顶部轮播图 -->
  <swiper class="banner" indicator-dots autoplay interval="3000" circular>
    <swiper-item wx:for="{{banners}}" wx:key="id">
      <image src="{{item.image}}" mode="aspectFill"></image>
    </swiper-item>
  </swiper>

  <!-- 会员卡信息 -->
  <view class="card-info" wx:if="{{membershipCard}}">
    <view class="card-header">
      <text class="card-name">{{membershipCard.name}}</text>
      <text class="card-status">{{membershipCard.statusText}}</text>
    </view>
    <view class="card-body">
      <view class="card-item">
        <text class="label">剩余天数</text>
        <text class="value">{{membershipCard.remainingDays}}天</text>
      </view>
      <view class="card-item">
        <text class="label">剩余次数</text>
        <text class="value">{{membershipCard.remainingTimes}}次</text>
      </view>
    </view>
    <view class="card-footer">
      <text>有效期至：{{membershipCard.expireDate}}</text>
    </view>
  </view>

  <!-- 快捷功能 -->
  <view class="quick-actions">
    <view class="action-item" bindtap="onQuickCheckin">
      <image src="/images/icon-checkin.png"></image>
      <text>快速签到</text>
    </view>
    <view class="action-item" bindtap="onQuickBooking">
      <image src="/images/icon-booking.png"></image>
      <text>预约课程</text>
    </view>
    <view class="action-item" bindtap="onViewCoaches">
      <image src="/images/icon-coach.png"></image>
      <text>选择教练</text>
    </view>
    <view class="action-item" bindtap="onViewRecords">
      <image src="/images/icon-record.png"></image>
      <text>训练记录</text>
    </view>
  </view>

  <!-- 训练统计 -->
  <view class="stats-section">
    <view class="section-title">本月训练统计</view>
    <view class="stats-grid">
      <view class="stat-item">
        <text class="stat-value">{{stats.monthTimes}}</text>
        <text class="stat-label">训练次数</text>
      </view>
      <view class="stat-item">
        <text class="stat-value">{{stats.continuousDays}}</text>
        <text class="stat-label">连续天数</text>
      </view>
      <view class="stat-item">
        <text class="stat-value">{{stats.totalHours}}</text>
        <text class="stat-label">训练时长(h)</text>
      </view>
    </view>
  </view>

  <!-- 即将到来的预约 -->
  <view class="upcoming-section" wx:if="{{upcomingBookings.length > 0}}">
    <view class="section-title">
      <text>即将到来的课程</text>
      <text class="more" bindtap="onViewAllBookings">查看全部</text>
    </view>
    <view class="booking-list">
      <view 
        class="booking-item" 
        wx:for="{{upcomingBookings}}" 
        wx:key="id"
        bindtap="onViewBookingDetail"
        data-id="{{item.id}}"
      >
        <view class="booking-time">
          <text class="date">{{item.date}}</text>
          <text class="time">{{item.time}}</text>
        </view>
        <view class="booking-info">
          <text class="course">{{item.courseName}}</text>
          <text class="coach">教练：{{item.coachName}}</text>
        </view>
        <view class="booking-status">
          <text class="status-tag">{{item.statusText}}</text>
        </view>
      </view>
    </view>
  </view>

  <!-- 推荐教练 -->
  <view class="coach-section">
    <view class="section-title">
      <text>明星教练</text>
      <text class="more" bindtap="onViewAllCoaches">查看全部</text>
    </view>
    <scroll-view class="coach-list" scroll-x>
      <view 
        class="coach-item" 
        wx:for="{{coaches}}" 
        wx:key="id"
        bindtap="onViewCoachDetail"
        data-id="{{item.id}}"
      >
        <image class="coach-avatar" src="{{item.avatar}}"></image>
        <text class="coach-name">{{item.name}}</text>
        <view class="coach-rating">
          <text class="rating">{{item.rating}}</text>
          <text class="star">★</text>
        </view>
        <text class="coach-specialty">{{item.specialty}}</text>
      </view>
    </scroll-view>
  </view>

  <!-- 公告通知 -->
  <view class="notice-section" wx:if="{{notices.length > 0}}">
    <view class="section-title">场馆公告</view>
    <view class="notice-list">
      <view 
        class="notice-item" 
        wx:for="{{notices}}" 
        wx:key="id"
        bindtap="onViewNotice"
        data-id="{{item.id}}"
      >
        <text class="notice-title">{{item.title}}</text>
        <text class="notice-time">{{item.time}}</text>
      </view>
    </view>
  </view>
</view>
```

**B. 首页逻辑 (pages/index/index.js)**
```javascript
const userApi = require('../../api/user.js')
const bookingApi = require('../../api/booking.js')
const coachApi = require('../../api/coach.js')

Page({
  data: {
    banners: [],
    membershipCard: null,
    stats: {},
    upcomingBookings: [],
    coaches: [],
    notices: []
  },

  onLoad() {
    this.loadData()
  },

  onShow() {
    // 每次显示页面时刷新数据
    this.loadData()
  },

  // 下拉刷新
  onPullDownRefresh() {
    this.loadData().then(() => {
      wx.stopPullDownRefresh()
    })
  },

  // 加载数据
  async loadData() {
    try {
      const [
        membershipCard,
        stats,
        upcomingBookings,
        coaches,
        notices
      ] = await Promise.all([
        userApi.getMembershipCard(),
        userApi.getTrainingStats(),
        bookingApi.getUpcoming(),
        coachApi.getRecommended(),
        userApi.getNotices()
      ])

      this.setData({
        membershipCard,
        stats,
        upcomingBookings,
        coaches,
        notices
      })
    } catch (err) {
      console.error('加载数据失败', err)
    }
  },

  // 快速签到
  onQuickCheckin() {
    wx.navigateTo({
      url: '/pages/checkin/scan/scan'
    })
  },

  // 快速预约
  onQuickBooking() {
    wx.navigateTo({
      url: '/pages/booking/create/create'
    })
  },

  // 查看教练
  onViewCoaches() {
    wx.navigateTo({
      url: '/pages/coach/list/list'
    })
  },

  // 查看训练记录
  onViewRecords() {
    wx.navigateTo({
      url: '/pages/user/records/records'
    })
  },

  // 查看所有预约
  onViewAllBookings() {
    wx.switchTab({
      url: '/pages/booking/list/list'
    })
  },

  // 查看预约详情
  onViewBookingDetail(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/booking/detail/detail?id=${id}`
    })
  },

  // 查看所有教练
  onViewAllCoaches() {
    wx.navigateTo({
      url: '/pages/coach/list/list'
    })
  },

  // 查看教练详情
  onViewCoachDetail(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/coach/detail/detail?id=${id}`
    })
  },

  // 查看公告
  onViewNotice(e) {
    const id = e.currentTarget.dataset.id
    wx.navigateTo({
      url: `/pages/notice/detail/detail?id=${id}`
    })
  }
})
```

**验收标准**:
- 首页布局美观
- 数据加载正常
- 快捷功能可用
- 支持下拉刷新

---

### 2.4 课程预约功能

#### 2.4.1 预约流程实现
**输入**: 教练和时间选择
**输出**: 预约记录

**执行内容**:

**A. 创建预约页面 (pages/booking/create/create.wxml)**
```xml
<view class="container">
  <!-- 步骤指示器 -->
  <view class="steps">
    <view class="step {{currentStep >= 1 ? 'active' : ''}}">
      <text class="step-number">1</text>
      <text class="step-text">选择课程</text>
    </view>
    <view class="step {{currentStep >= 2 ? 'active' : ''}}">
      <text class="step-number">2</text>
      <text class="step-text">选择教练</text>
    </view>
    <view class="step {{currentStep >= 3 ? 'active' : ''}}">
      <text class="step-number">3</text>
      <text class="step-text">选择时间</text>
    </view>
    <view class="step {{currentStep >= 4 ? 'active' : ''}}">
      <text class="step-number">4</text>
      <text class="step-text">确认预约</text>
    </view>
  </view>

  <!-- 步骤1：选择课程 -->
  <view class="step-content" wx:if="{{currentStep === 1}}">
    <view class="course-list">
      <view 
        class="course-item {{selectedCourse.id === item.id ? 'selected' : ''}}" 
        wx:for="{{courses}}" 
        wx:key="id"
        bindtap="onSelectCourse"
        data-course="{{item}}"
      >
        <image class="course-icon" src="{{item.icon}}"></image>
        <view class="course-info">
          <text class="course-name">{{item.name}}</text>
          <text class="course-desc">{{item.description}}</text>
          <view class="course-meta">
            <text class="duration">{{item.duration}}分钟</text>
            <text class="price">¥{{item.price}}</text>
          </view>
        </view>
        <image class="check-icon" src="/images/icon-check.png" wx:if="{{selectedCourse.id === item.id}}"></image>
      </view>
    </view>
  </view>

  <!-- 步骤2：选择教练 -->
  <view class="step-content" wx:if="{{currentStep === 2}}">
    <view class="coach-list">
      <view 
        class="coach-item {{selectedCoach.id === item.id ? 'selected' : ''}}" 
        wx:for="{{coaches}}" 
        wx:key="id"
        bindtap="onSelectCoach"
        data-coach="{{item}}"
      >
        <image class="coach-avatar" src="{{item.avatar}}"></image>
        <view class="coach-info">
          <text class="coach-name">{{item.name}}</text>
          <view class="coach-rating">
            <text class="rating">{{item.rating}}</text>
            <text class="star">★</text>
            <text class="count">({{item.ratingCount}})</text>
          </view>
          <text class="coach-specialty">擅长：{{item.specialty}}</text>
        </view>
        <image class="check-icon" src="/images/icon-check.png" wx:if="{{selectedCoach.id === item.id}}"></image>
      </view>
    </view>
  </view>

  <!-- 步骤3：选择时间 -->
  <view class="step-content" wx:if="{{currentStep === 3}}">
    <!-- 日期选择 -->
    <scroll-view class="date-scroll" scroll-x>
      <view 
        class="date-item {{selectedDate === item.date ? 'selected' : ''}}" 
        wx:for="{{dates}}" 
        wx:key="date"
        bindtap="onSelectDate"
        data-date="{{item.date}}"
      >
        <text class="weekday">{{item.weekday}}</text>
        <text class="day">{{item.day}}</text>
        <text class="month">{{item.month}}</text>
      </view>
    </scroll-view>

    <!-- 时间段选择 -->
    <view class="time-slots">
      <view 
        class="time-slot {{item.isBooked ? 'disabled' : ''}} {{selectedTime === item.time ? 'selected' : ''}}" 
        wx:for="{{timeSlots}}" 
        wx:key="time"
        bindtap="onSelectTime"
        data-time="{{item.time}}"
      >
        <text>{{item.time}}</text>
        <text class="status" wx:if="{{item.isBooked}}">已约</text>
      </view>
    </view>
  </view>

  <!-- 步骤4：确认预约 -->
  <view class="step-content" wx:if="{{currentStep === 4}}">
    <view class="confirm-info">
      <view class="info-item">
        <text class="label">课程类型</text>
        <text class="value">{{selectedCourse.name}}</text>
      </view>
      <view class="info-item">
        <text class="label">教练</text>
        <text class="value">{{selectedCoach.name}}</text>
      </view>
      <view class="info-item">
        <text class="label">预约时间</text>
        <text class="value">{{selectedDate}} {{selectedTime}}</text>
      </view>
      <view class="info-item">
        <text class="label">课程时长</text>
        <text class="value">{{selectedCourse.duration}}分钟</text>
      </view>
      <view class="info-item">
        <text class="label">使用会员卡</text>
        <switch checked="{{useMembershipCard}}" bindchange="onToggleMembershipCard"/>
      </view>
    </view>

    <view class="remark-section">
      <text class="label">备注</text>
      <textarea 
        class="remark-input" 
        placeholder="请输入备注信息（选填）"
        value="{{remark}}"
        bindinput="onRemarkInput"
        maxlength="200"
      ></textarea>
    </view>
  </view>

  <!-- 底部按钮 -->
  <view class="bottom-bar">
    <button class="btn btn-back" bindtap="onPrevStep" wx:if="{{currentStep > 1}}">上一步</button>
    <button class="btn btn-next" bindtap="onNextStep" wx:if="{{currentStep < 4}}">下一步</button>
    <button class="btn btn-submit" bindtap="onSubmit" wx:if="{{currentStep === 4}}">确认预约</button>
  </view>
</view>
```

**B. 预约逻辑 (pages/booking/create/create.js)**
```javascript
const bookingApi = require('../../../api/booking.js')
const coachApi = require('../../../api/coach.js')

Page({
  data: {
    currentStep: 1,
    courses: [],
    coaches: [],
    dates: [],
    timeSlots: [],
    selectedCourse: {},
    selectedCoach: {},
    selectedDate: '',
    selectedTime: '',
    useMembershipCard: true,
    remark: ''
  },

  onLoad() {
    this.loadCourses()
  },

  // 加载课程列表
  async loadCourses() {
    try {
      const courses = await bookingApi.getCourseTypes()
      this.setData({ courses })
    } catch (err) {
      wx.showToast({
        title: '加载课程失败',
        icon: 'none'
      })
    }
  },

  // 选择课程
  onSelectCourse(e) {
    const course = e.currentTarget.dataset.course
    this.setData({ selectedCourse: course })
  },

  // 选择教练
  onSelectCoach(e) {
    const coach = e.currentTarget.dataset.coach
    this.setData({ selectedCoach: coach })
  },

  // 选择日期
  async onSelectDate(e) {
    const date = e.currentTarget.dataset.date
    this.setData({ selectedDate: date })
    
    // 加载该日期的可用时间段
    await this.loadTimeSlots(date)
  },

  // 加载时间段
  async loadTimeSlots(date) {
    try {
      wx.showLoading({ title: '加载中...' })
      const timeSlots = await coachApi.getAvailableSlots({
        coachId: this.data.selectedCoach.id,
        date: date
      })
      this.setData({ timeSlots })
      wx.hideLoading()
    } catch (err) {
      wx.hideLoading()
      wx.showToast({
        title: '加载时间段失败',
        icon: 'none'
      })
    }
  },

  // 选择时间
  onSelectTime(e) {
    const time = e.currentTarget.dataset.time
    this.setData({ selectedTime: time })
  },

  // 切换会员卡使用
  onToggleMembershipCard(e) {
    this.setData({
      useMembershipCard: e.detail.value
    })
  },

  // 备注输入
  onRemarkInput(e) {
    this.setData({
      remark: e.detail.value
    })
  },

  // 上一步
  onPrevStep() {
    this.setData({
      currentStep: this.data.currentStep - 1
    })
  },

  // 下一步
  async onNextStep() {
    const { currentStep, selectedCourse, selectedCoach, selectedDate, selectedTime } = this.data

    // 验证当前步骤
    if (currentStep === 1 && !selectedCourse.id) {
      wx.showToast({
        title: '请选择课程',
        icon: 'none'
      })
      return
    }

    if (currentStep === 2 && !selectedCoach.id) {
      wx.showToast({
        title: '请选择教练',
        icon: 'none'
      })
      return
    }

    // 进入步骤3时加载日期和教练
    if (currentStep === 2) {
      await this.loadDates()
      await this.loadCoaches()
    }

    if (currentStep === 3 && (!selectedDate || !selectedTime)) {
      wx.showToast({
        title: '请选择预约时间',
        icon: 'none'
      })
      return
    }

    this.setData({
      currentStep: currentStep + 1
    })
  },

  // 加载日期列表（未来7天）
  loadDates() {
    const dates = []
    const today = new Date()
    const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']

    for (let i = 0; i < 7; i++) {
      const date = new Date(today)
      date.setDate(today.getDate() + i)
      
      dates.push({
        date: this.formatDate(date),
        weekday: i === 0 ? '今天' : weekdays[date.getDay()],
        day: date.getDate(),
        month: date.getMonth() + 1
      })
    }

    this.setData({ dates })
  },

  // 加载教练列表
  async loadCoaches() {
    try {
      const coaches = await coachApi.getCoachesByCourse(this.data.selectedCourse.id)
      this.setData({ coaches })
    } catch (err) {
      wx.showToast({
        title: '加载教练失败',
        icon: 'none'
      })
    }
  },

  // 提交预约
  async onSubmit() {
    try {
      wx.showLoading({ title: '提交中...' })

      const data = {
        courseTypeId: this.data.selectedCourse.id,
        coachId: this.data.selectedCoach.id,
        bookingDate: this.data.selectedDate,
        startTime: this.data.selectedTime,
        useMembershipCard: this.data.useMembershipCard,
        remark: this.data.remark
      }

      await bookingApi.createBooking(data)

      wx.hideLoading()
      wx.showToast({
        title: '预约成功',
        icon: 'success'
      })

      // 跳转到预约列表
      setTimeout(() => {
        wx.switchTab({
          url: '/pages/booking/list/list'
        })
      }, 1500)

    } catch (err) {
      wx.hideLoading()
      wx.showToast({
        title: err.message || '预约失败',
        icon: 'none'
      })
    }
  },

  // 格式化日期
  formatDate(date) {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${year}-${month}-${day}`
  }
})
```

**验收标准**:
- 预约流程清晰
- 步骤切换流畅
- 数据验证完整
- 预约提交成功

---

### 2.5 签到功能

#### 2.5.1 扫码签到实现
**输入**: 二维码扫描
**输出**: 签到记录

**执行内容**:

**A. 扫码签到页面 (pages/checkin/scan/scan.wxml)**
```xml
<view class="container">
  <view class="scan-area">
    <view class="scan-tips">请将二维码放入框内</view>
    <view class="scan-frame"></view>
  </view>

  <view class="actions">
    <button class="btn-scan" bindtap="onScan">
      <image src="/images/icon-scan.png"></image>
      <text>扫码签到</text>
    </button>
    
    <button class="btn-face" bindtap="onFaceCheckin">
      <image src="/images/icon-face.png"></image>
      <text>人脸签到</text>
    </button>
  </view>

  <view class="today-record" wx:if="{{todayRecord}}">
    <view class="record-title">今日签到记录</view>
    <view class="record-info">
      <view class="record-item">
        <text class="label">签到时间</text>
        <text class="value">{{todayRecord.checkInTime}}</text>
      </view>
      <view class="record-item" wx:if="{{todayRecord.checkOutTime}}">
        <text class="label">签出时间</text>
        <text class="value">{{todayRecord.checkOutTime}}</text>
      </view>
      <view class="record-item" wx:if="{{todayRecord.duration}}">
        <text class="label">训练时长</text>
        <text class="value">{{todayRecord.duration}}分钟</text>
      </view>
    </view>
    
    <button 
      class="btn-checkout" 
      bindtap="onCheckout"
      wx:if="{{!todayRecord.checkOutTime}}"
    >
      签出
    </button>
  </view>

  <view class="stats-section">
    <view class="stat-item">
      <text class="stat-value">{{stats.totalDays}}</text>
      <text class="stat-label">累计训练天数</text>
    </view>
    <view class="stat-item">
      <text class="stat-value">{{stats.continuousDays}}</text>
      <text class="stat-label">连续训练天数</text>
    </view>
    <view class="stat-item">
      <text class="stat-value">{{stats.monthTimes}}</text>
      <text class="stat-label">本月训练次数</text>
    </view>
  </view>
</view>
```

**B. 签到逻辑 (pages/checkin/scan/scan.js)**
```javascript
const checkinApi = require('../../../api/checkin.js')

Page({
  data: {
    todayRecord: null,
    stats: {}
  },

  onShow() {
    this.loadData()
  },

  // 加载数据
  async loadData() {
    try {
      const [todayRecord, stats] = await Promise.all([
        checkinApi.getTodayRecord(),
        checkinApi.getStats()
      ])

      this.setData({
        todayRecord,
        stats
      })
    } catch (err) {
      console.error('加载数据失败', err)
    }
  },

  // 扫码签到
  onScan() {
    wx.scanCode({
      onlyFromCamera: true,
      scanType: ['qrCode'],
      success: (res) => {
        this.handleCheckin(res.result)
      },
      fail: (err) => {
        wx.showToast({
          title: '扫码失败',
          icon: 'none'
        })
      }
    })
  },

  // 处理签到
  async handleCheckin(qrCode) {
    try {
      wx.showLoading({ title: '签到中...' })

      await checkinApi.checkin({
        type: 'qrcode',
        qrCode: qrCode
      })

      wx.hideLoading()
      wx.showToast({
        title: '签到成功',
        icon: 'success'
      })

      // 刷新数据
      setTimeout(() => {
        this.loadData()
      }, 1500)

    } catch (err) {
      wx.hideLoading()
      wx.showToast({
        title: err.message || '签到失败',
        icon: 'none'
      })
    }
  },

  // 人脸签到
  onFaceCheckin() {
    wx.navigateTo({
      url: '/pages/checkin/face/face'
    })
  },

  // 签出
  async onCheckout() {
    try {
      wx.showLoading({ title: '签出中...' })

      await checkinApi.checkout()

      wx.hideLoading()
      wx.showToast({
        title: '签出成功',
        icon: 'success'
      })

      // 刷新数据
      setTimeout(() => {
        this.loadData()
      }, 1500)

    } catch (err) {
      wx.hideLoading()
      wx.showToast({
        title: err.message || '签出失败',
        icon: 'none'
      })
    }
  }
})
```

**验收标准**:
- 扫码功能正常
- 签到签出成功
- 数据统计准确
- 今日记录显示正确

---

## 三、API接口定义

### 3.1 用户相关接口

```javascript
// api/user.js
const request = require('../utils/request.js')

module.exports = {
  // 登录
  login(data) {
    return request.post('/api/v1/miniprogram/login', data)
  },

  // 获取用户信息
  getUserInfo() {
    return request.get('/api/v1/miniprogram/user/info')
  },

  // 获取会员卡信息
  getMembershipCard() {
    return request.get('/api/v1/miniprogram/user/membership-card')
  },

  // 获取训练统计
  getTrainingStats() {
    return request.get('/api/v1/miniprogram/user/stats')
  },

  // 获取公告列表
  getNotices() {
    return request.get('/api/v1/miniprogram/notices')
  }
}
```

### 3.2 预约相关接口

```javascript
// api/booking.js
const request = require('../utils/request.js')

module.exports = {
  // 获取课程类型
  getCourseTypes() {
    return request.get('/api/v1/miniprogram/course-types')
  },

  // 创建预约
  createBooking(data) {
    return request.post('/api/v1/miniprogram/bookings', data)
  },

  // 获取预约列表
  getBookings(params) {
    return request.get('/api/v1/miniprogram/bookings', params)
  },

  // 获取即将到来的预约
  getUpcoming() {
    return request.get('/api/v1/miniprogram/bookings/upcoming')
  },

  // 获取预约详情
  getBookingDetail(id) {
    return request.get(`/api/v1/miniprogram/bookings/${id}`)
  },

  // 取消预约
  cancelBooking(id, reason) {
    return request.put(`/api/v1/miniprogram/bookings/${id}/cancel`, { reason })
  },

  // 评价课程
  rateBooking(id, data) {
    return request.put(`/api/v1/miniprogram/bookings/${id}/rate`, data)
  }
}
```

### 3.3 签到相关接口

```javascript
// api/checkin.js
const request = require('../utils/request.js')

module.exports = {
  // 签到
  checkin(data) {
    return request.post('/api/v1/miniprogram/checkin', data)
  },

  // 签出
  checkout() {
    return request.post('/api/v1/miniprogram/checkout')
  },

  // 获取今日签到记录
  getTodayRecord() {
    return request.get('/api/v1/miniprogram/checkin/today')
  },

  // 获取签到统计
  getStats() {
    return request.get('/api/v1/miniprogram/checkin/stats')
  },

  // 获取签到记录列表
  getRecords(params) {
    return request.get('/api/v1/miniprogram/checkin/records', params)
  }
}
```

## 四、测试用例

### 4.1 功能测试
- 登录授权测试
- 首页数据加载测试
- 课程预约流程测试
- 签到签出测试
- 个人信息查看测试

### 4.2 兼容性测试
- iOS系统测试
- Android系统测试
- 不同微信版本测试
- 不同屏幕尺寸测试

### 4.3 性能测试
- 页面加载速度测试
- 网络请求性能测试
- 图片加载优化测试
- 内存占用测试

## 五、上线检查清单

### 5.1 开发检查
- [ ] 所有页面开发完成
- [ ] API接口联调完成
- [ ] 功能测试通过
- [ ] 兼容性测试通过
- [ ] 性能优化完成

### 5.2 配置检查
- [ ] AppID配置正确
- [ ] 服务器域名配置
- [ ] 业务域名配置
- [ ] 支付配置（如需要）
- [ ] 消息推送配置

### 5.3 审核准备
- [ ] 小程序信息完善
- [ ] 隐私政策完善
- [ ] 用户协议完善
- [ ] 测试账号准备
- [ ] 审核说明准备

## 六、后续优化方向

1. **社交功能**: 添加好友、训练打卡分享
2. **智能推荐**: 基于用户数据推荐课程和教练
3. **在线支付**: 支持会员卡购买和续费
4. **直播课程**: 支持线上直播课程
5. **积分系统**: 训练积分兑换奖励
6. **社区功能**: 用户交流社区

---
**任务优先级**: P0（核心功能）  
**预计工期**: 3-4周  
**依赖任务**: TASK007（系统基础设施）  
**后续任务**: TASK010（系统测试与上线）
