import React from 'react'
import { Table, Button, Space, Tag } from 'antd'
import { PlusOutlined } from '@ant-design/icons'

const CardList: React.FC = () => {
  const columns = [
    { title: '卡号', dataIndex: 'cardNo', key: 'cardNo' },
    { title: '会员姓名', dataIndex: 'userName', key: 'userName' },
    { title: '卡类型', dataIndex: 'cardType', key: 'cardType' },
    { 
      title: '状态', 
      dataIndex: 'status', 
      key: 'status',
      render: (status: number) => {
        const statusMap: any = {
          1: <Tag color="green">正常</Tag>,
          2: <Tag color="red">已过期</Tag>,
          3: <Tag color="orange">已冻结</Tag>,
        }
        return statusMap[status] || '-'
      }
    },
    { title: '到期日期', dataIndex: 'endDate', key: 'endDate' },
    {
      title: '操作',
      key: 'action',
      render: () => (
        <Space>
          <Button type="link">详情</Button>
          <Button type="link">续费</Button>
          <Button type="link">冻结</Button>
        </Space>
      ),
    },
  ]

  // TODO: Fetch data from API
  const data: any[] = []

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />}>
          办理会员卡
        </Button>
      </div>
      <Table columns={columns} dataSource={data} />
    </div>
  )
}

export default CardList
