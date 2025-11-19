# TASK009 - 数据统计与报表

## 一、功能概述

### 1.1 功能目标
构建健身房管理系统的数据统计与报表模块，提供全面的数据分析和可视化展示，帮助管理者了解运营状况、做出科学决策，支持多维度数据统计和报表导出。

### 1.2 核心价值
- 提供实时运营数据监控和分析
- 支持多维度数据统计和对比
- 生成可视化图表，直观展示数据趋势
- 导出各类报表，支持决策分析
- 发现运营问题，优化业务流程

### 1.3 涉及角色
- **管理员**: 查看所有统计数据和报表
- **教练**: 查看个人业绩和学员数据
- **前台人员**: 查看会员和预约数据
- **财务人员**: 查看收入和财务报表

## 二、功能详细拆解

### 2.1 运营概览仪表板

#### 2.1.1 核心指标统计
**输入**: 时间范围
**输出**: 核心运营指标

**执行内容**:

**A. 数据库统计查询**
```sql
-- 核心指标统计视图
CREATE VIEW v_dashboard_stats AS
SELECT
    -- 会员统计
    (SELECT COUNT(*) FROM users WHERE status = 1) AS total_members,
    (SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURDATE()) AS today_new_members,
    (SELECT COUNT(*) FROM users WHERE MONTH(created_at) = MONTH(CURDATE())) AS month_new_members,
    
    -- 会员卡统计
    (SELECT COUNT(*) FROM membership_cards WHERE status = 1) AS active_cards,
    (SELECT COUNT(*) FROM membership_cards WHERE expire_date < CURDATE()) AS expired_cards,
    (SELECT COUNT(*) FROM membership_cards WHERE DATEDIFF(expire_date, CURDATE()) <= 7) AS expiring_soon_cards,
    
    -- 签到统计
    (SELECT COUNT(*) FROM check_in_records WHERE DATE(check_in_time) = CURDATE()) AS today_checkins,
    (SELECT COUNT(*) FROM check_in_records WHERE MONTH(check_in_time) = MONTH(CURDATE())) AS month_checkins,
    
    -- 预约统计
    (SELECT COUNT(*) FROM course_bookings WHERE status IN (1,2)) AS active_bookings,
    (SELECT COUNT(*) FROM course_bookings WHERE DATE(booking_date) = CURDATE()) AS today_bookings,
    (SELECT COUNT(*) FROM course_bookings WHERE status = 3) AS completed_bookings,
    
    -- 收入统计
    (SELECT COALESCE(SUM(amount), 0) FROM orders WHERE DATE(created_at) = CURDATE()) AS today_revenue,
    (SELECT COALESCE(SUM(amount), 0) FROM orders WHERE MONTH(created_at) = MONTH(CURDATE())) AS month_revenue,
    
    -- 教练统计
    (SELECT COUNT(*) FROM coaches WHERE status = 1) AS active_coaches,
    (SELECT AVG(rating) FROM coaches WHERE status = 1) AS avg_coach_rating;
```

