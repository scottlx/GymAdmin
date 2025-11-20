import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { userService, User } from '../../services/userService';
import { Button, Card, Descriptions, message, Spin } from 'antd';
import UserStats from './components/UserStats';

const UserDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUser = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const response = await userService.get(Number(id));
        setUser(response);
      } catch (error) {
        message.error('获取用户详情失败');
      } finally {
        setLoading(false);
      }
    };
    fetchUser();
  }, [id]);

  if (loading) {
    return <Spin />;
  }

  if (!user) {
    return <div>用户不存在</div>;
  }

  return (
    <div>
      <Button onClick={() => navigate(-1)} style={{ marginBottom: 16 }}>
        返回
      </Button>
      <Card title="用户详情">
        <Descriptions bordered>
          <Descriptions.Item label="姓名">{user.name}</Descriptions.Item>
          <Descriptions.Item label="手机号">{user.phone}</Descriptions.Item>
          <Descriptions.Item label="用户编号">{user.user_no}</Descriptions.Item>
          <Descriptions.Item label="性别">{user.gender === 1 ? '男' : '女'}</Descriptions.Item>
          <Descriptions.Item label="邮箱">{user.email}</Descriptions.Item>
          <Descriptions.Item label="状态">{user.status === 1 ? '正常' : '冻结'}</Descriptions.Item>
        </Descriptions>
      </Card>
      <Card title="训练统计" style={{ marginTop: 16 }}>
        <UserStats />
      </Card>
    </div>
  );
};

export default UserDetail;
