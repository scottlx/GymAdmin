import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { userService } from '../../../services/userService';
import { message, Spin, Descriptions, Row, Col } from 'antd';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

const UserStats: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [stats, setStats] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const response = await userService.getStats(Number(id));
        setStats(response.data);
      } catch (error) {
        message.error('获取用户统计失败');
      } finally {
        setLoading(false);
      }
    };
    fetchStats();
  }, [id]);

  if (loading) {
    return <Spin />;
  }

  if (!stats) {
    return <div>暂无统计数据</div>;
  }

  const chartData = [
    {
      name: '最近7天',
      训练次数: stats.last_7_days_checkin_count,
    },
    {
      name: '最近30天',
      训练次数: stats.last_30_days_checkin_count,
    },
  ];

  return (
    <Row gutter={[16, 16]}>
      <Col span={24}>
        <Descriptions bordered column={2}>
          <Descriptions.Item label="总签到天数">{stats.total_checkin_days}</Descriptions.Item>
          <Descriptions.Item label="连续签到天数">{stats.consecutive_checkin_days}</Descriptions.Item>
          <Descriptions.Item label="本月签到天数">{stats.current_month_checkin_days}</Descriptions.Item>
          <Descriptions.Item label="今年签到天数">{stats.current_year_checkin_days}</Descriptions.Item>
        </Descriptions>
      </Col>
      <Col span={24}>
        <ResponsiveContainer width="100%" height={300}>
          <BarChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="name" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Bar dataKey="训练次数" fill="#8884d8" />
          </BarChart>
        </ResponsiveContainer>
      </Col>
    </Row>
  );
};

export default UserStats;