**B. 后端API实现 (services/dashboard_service.go)**
```go
package service

type DashboardStats struct {
    // 会员统计
    TotalMembers      int     `json:"total_members"`
    TodayNewMembers   int     `json:"today_new_members"`
    MonthNewMembers   int     `json:"month_new_members"`
    MemberGrowthRate  float64 `json:"member_growth_rate"`
    
    // 会员卡统计
    ActiveCards       int     `json:"active_cards"`
    ExpiredCards      int     `json:"expired_cards"`
    ExpiringSoonCards int     `json:"expiring_soon_cards"`
    
    // 签到统计
    TodayCheckins     int     `json:"today_checkins"`
    MonthCheckins     int     `json:"month_checkins"`
    CheckinRate       float64 `json:"checkin_rate"`
    
    // 预约统计
    ActiveBookings    int     `json:"active_bookings"`
    TodayBookings     int     `json:"today_bookings"`
    CompletedBookings int     `json:"completed_bookings"`
    CompletionRate    float64 `json:"completion_rate"`
    
    // 收入统计
    TodayRevenue      float64 `json:"today_revenue"`
    MonthRevenue      float64 `json:"month_revenue"`
    RevenueGrowthRate float64 `json:"revenue_growth_rate"`
    
    // 教练统计
    ActiveCoaches     int     `json:"active_coaches"`
    AvgCoachRating    float64 `json:"avg_coach_rating"`
}

func GetDashboardStats() (*DashboardStats, error) {
    var stats DashboardStats
    
    // 查询核心指标
    err := database.DB.Raw("SELECT * FROM v_dashboard_stats").Scan(&stats).Error
    if err != nil {
        return nil, err
    }
    
    // 计算增长率
    stats.MemberGrowthRate = calculateGrowthRate("users", "month")
    stats.RevenueGrowthRate = calculateGrowthRate("orders", "month")
    stats.CheckinRate = calculateCheckinRate()
    stats.CompletionRate = calculateCompletionRate()
    
    return &stats, nil
}

// 计算增长率
func calculateGrowthRate(table string, period string) float64 {
    var current, previous int
    
    // 当前周期数据
    database.DB.Raw(fmt.Sprintf(
        "SELECT COUNT(*) FROM %s WHERE MONTH(created_at) = MONTH(CURDATE())",
        table,
    )).Scan(&current)
    
    // 上一周期数据
    database.DB.Raw(fmt.Sprintf(
        "SELECT COUNT(*) FROM %s WHERE MONTH(created_at) = MONTH(DATE_SUB(CURDATE(), INTERVAL 1 MONTH))",
        table,
    )).Scan(&previous)
    
    if previous == 0 {
        return 0
    }
    
    return float64(current-previous) / float64(previous) * 100
}

// 计算签到率
func calculateCheckinRate() float64 {
    var activeMembers, checkins int
    
    database.DB.Raw("SELECT COUNT(*) FROM users WHERE status = 1").Scan(&activeMembers)
    database.DB.Raw("SELECT COUNT(DISTINCT user_id) FROM check_in_records WHERE MONTH(check_in_time) = MONTH(CURDATE())").Scan(&checkins)
    
    if activeMembers == 0 {
        return 0
    }
    
    return float64(checkins) / float64(activeMembers) * 100
}

// 计算完成率
func calculateCompletionRate() float64 {
    var total, completed int
    
    database.DB.Raw("SELECT COUNT(*) FROM course_bookings WHERE MONTH(booking_date) = MONTH(CURDATE())").Scan(&total)
    database.DB.Raw("SELECT COUNT(*) FROM course_bookings WHERE status = 3 AND MONTH(booking_date) = MONTH(CURDATE())").Scan(&completed)
    
    if total == 0 {
        return 0
    }
    
    return float64(completed) / float64(total) * 100
}
```

