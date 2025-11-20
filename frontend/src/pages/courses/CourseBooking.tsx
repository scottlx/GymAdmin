import React, { useEffect, useState } from 'react';
import { Calendar, Badge, List, message } from 'antd';
import { courseService, Course } from '../../services/courseService';
import dayjs, { Dayjs } from 'dayjs';

const CourseBooking: React.FC = () => {
  const [courses, setCourses] = useState<Course[]>([]);

  useEffect(() => {
    const fetchCourses = async () => {
      try {
        const response = await courseService.list({ page: 1, page_size: 1000 }); // Fetch all courses
        setCourses(response.data.list || []);
      } catch (error) {
        message.error('获取课程列表失败');
      }
    };
    fetchCourses();
  }, []);

  const getListData = (value: Dayjs) => {
    let listData = courses.filter(course => {
      return dayjs(course.start_time).isSame(value, 'day');
    });
    return listData || [];
  };

  const dateCellRender = (value: Dayjs) => {
    const listData = getListData(value);
    return (
      <List
        dataSource={listData}
        renderItem={(item) => (
          <List.Item>
            <Badge status={item.status === 1 ? 'success' : 'error'} text={item.course_name} />
          </List.Item>
        )}
      />
    );
  };

  return (
    <div>
      <h2>Course Booking</h2>
      <Calendar dateCellRender={dateCellRender} />
    </div>
  );
};

export default CourseBooking;