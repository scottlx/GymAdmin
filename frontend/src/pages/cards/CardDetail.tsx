import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { cardService, MembershipCard } from '../../services/cardService';
import { Button, Card, Descriptions, message, Spin, Tag } from 'antd';
import dayjs from 'dayjs';

const CardDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [card, setCard] = useState<MembershipCard | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchCard = async () => {
      if (!id) return;
      setLoading(true);
      try {
        const response = await cardService.get(Number(id));
        setCard(response);
      } catch (error) {
        message.error('获取会员卡详情失败');
      } finally {
        setLoading(false);
      }
    };
    fetchCard();
  }, [id]);

  const getStatusTag = (status: number) => {
    const statusMap: any = {
      1: <Tag color="green">正常</Tag>,
      2: <Tag color="red">已过期</Tag>,
      3: <Tag color="orange">已冻结</Tag>,
      4: <Tag color="purple">已转出</Tag>,
      5: <Tag color="gray">已退卡</Tag>,
    };
    return statusMap[status] || '-';
  };

  if (loading) {
    return <Spin />;
  }

  if (!card) {
    return <div>会员卡不存在</div>;
  }

  return (
    <div>
      <Button onClick={() => navigate(-1)} style={{ marginBottom: 16 }}>
        返回
      </Button>
      <Card title="会员卡详情">
        <Descriptions bordered>
          <Descriptions.Item label="卡号">{card.card_no}</Descriptions.Item>
          <Descriptions.Item label="会员ID">{card.user_id}</Descriptions.Item>
          <Descriptions.Item label="状态">{getStatusTag(card.status)}</Descriptions.Item>
          <Descriptions.Item label="开始日期">{dayjs(card.start_date).format('YYYY-MM-DD')}</Descriptions.Item>
          <Descriptions.Item label="到期日期">{dayjs(card.end_date).format('YYYY-MM-DD')}</Descriptions.Item>
          <Descriptions.Item label="购买价格">¥{card.purchase_price}</Descriptions.Item>
          <Descriptions.Item label="剩余次数">{card.remaining_times ?? 'N/A'}</Descriptions.Item>
          <Descriptions.Item label="总次数">{card.total_times ?? 'N/A'}</Descriptions.Item>
          <Descriptions.Item label="冻结次数">{card.freeze_times}</Descriptions.Item>
          <Descriptions.Item label="冻结天数">{card.freeze_days}</Descriptions.Item>
          <Descriptions.Item label="备注">{card.remark}</Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  );
};

export default CardDetail;