**C. 前端仪表板页面 (src/pages/Dashboard/Dashboard.tsx)**
```typescript
import React, { useEffect, useState } from 'react'
import { Card, Row, Col, Statistic, Progress } from 'antd'
import { ArrowUpOutlined, ArrowDownOutlined } from '@ant-design/icons'
import { getDashboardStats } from '@/api/dashboard'

const Dashboard: React.FC = () => {
  const [stats, setStats] = useState<any>({})
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadStats()
  }, [])

  const loadStats = async () => {
    try {
      const data = await getDashboardStats()
      setStats(data)
    } catch (error) {
      console.error('加载统计数据失败', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="dashboard">
      {/* 核心指标卡片 */}
      <Row gutter={16}>
        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="总会员数"
              value={stats.total_members}
              prefix={<UserOutlined />}
              suffix={
                <span className="growth-rate">
                  {stats.member_growth_rate > 0 ? (
                    <ArrowUpOutlined style={{ color: '#3f8600' }} />
                  ) : (
                    <ArrowDownOutlined style={{ color: '#cf1322' }} />
                  )}
                  {Math.abs(stats.member_growth_rate).toFixed(1)}%
                </span>
              }
            />
            <div className="stat-detail">
              今日新增: {stats.today_new_members} | 本月新增: {stats.month_new_members}
            </div>
          </Card>
        </Col>

        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="今日签到"
              value={stats.today_checkins}
              prefix={<CheckCircleOutlined />}
            />
            <div className="stat-detail">
              本月签到: {stats.month_checkins}
              <Progress 
                percent={stats.checkin_rate} 
                size="small" 
                status="active"
              />
            </div>
          </Card>
        </Col>

        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="今日收入"
              value={stats.today_revenue}
              precision={2}
              prefix="¥"
              suffix={
                <span className="growth-rate">
                  {stats.revenue_growth_rate > 0 ? (
                    <ArrowUpOutlined style={{ color: '#3f8600' }} />
                  ) : (
                    <ArrowDownOutlined style={{ color: '#cf1322' }} />
                  )}
                  {Math.abs(stats.revenue_growth_rate).toFixed(1)}%
                </span>
              }
            />
            <div className="stat-detail">
              本月收入: ¥{stats.month_revenue?.toFixed(2)}
            </div>
          </Card>
        </Col>

        <Col span={6}>
          <Card loading={loading}>
            <Statistic
              title="活跃预约"
              value={stats.active_bookings}
              prefix={<CalendarOutlined />}
            />
            <div className="stat-detail">
              今日预约: {stats.today_bookings}
              <Progress 
                percent={stats.completion_rate} 
                size="small"
                status="active"
              />
            </div>
          </Card>
        </Col>
      </Row>

      {/* 图表区域 */}
      <Row gutter={16} style={{ marginTop: 16 }}>
        <Col span={12}>
          <Card title="会员增长趋势" loading={loading}>
            <MemberTrendChart />
          </Card>
        </Col>
        <Col span={12}>
          <Card title="收入趋势" loading={loading}>
            <RevenueTrendChart />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginTop: 16 }}>
        <Col span={12}>
          <Card title="签到热力图" loading={loading}>
            <CheckinHeatmap />
          </Card>
        </Col>
        <Col span={12}>
          <Card title="课程预约分布" loading={loading}>
            <BookingDistributionChart />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default Dashboard
```

**验收标准**:
- 核心指标统计准确
- 增长率计算正确
- 数据实时更新
- 图表展示清晰

---

### 2.2 会员数据分析

#### 2.2.1 会员统计报表
**输入**: 筛选条件、时间范围
**输出**: 会员统计数据和图表

**执行内容**:

**A. 会员统计API (services/member_stats_service.go)**
```go
// 会员来源分析
func GetMemberSourceStats(startDate, endDate time.Time) ([]SourceStat, error) {
    var stats []SourceStat
    
    err := database.DB.Raw(`
        SELECT 
            source,
            COUNT(*) as count,
            COUNT(*) * 100.0 / (SELECT COUNT(*) FROM users WHERE created_at BETWEEN ? AND ?) as percentage
        FROM users
        WHERE created_at BETWEEN ? AND ?
        GROUP BY source
        ORDER BY count DESC
    `, startDate, endDate, startDate, endDate).Scan(&stats).Error
    
    return stats, err
}

// 会员活跃度分析
func GetMemberActivityStats(startDate, endDate time.Time) (*ActivityStats, error) {
    var stats ActivityStats
    
    // 活跃会员（本月有签到记录）
    database.DB.Raw(`
        SELECT COUNT(DISTINCT user_id) 
        FROM check_in_records 
        WHERE check_in_time BETWEEN ? AND ?
    `, startDate, endDate).Scan(&stats.ActiveMembers)
    
    // 沉睡会员（超过30天未签到）
    database.DB.Raw(`
        SELECT COUNT(*) 
        FROM users u
        WHERE status = 1
        AND NOT EXISTS (
            SELECT 1 FROM check_in_records 
            WHERE user_id = u.id 
            AND check_in_time > DATE_SUB(NOW(), INTERVAL 30 DAY)
        )
    `).Scan(&stats.InactiveMembers)
    
    // 流失会员（超过90天未签到）
    database.DB.Raw(`
        SELECT COUNT(*) 
        FROM users u
        WHERE status = 1
        AND NOT EXISTS (
            SELECT 1 FROM check_in_records 
            WHERE user_id = u.id 
            AND check_in_time > DATE_SUB(NOW(), INTERVAL 90 DAY)
        )
    `).Scan(&stats.ChurnedMembers)
    
    return &stats, nil
}

// 会员年龄分布
func GetMemberAgeDistribution() ([]AgeDistribution, error) {
    var distribution []AgeDistribution
    
    err := database.DB.Raw(`
        SELECT 
            CASE 
                WHEN TIMESTAMPDIFF(YEAR, birthday, CURDATE()) < 20 THEN '20岁以下'
                WHEN TIMESTAMPDIFF(YEAR, birthday, CURDATE()) BETWEEN 20 AND 29 THEN '20-29岁'
                WHEN TIMESTAMPDIFF(YEAR, birthday, CURDATE()) BETWEEN 30 AND 39 THEN '30-39岁'
                WHEN TIMESTAMPDIFF(YEAR, birthday, CURDATE()) BETWEEN 40 AND 49 THEN '40-49岁'
                ELSE '50岁以上'
            END as age_range,
            COUNT(*) as count
        FROM users
        WHERE birthday IS NOT NULL
        GROUP BY age_range
        ORDER BY age_range
    `).Scan(&distribution).Error
    
    return distribution, err
}

// 会员性别分布
func GetMemberGenderDistribution() ([]GenderDistribution, error) {
    var distribution []GenderDistribution
    
    err := database.DB.Raw(`
        SELECT 
            CASE gender
                WHEN 1 THEN '男'
                WHEN 2 THEN '女'
                ELSE '未知'
            END as gender,
            COUNT(*) as count
        FROM users
        GROUP BY gender
    `).Scan(&distribution).Error
    
    return distribution, err
}
```

