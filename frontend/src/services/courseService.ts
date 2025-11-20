import apiClient from './api'

export interface Course {
  id: number
  coach_id: number
  course_name: string
  course_type: number
  start_time: string
  end_time: string
  max_capacity: number
  current_count: number
  price: number
  status: number
  description?: string
  remark?: string
  created_at: string
}

export interface CourseListResponse {
  code: number
  message: string
  data: {
    list: Course[]
    total: number
    page: number
    page_size: number
  }
}

export const courseService = {
  list: async (params?: any): Promise<CourseListResponse> => {
    return apiClient.get('/courses', { params })
  },

  get: async (id: number): Promise<Course> => {
    return apiClient.get(`/courses/${id}`)
  },

  create: async (data: Partial<Course>): Promise<Course> => {
    return apiClient.post('/courses', data)
  },

  update: async (id: number, data: Partial<Course>): Promise<any> => {
    return apiClient.put(`/courses/${id}`, data)
  },

  delete: async (id: number): Promise<any> => {
    return apiClient.delete(`/courses/${id}`)
  },
}
