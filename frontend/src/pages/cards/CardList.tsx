import React, { useEffect, useState } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Select, DatePicker, Tag, InputNumber } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { cardService, MembershipCard } from '../../services/cardService'
import { userService } from '../../services/userService'
import dayjs from 'dayjs'
import { useNavigate } from 'react-router-dom'

const CardList: React.FC = () => {
  const [cards, setCards] = useState<MembershipCard[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalVisible, setModalVisible] = useState(false)
  const [users, setUsers] = useState<any[]>([])
  const [form] = Form.useForm()
  const navigate = useNavigate()

  const columns = [
    { title: '卡号', dataIndex: 'card_no', key: 'card_no' },
    { title: '会员ID', dataIndex: 'user_id', key: 'user_id' },
    { 
      title: '状态', 
      dataIndex: 'status', 
      key: 'status',
      render: (status: number) => {
        const statusMap: any = {
          1: <Tag color="green">正常</Tag>,
          2: <Tag color="red">已过期</Tag>,
          3: <Tag color="orange">已冻结</Tag>,
          4: <Tag color="purple">已转出</Tag>,
          5: <Tag color="gray">已退卡</Tag>,
        }
        return statusMap[status] || '-'
      }
    },
    { 
      title: '开始日期', 
      dataIndex: 'start_date', 
      key: 'start_date',
      render: (date: string) => date ? dayjs(date).format('YYYY-MM-DD') : '-'
    },
    { 
      title: '到期日期', 
      dataIndex: 'end_date', 
      key: 'end_date',
      render: (date: string) => date ? dayjs(date).format('YYYY-MM-DD') : '-'
    },
    { 
      title: '购买价格', 
      dataIndex: 'purchase_price', 
      key: 'purchase_price',
      render: (price: number) => `¥${price}`
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: MembershipCard) => (
        <Space>
          <Button type="link" onClick={() => handleEdit(record)}>编辑</Button>
          <Button type="link" onClick={() => handleView(record)}>查看</Button>
          <Button type="link" danger onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  const fetchCards = async () => {
    setLoading(true)
    try {
      const response = await cardService.list({ page, page_size: pageSize })
      setCards(response.data.list || [])
      setTotal(response.data.total || 0)
    } catch (error) {
      message.error('获取会员卡列表失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchUsers = async () => {
    try {
      const response = await userService.list({ page: 1, page_size: 100 })
      setUsers(response.data.list || [])
    } catch (error) {
      console.error('获取用户列表失败')
    }
  }

  useEffect(() => {
    fetchCards()
  }, [page, pageSize])

  useEffect(() => {
    fetchUsers()
  }, [])

  const handleAdd = () => {
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (record: MembershipCard) => {
    form.setFieldsValue({
      ...record,
      start_date: record.start_date ? dayjs(record.start_date) : null,
      end_date: record.end_date ? dayjs(record.end_date) : null,
    })
    setModalVisible(true)
  }

  const handleView = (record: MembershipCard) => {
    navigate(`/cards/${record.id}`)
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这张会员卡吗？',
      onOk: async () => {
        try {
          await cardService.delete(id)
          message.success('删除成功')
          fetchCards()
        } catch (error) {
          message.error('删除失败')
        }
      },
    })
  }

  const handleSubmit = async () => {
    try {
      const values = await form.validateFields()
      const submitData = {
        ...values,
        start_date: values.start_date ? values.start_date.format('YYYY-MM-DD') : undefined,
        end_date: values.end_date ? values.end_date.format('YYYY-MM-DD') : undefined,
      }
      
      if (values.id) {
        await cardService.update(values.id, submitData)
        message.success('更新成功')
      } else {
        await cardService.create(submitData)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchCards()
    } catch (error) {
      message.error('操作失败')
    }
  }

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          办理会员卡
        </Button>
      </div>
      <Table 
        columns={columns} 
        dataSource={cards}
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
        title="会员卡信息"
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="id" hidden>
            <Input />
          </Form.Item>
          <Form.Item name="user_id" label="会员" rules={[{ required: true, message: '请选择会员' }]}>
            <Select
              showSearch
              placeholder="选择会员"
              optionFilterProp="children"
              filterOption={(input, option) =>
                (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
              }
              options={users.map(user => ({
                value: user.id,
                label: `${user.name} (${user.phone})`,
              }))}
            />
          </Form.Item>
          <Form.Item name="card_type_id" label="卡类型ID" rules={[{ required: true }]}>
            <InputNumber style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="start_date" label="开始日期" rules={[{ required: true }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="end_date" label="到期日期" rules={[{ required: true }]}>
            <DatePicker style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="purchase_price" label="购买价格" rules={[{ required: true }]}>
            <InputNumber style={{ width: '100%' }} min={0} precision={2} />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select>
              <Select.Option value={1}>正常</Select.Option>
              <Select.Option value={2}>已过期</Select.Option>
              <Select.Option value={3}>已冻结</Select.Option>
              <Select.Option value={4}>已转出</Select.Option>
              <Select.Option value={5}>已退卡</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="remark" label="备注">
            <Input.TextArea rows={3} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default CardList