**B. 前端会员分析页面 (src/pages/Stats/MemberStats.tsx)**
```typescript
import React, { useEffect, useState } from 'react'
import { Card, Row, Col, DatePicker, Select } from 'antd'
import { Pie, Column, Line } from '@ant-design/charts'
import { getMemberStats } from '@/api/stats'

const MemberStats: React.FC = () => {
  const [dateRange, setDateRange] = useState<any>([])
  const [sourceData, setSourceData] = useState([])
  const [activityData, setActivityData] = useState({})
  const [ageData, setAgeData] = useState([])
  const [genderData, setGenderData] = useState([])

  useEffect(() => {
    loadData()
  }, [dateRange])

  const loadData = async () => {
    try {
      const data = await getMemberStats({
        startDate: dateRange[0],
        endDate: dateRange[1]
      })
      
      setSourceData(data.source)
      setActivityData(data.activity)
      setAgeData(data.age)
      setGenderData(data.gender)
    } catch (error) {
      console.error('加载数据失败', error)
    }
  }

  // 来源分布饼图配置
  const sourceConfig = {
    data: sourceData,
    angleField: 'count',
    colorField: 'source',
    radius: 0.8,
    label: {
      type: 'outer',
      content: '{name} {percentage}'
    }
  }

  // 年龄分布柱状图配置
  const ageConfig = {
    data: ageData,
    xField: 'age_range',
    yField: 'count',
    label: {
      position: 'top',
      style: {
        fill: '#000000',
        opacity: 0.6
      }
    }
  }

  return (
    <div className="member-stats">
      <Card>
        <DatePicker.RangePicker 
          onChange={setDateRange}
          style={{ marginBottom: 16 }}
        />
      </Card>

      <Row gutter={16}>
        <Col span={8}>
          <Card title="会员来源分布">
            <Pie {...sourceConfig} />
          </Card>
        </Col>
        <Col span={8}>
          <Card title="年龄分布">
            <Column {...ageConfig} />
          </Card>
        </Col>
        <Col span={8}>
          <Card title="性别分布">
            <Pie 
              data={genderData}
              angleField="count"
              colorField="gender"
            />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginTop: 16 }}>
        <Col span={24}>
          <Card title="会员活跃度分析">
            <Row gutter={16}>
              <Col span={8}>
                <Statistic 
                  title="活跃会员" 
                  value={activityData.active_members}
                  valueStyle={{ color: '#3f8600' }}
                />
              </Col>
              <Col span={8}>
                <Statistic 
                  title="沉睡会员" 
                  value={activityData.inactive_members}
                  valueStyle={{ color: '#faad14' }}
                />
              </Col>
              <Col span={8}>
                <Statistic 
                  title="流失会员" 
                  value={activityData.churned_members}
                  valueStyle={{ color: '#cf1322' }}
                />
              </Col>
            </Row>
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default MemberStats
```

