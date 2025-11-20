import React, { useEffect, useState } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Select } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { userService, User } from '../../services/userService'
import { useNavigate } from 'react-router-dom'

const UserList: React.FC = () => {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalVisible, setModalVisible] = useState(false)
  const [form] = Form.useForm()
  const navigate = useNavigate()

  const columns = [
    { title: '用户编号', dataIndex: 'user_no', key: 'user_no' },
    { title: '姓名', dataIndex: 'name', key: 'name' },
    { title: '手机号', dataIndex: 'phone', key: 'phone' },
    { 
      title: '状态', 
      dataIndex: 'status', 
      key: 'status',
      render: (status: number) => status === 1 ? '正常' : '冻结'
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: User) => (
        <Space>
          <Button type="link" onClick={() => handleEdit(record)}>编辑</Button>
          <Button type="link" onClick={() => handleView(record)}>查看</Button>
          <Button type="link" danger onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  const fetchUsers = async () => {
    setLoading(true)
    try {
      const response = await userService.list({ page, page_size: pageSize })
      setUsers(response.data.list || [])
      setTotal(response.data.total || 0)
    } catch (error) {
      message.error('获取用户列表失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchUsers()
  }, [page, pageSize])

  const handleAdd = () => {
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (record: User) => {
    form.setFieldsValue(record)
    setModalVisible(true)
  }

  const handleView = (record: User) => {
    navigate(`/users/${record.id}`)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个用户吗？',
      onOk: async () => {
        try {
          await userService.delete(id)
          message.success('删除成功')
          fetchUsers()
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
        await userService.update(values.id, values)
        message.success('更新成功')
      } else {
        await userService.create(values)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchUsers()
    } catch (error) {
      message.error('操作失败')
    }
  }

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          新增会员
        </Button>
      </div>
      <Table 
        columns={columns} 
        dataSource={users}
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
        title="用户信息"
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
        </Form>
      </Modal>
    </div>
  )
}

export default UserList