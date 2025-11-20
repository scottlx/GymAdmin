import React, { useEffect, useState } from 'react'
import { Table, Button, Space, message, Modal, Form, Input, Select, Tag, InputNumber, DatePicker } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { courseService, Course } from '../../services/courseService'
import { coachService } from '../../services/coachService'
import dayjs from 'dayjs'

const CourseList: React.FC = () => {
  const [courses, setCourses] = useState<Course[]>([])
  const [loading, setLoading] = useState(false)
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [modalVisible, setModalVisible] = useState(false)
  const [coaches, setCoaches] = useState<any[]>([])
  const [form] = Form.useForm()

  const columns = [
    { title: '课程名称', dataIndex: 'course_name', key: 'course_name' },
    { title: '教练ID', dataIndex: 'coach_id', key: 'coach_id' },
    { 
      title: '课程类型', 
      dataIndex: 'course_type', 
      key: 'course_type',
      render: (type: number) => type === 1 ? <Tag>私教课</Tag> : <Tag color="blue">团课</Tag>
    },
    { 
      title: '开始时间', 
      dataIndex: 'start_time', 
      key: 'start_time',
      render: (time: string) => time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'
    },
    { 
      title: '结束时间', 
      dataIndex: 'end_time', 
      key: 'end_time',
      render: (time: string) => time ? dayjs(time).format('YYYY-MM-DD HH:mm') : '-'
    },
    { 
      title: '人数', 
      dataIndex: 'current_count', 
      key: 'current_count', 
      render: (count: number, record: Course) => `${count}/${record.max_capacity}` 
    },
    { 
      title: '价格', 
      dataIndex: 'price', 
      key: 'price',
      render: (price: number) => `¥${price}`
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status: number) => {
        const statusMap: any = {
          1: <Tag color="green">可预约</Tag>,
          2: <Tag color="orange">已满员</Tag>,
          3: <Tag color="red">已取消</Tag>,
          4: <Tag color="gray">已完成</Tag>,
        }
        return statusMap[status] || '-'
      }
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: Course) => (
        <Space>
          <Button type="link" onClick={() => handleEdit(record)}>编辑</Button>
          <Button type="link" onClick={() => handleView(record)}>查看</Button>
          <Button type="link" danger onClick={() => handleDelete(record.id)}>删除</Button>
        </Space>
      ),
    },
  ]

  const fetchCourses = async () => {
    setLoading(true)
    try {
      const response = await courseService.list({ page, page_size: pageSize })
      setCourses(response.data.list || [])
      setTotal(response.data.total || 0)
    } catch (error) {
      message.error('获取课程列表失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchCoaches = async () => {
    try {
      const response = await coachService.list({ page: 1, page_size: 100 })
      setCoaches(response.data.list || [])
    } catch (error) {
      console.error('获取教练列表失败')
    }
  }

  useEffect(() => {
    fetchCourses()
  }, [page, pageSize])

  useEffect(() => {
    fetchCoaches()
  }, [])

  const handleAdd = () => {
    form.resetFields()
    setModalVisible(true)
  }

  const handleEdit = (record: Course) => {
    form.setFieldsValue({
      ...record,
      start_time: record.start_time ? dayjs(record.start_time) : null,
      end_time: record.end_time ? dayjs(record.end_time) : null,
    })
    setModalVisible(true)
  }

  const handleView = (record: Course) => {
    Modal.info({
      title: '课程详情',
      width: 600,
      content: (
        <div>
          <p>课程名称: {record.course_name}</p>
          <p>教练ID: {record.coach_id}</p>
          <p>课程类型: {record.course_type === 1 ? '私教课' : '团课'}</p>
          <p>开始时间: {dayjs(record.start_time).format('YYYY-MM-DD HH:mm')}</p>
          <p>结束时间: {dayjs(record.end_time).format('YYYY-MM-DD HH:mm')}</p>
          <p>人数: {record.current_count}/{record.max_capacity}</p>
          <p>价格: ¥{record.price}</p>
          {record.description && <p>描述: {record.description}</p>}
          {record.remark && <p>备注: {record.remark}</p>}
        </div>
      ),
    })
  }

  const handleDelete = async (id: number) => {
    Modal.confirm({
      title: '确认删除',
      content: '确定要删除这个课程吗？',
      onOk: async () => {
        try {
          await courseService.delete(id)
          message.success('删除成功')
          fetchCourses()
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
        start_time: values.start_time ? values.start_time.toISOString() : undefined,
        end_time: values.end_time ? values.end_time.toISOString() : undefined,
      }
      
      if (values.id) {
        await courseService.update(values.id, submitData)
        message.success('更新成功')
      } else {
        await courseService.create(submitData)
        message.success('创建成功')
      }
      setModalVisible(false)
      fetchCourses()
    } catch (error) {
      message.error('操作失败')
    }
  }

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
          新增课程
        </Button>
      </div>
      <Table 
        columns={columns} 
        dataSource={courses}
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
        title="课程信息"
        open={modalVisible}
        onOk={handleSubmit}
        onCancel={() => setModalVisible(false)}
        width={600}
      >
        <Form form={form} layout="vertical">
          <Form.Item name="id" hidden>
            <Input />
          </Form.Item>
          <Form.Item name="course_name" label="课程名称" rules={[{ required: true, message: '请输入课程名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="coach_id" label="教练" rules={[{ required: true, message: '请选择教练' }]}>
            <Select
              showSearch
              placeholder="选择教练"
              optionFilterProp="children"
              filterOption={(input, option) =>
                (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
              }
              options={coaches.map(coach => ({
                value: coach.id,
                label: `${coach.name} (${coach.coach_no})`,
              }))}
            />
          </Form.Item>
          <Form.Item name="course_type" label="课程类型" rules={[{ required: true }]}>
            <Select>
              <Select.Option value={1}>私教课</Select.Option>
              <Select.Option value={2}>团课</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="start_time" label="开始时间" rules={[{ required: true }]}>
            <DatePicker showTime style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="end_time" label="结束时间" rules={[{ required: true }]}>
            <DatePicker showTime style={{ width: '100%' }} />
          </Form.Item>
          <Form.Item name="max_capacity" label="最大容量" rules={[{ required: true }]}>
            <InputNumber style={{ width: '100%' }} min={1} />
          </Form.Item>
          <Form.Item name="price" label="价格" rules={[{ required: true }]}>
            <InputNumber style={{ width: '100%' }} min={0} precision={2} />
          </Form.Item>
          <Form.Item name="status" label="状态">
            <Select>
              <Select.Option value={1}>可预约</Select.Option>
              <Select.Option value={2}>已满员</Select.Option>
              <Select.Option value={3}>已取消</Select.Option>
              <Select.Option value={4}>已完成</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={3} />
          </Form.Item>
          <Form.Item name="remark" label="备注">
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default CourseList