**验收标准**:
- 会员统计数据准确
- 图表展示清晰
- 支持时间范围筛选
- 数据可导出

---

### 2.3 财务报表

#### 2.3.1 收入统计
**输入**: 时间范围、统计维度
**输出**: 收入报表和图表

**执行内容**:

**A. 财务统计API (services/finance_stats_service.go)**
```go
// 收入统计
func GetRevenueStats(startDate, endDate time.Time, dimension string) ([]RevenueStat, error) {
    var stats []RevenueStat
    var groupBy string
    
    switch dimension {
    case "day":
        groupBy = "DATE(created_at)"
    case "week":
        groupBy = "YEARWEEK(created_at)"
    case "month":
        groupBy = "DATE_FORMAT(created_at, '%Y-%m')"
    default:
        groupBy = "DATE(created_at)"
    }
    
    err := database.DB.Raw(fmt.Sprintf(`
        SELECT 
            %s as period,
            SUM(amount) as total_amount,
            COUNT(*) as order_count,
            AVG(amount) as avg_amount
        FROM orders
        WHERE created_at BETWEEN ? AND ?
        AND status = 'paid'
        GROUP BY %s
        ORDER BY period
    `, groupBy, groupBy), startDate, endDate).Scan(&stats).Error
    
    return stats, err
}

// 收入来源分析
func GetRevenueSourceStats(startDate, endDate time.Time) ([]RevenueSource, error) {
    var stats []RevenueSource
    
    err := database.DB.Raw(`
        SELECT 
            order_type,
            SUM(amount) as amount,
            COUNT(*) as count
        FROM orders
        WHERE created_at BETWEEN ? AND ?
        AND status = 'paid'
        GROUP BY order_type
    `, startDate, endDate).Scan(&stats).Error
    
    return stats, err
}

// 会员卡销售统计
func GetMembershipCardSalesStats(startDate, endDate time.Time) (*CardSalesStats, error) {
    var stats CardSalesStats
    
    // 新办卡数量和金额
    database.DB.Raw(`
        SELECT 
            COUNT(*) as new_cards,
            SUM(price) as new_cards_revenue
        FROM membership_cards
        WHERE created_at BETWEEN ? AND ?
    `, startDate, endDate).Scan(&stats)
    
    // 续卡数量和金额
    database.DB.Raw(`
        SELECT 
            COUNT(*) as renewal_cards,
            SUM(price) as renewal_revenue
        FROM membership_card_renewals
        WHERE created_at BETWEEN ? AND ?
    `, startDate, endDate).Scan(&stats)
    
    return &stats, nil
}

// 导出财务报表
func ExportFinanceReport(startDate, endDate time.Time) (string, error) {
    // 创建Excel文件
    f := excelize.NewFile()
    
    // 收入汇总表
    sheet1 := "收入汇总"
    f.NewSheet(sheet1)
    
    // 写入表头
    headers := []string{"日期", "订单数", "总收入", "平均订单金额"}
    for i, header := range headers {
        cell := fmt.Sprintf("%s1", string(rune('A'+i)))
        f.SetCellValue(sheet1, cell, header)
    }
    
    // 查询数据
    stats, _ := GetRevenueStats(startDate, endDate, "day")
    
    // 写入数据
    for i, stat := range stats {
        row := i + 2
        f.SetCellValue(sheet1, fmt.Sprintf("A%d", row), stat.Period)
        f.SetCellValue(sheet1, fmt.Sprintf("B%d", row), stat.OrderCount)
        f.SetCellValue(sheet1, fmt.Sprintf("C%d", row), stat.TotalAmount)
        f.SetCellValue(sheet1, fmt.Sprintf("D%d", row), stat.AvgAmount)
    }
    
    // 保存文件
    filename := fmt.Sprintf("finance_report_%s_%s.xlsx", 
        startDate.Format("20060102"), 
        endDate.Format("20060102"))
    filepath := fmt.Sprintf("./exports/%s", filename)
    
    if err := f.SaveAs(filepath); err != nil {
        return "", err
    }
    
    return filepath, nil
}
```

