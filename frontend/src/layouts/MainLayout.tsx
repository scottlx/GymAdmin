import React from 'react'
import { Outlet, useNavigate } from 'react-router-dom'
import { Layout, Menu } from 'antd'
import {
  DashboardOutlined,
  UserOutlined,
  CreditCardOutlined,
  TeamOutlined,
  CalendarOutlined,
  CheckCircleOutlined,
  GiftOutlined,
} from '@ant-design/icons'

const { Header, Sider, Content } = Layout

const MainLayout: React.FC = () => {
  const navigate = useNavigate()

  const menuItems = [
    { key: 'dashboard', icon: <DashboardOutlined />, label: '仪表盘' },
    { key: 'users', icon: <UserOutlined />, label: '会员管理' },
    { key: 'cards', icon: <CreditCardOutlined />, label: '会员卡管理' },
    { key: 'coaches', icon: <TeamOutlined />, label: '教练管理' },
    { key: 'courses', icon: <CalendarOutlined />, label: '课程管理' },
    { key: 'checkins', icon: <CheckCircleOutlined />, label: '签到记录' },
    { key: 'vouchers', icon: <GiftOutlined />, label: '券核销' },
  ]

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="dark">
        <div style={{ height: 64, display: 'flex', alignItems: 'center', justifyContent: 'center', color: 'white', fontSize: 18, fontWeight: 'bold' }}>
          健身房管理系统
        </div>
        <Menu
          theme="dark"
          mode="inline"
          defaultSelectedKeys={['dashboard']}
          items={menuItems}
          onClick={({ key }) => navigate(`/${key}`)}
        />
      </Sider>
      <Layout>
        <Header style={{ background: '#fff', padding: '0 24px' }}>
          <h2 style={{ margin: 0 }}>健身房管理系统</h2>
        </Header>
        <Content style={{ margin: '24px 16px', padding: 24, background: '#fff' }}>
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  )
}

export default MainLayout
