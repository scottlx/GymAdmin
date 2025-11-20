import apiClient from './api'

export interface Coach {
  id: number
  coach_no: string
  name: string
  phone: string
  gender?: number
  email?: string
  specialties?: string
  experience?: number
  status: number
  created_at: string
}

export interface CoachListResponse {
  code: number
  message: string
  data: {
    list: Coach[]
    total: number
    page: number
    page_size: number
  }
}

export const coachService = {
  list: async (params?: any): Promise<CoachListResponse> => {
    return apiClient.get('/coaches', { params })
  },

  get: async (id: number): Promise<Coach> => {
    return apiClient.get(`/coaches/${id}`)
  },

  create: async (data: Partial<Coach>): Promise<Coach> => {
    return apiClient.post('/coaches', data)
  },

  update: async (id: number, data: Partial<Coach>): Promise<any> => {
    return apiClient.put(`/coaches/${id}`, data)
  },

  delete: async (id: number): Promise<any> => {
    return apiClient.delete(`/coaches/${id}`)
  },
}