**B. 前端财务报表页面 (src/pages/Stats/FinanceStats.tsx)**
```typescript
import React, { useEffect, useState } from 'react'
import { Card, Row, Col, DatePicker, Button, Table } from 'antd'
import { Line, Pie } from '@ant-design/charts'
import { getFinanceStats, exportFinanceReport } from '@/api/stats'

const FinanceStats: React.FC = () => {
  const [dateRange, setDateRange] = useState<any>([])
  const [revenueData, setRevenueData] = useState([])
  const [sourceData, setSourceData] = useState([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    loadData()
  }, [dateRange])

  const loadData = async () => {
    setLoading(true)
    try {
      const data = await getFinanceStats({
        startDate: dateRange[0],
        endDate: dateRange[1]
      })
      
      setRevenueData(data.revenue)
      setSourceData(data.source)
    } catch (error) {
      console.error('加载数据失败', error)
    } finally {
      setLoading(false)
    }
  }

  const handleExport = async () => {
    try {
      const result = await exportFinanceReport({
        startDate: dateRange[0],
        endDate: dateRange[1]
      })
      
      // 下载文件
      window.open(result.download_url)
    } catch (error) {
      console.error('导出失败', error)
    }
  }

  // 收入趋势图配置
  const revenueConfig = {
    data: revenueData,
    xField: 'period',
    yField: 'total_amount',
    point: {
      size: 5,
      shape: 'diamond'
    },
    label: {
      style: {
        fill: '#aaa'
      }
    }
  }

  // 收入来源饼图配置
  const sourceConfig = {
    data: sourceData,
    angleField: 'amount',
    colorField: 'order_type',
    radius: 0.8,
    label: {
      type: 'outer',
      content: '{name} ¥{value}'
    }
  }

  const columns = [
    { title: '日期', dataIndex: 'period', key: 'period' },
    { title: '订单数', dataIndex: 'order_count', key: 'order_count' },
    { 
      title: '总收入', 
      dataIndex: 'total_amount', 
      key: 'total_amount',
      render: (val: number) => `¥${val.toFixed(2)}`
    },
    { 
      title: '平均订单金额', 
      dataIndex: 'avg_amount', 
      key: 'avg_amount',
      render: (val: number) => `¥${val.toFixed(2)}`
    }
  ]

  return (
    <div className="finance-stats">
      <Card>
        <Row justify="space-between" align="middle">
          <Col>
            <DatePicker.RangePicker 
              onChange={setDateRange}
            />
          </Col>
          <Col>
            <Button type="primary" onClick={handleExport}>
              导出报表
            </Button>
          </Col>
        </Row>
      </Card>

      <Row gutter={16} style={{ marginTop: 16 }}>
        <Col span={16}>
          <Card title="收入趋势" loading={loading}>
            <Line {...revenueConfig} />
          </Card>
        </Col>
        <Col span={8}>
          <Card title="收入来源" loading={loading}>
            <Pie {...sourceConfig} />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginTop: 16 }}>
        <Col span={24}>
          <Card title="收入明细" loading={loading}>
            <Table 
              columns={columns}
              dataSource={revenueData}
              pagination={{ pageSize: 10 }}
            />
          </Card>
        </Col>
      </Row>
    </div>
  )
}

export default FinanceStats
```

**验收标准**:
- 财务数据统计准确
- 支持多维度分析
- 报表导出成功
- 数据可视化清晰

---

### 2.4 教练业绩报表

#### 2.4.1 教练业绩统计
**输入**: 教练ID、时间范围
**输出**: 教练业绩数据

