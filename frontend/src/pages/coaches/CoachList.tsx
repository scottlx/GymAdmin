import React, { useEffect, useState } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Select } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { coachService, Coach } from '../../services/coachService'

const CoachList: React.FC = () => {
  const [coaches, setCoaches] = useState<Coach[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()

  const columns = [
    { title: '教练编号', dataIndex: 'coach_no', key: 'coach_no' },
    { title: '姓名', dataIndex: 'name', key: 'name' },
    { title: '手机号', dataIndex: 'phone', key: 'phone' },
    { title: '专长', dataIndex: 'specialties', key: 'specialties' },
    { title: '工作年限', dataIndex: 'experience', key: 'experience' },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Coach) => (
        <Space>
          <Button type="link" onClick={() => handleEdit(record)}>编辑</Button>
          <Button type="link" onClick={() => handleView(record)}>查看</Button>
          <Button type="link" danger onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  const fetchCoaches = async () => {
    setLoading(true)
    try {
      const response = await coachService.list({ page, page_size: pageSize })
      setCoaches(response.data.list || [])
      setTotal(response.data.total || 0)
    } catch (error) {
      message.error('获取教练列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchCoaches()
  }, [page, pageSize])

  const handleAdd = () => {
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (record: Coach) => {
    form.setFieldsValue(record)
    setModalVisible(true)
  }

  const handleView = (record: Coach) => {
    Modal.info({
      title: '教练详情',
      content: (
        <div>
          <p>姓名: {record.name}</p>
          <p>手机号: {record.phone}</p>
          <p>教练编号: {record.coach_no}</p>
          <p>工作年限: {record.experience}年</p>
        </div>
      ),
    })
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个教练吗？',
      onOk: async () => {
        try {
          await coachService.delete(id)
          message.success('删除成功')
          fetchCoaches()
        } catch (error) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      if (values.id) {
        await coachService.update(values.id, values)
        message.success('更新成功')
      } else {
        await coachService.create(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchCoaches()
    } catch (error) {
      message.error('操作失败')
    }
  }

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          新增教练
        </Button>
      </div>
      <Table 
        columns={columns} 
        dataSource={coaches}
        rowKey="id"
        loading={loading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: total,
          onChange: (p, ps) => {
            setPage(p)
            setPageSize(ps || 10)
          },
        }}
      />

      <Modal
        title="教练信息"
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="id" hidden>
            <Input />
          </Form.Item>
          <Form.Item name="name" label="姓名" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="phone" label="手机号" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="gender" label="性别">
            <Select>
              <Select.Option value={1}>男</Select.Option>
              <Select.Option value={2}>女</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="email" label="邮箱">
            <Input />
          </Form.Item>
          <Form.Item name="specialties" label="专长">
            <Input.TextArea />
          </Form.Item>
          <Form.Item name="experience" label="工作年限">
            <Input type="number" />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default CoachList
