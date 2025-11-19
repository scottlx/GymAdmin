import React from 'react'
import { Table, Button, Space, Tag } from 'antd'
import { PlusOutlined } from '@ant-design/icons'

const CourseList: React.FC = () => {
  const columns = [
    { title: '课程名称', dataIndex: 'courseName', key: 'courseName' },
    { title: '教练', dataIndex: 'coachName', key: 'coachName' },
    { 
      title: '课程类型', 
      dataIndex: 'courseType', 
      key: 'courseType',
      render: (type: number) => type === 1 ? <Tag>私教课</Tag> : <Tag color="blue">团课</Tag>
    },
    { title: '开始时间', dataIndex: 'startTime', key: 'startTime' },
    { title: '人数', dataIndex: 'currentCount', key: 'currentCount', render: (count: number, record: any) => `${count}/${record.maxCapacity}` },
    {
      title: '操作',
      key: 'action',
      render: () => (
        <Space>
          <Button type="link">详情</Button>
          <Button type="link">编辑</Button>
          <Button type="link" danger>取消</Button>
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
          新增课程
        </Button>
      </div>
      <Table columns={columns} dataSource={data} />
    </div>
  )
}

export default CourseList