**执行内容**:

**A. 教练业绩统计API**
```go
// 教练业绩排行
func GetCoachPerformanceRanking(startDate, endDate time.Time, rankBy string) ([]CoachRanking, error) {
    var rankings []CoachRanking
    var orderBy string
    
    switch rankBy {
    case "bookings":
        orderBy = "total_bookings DESC"
    case "revenue":
        orderBy = "total_revenue DESC"
    case "rating":
        orderBy = "avg_rating DESC"
    default:
        orderBy = "total_bookings DESC"
    }
    
    err := database.DB.Raw(fmt.Sprintf(`
        SELECT 
            c.id,
            c.real_name as name,
            c.avatar_url,
            COUNT(cb.id) as total_bookings,
            SUM(CASE WHEN cb.status = 3 THEN 1 ELSE 0 END) as completed_bookings,
            AVG(cb.rating) as avg_rating,
            SUM(ct.price) as total_revenue
        FROM coaches c
        LEFT JOIN course_bookings cb ON c.id = cb.coach_id
        LEFT JOIN course_types ct ON cb.course_type_id = ct.id
        WHERE cb.booking_date BETWEEN ? AND ?
        GROUP BY c.id
        ORDER BY %s
        LIMIT 10
    `, orderBy), startDate, endDate).Scan(&rankings).Error
    
    return rankings, err
}

// 教练课程完成率
func GetCoachCompletionRate(coachID int64, startDate, endDate time.Time) (float64, error) {
    var total, completed int
    
    database.DB.Raw(`
        SELECT COUNT(*) 
        FROM course_bookings 
        WHERE coach_id = ? 
        AND booking_date BETWEEN ? AND ?
    `, coachID, startDate, endDate).Scan(&total)
    
    database.DB.Raw(`
        SELECT COUNT(*) 
        FROM course_bookings 
        WHERE coach_id = ? 
        AND booking_date BETWEEN ? AND ?
        AND status = 3
    `, coachID, startDate, endDate).Scan(&completed)
    
    if total == 0 {
        return 0, nil
    }
    
    return float64(completed) / float64(total) * 100, nil
}
```

**验收标准**:
- 教练业绩统计准确
- 排行榜实时更新
- 支持多维度排序
- 数据可导出

---

## 三、API接口文档

### 3.1 统计相关接口

| 接口路径 | 方法 | 说明 | 权限 |
|---------|------|------|------|
| /api/v1/stats/dashboard | GET | 获取仪表板统计 | 管理员 |
| /api/v1/stats/members | GET | 获取会员统计 | 管理员 |
| /api/v1/stats/finance | GET | 获取财务统计 | 管理员、财务 |
| /api/v1/stats/coaches | GET | 获取教练统计 | 管理员 |
| /api/v1/stats/bookings | GET | 获取预约统计 | 管理员 |
| /api/v1/stats/export | POST | 导出报表 | 管理员 |

## 四、测试用例

### 4.1 功能测试
- 统计数据准确性测试
- 图表展示测试
- 报表导出测试
- 时间范围筛选测试

### 4.2 性能测试
- 大数据量统计性能测试
- 并发查询测试
- 报表生成速度测试

## 五、上线检查清单

### 5.1 开发检查
- [ ] 所有统计功能开发完成
- [ ] 数据准确性验证通过
- [ ] 图表展示正常
- [ ] 报表导出功能正常

### 5.2 性能检查
- [ ] 统计查询优化完成
- [ ] 索引创建完成
- [ ] 缓存策略配置
- [ ] 大数据量测试通过

## 六、后续优化方向

1. **实时数据**: 引入实时数据流处理
2. **预测分析**: 基于历史数据进行趋势预测
3. **自定义报表**: 支持用户自定义报表模板
4. **数据对比**: 支持同比、环比数据对比
5. **智能告警**: 异常数据自动告警
6. **移动端**: 开发移动端数据看板

---
**任务优先级**: P1（重要功能）  
**预计工期**: 2-3周  
**依赖任务**: TASK001-TASK006（业务功能）  
**后续任务**: TASK010（系统测试与上线）
