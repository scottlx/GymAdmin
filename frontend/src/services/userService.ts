import apiClient from './api'

export interface User {
  id: number
  user_no: string
  name: string
  phone: string
  gender?: number
  email?: string
  status: number
  created_at: string
}

export interface UserListResponse {
  code: number
  message: string
  data: {
    list: User[]
    total: number
    page: number
    page_size: number
  }
}

export const userService = {
  list: async (params?: any): Promise<UserListResponse> => {
    return apiClient.get('/users', { params })
  },

  get: async (id: number): Promise<User> => {
    return apiClient.get(`/users/${id}`)
  },

  create: async (data: Partial<User>): Promise<User> => {
    return apiClient.post('/users', data)
  },

  update: async (id: number, data: Partial<User>): Promise<any> => {
    return apiClient.put(`/users/${id}`, data)
  },

  delete: async (id: number): Promise<any> => {
    return apiClient.delete(`/users/${id}`)
  },

  getStats: async (id: number): Promise<any> => {
    return apiClient.get(`/users/${id}/stats`)
  },
}