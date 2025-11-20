import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { coachService, Coach } from '../../services/coachService';
import { Button, Card, Descriptions, message, Spin } from 'antd';

const CoachDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [coach, setCoach] = useState<Coach | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchCoach = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const response = await coachService.get(Number(id));
        setCoach(response);
      } catch (error) {
        message.error('获取教练详情失败');
      } finally {
        setLoading(false);
      }
    };
    fetchCoach();
  }, [id]);

  if (loading) {
    return <Spin />;
  }

  if (!coach) {
    return <div>教练不存在</div>;
  }

  return (
    <div>
      <Button onClick={() => navigate(-1)} style={{ marginBottom: 16 }}>
        返回
      </Button>
      <Card title="教练详情">
        <Descriptions bordered>
          <Descriptions.Item label="姓名">{coach.name}</Descriptions.Item>
          <Descriptions.Item label="手机号">{coach.phone}</Descriptions.Item>
          <Descriptions.Item label="教练编号">{coach.coach_no}</Descriptions.Item>
          <Descriptions.Item label="性别">{coach.gender === 1 ? '男' : '女'}</Descriptions.Item>
          <Descriptions.Item label="邮箱">{coach.email}</Descriptions.Item>
          <Descriptions.Item label="专长">{coach.specialties}</Descriptions.Item>
          <Descriptions.Item label="工作年限">{coach.experience}年</Descriptions.Item>
          <Descriptions.Item label="状态">{coach.status === 1 ? '正常' : '停用'}</Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  );
};

export default CoachDetail;
