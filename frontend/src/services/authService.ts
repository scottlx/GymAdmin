import apiClient from './api'

export interface LoginRequest {
  phone: string
  password: string
}

export interface LoginResponse {
  code: number
  message: string
  data: {
    token: string
    user: any
  }
}

export const authService = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    return apiClient.post('/login', data)
  },

  register: async (data: any): Promise<any> => {
    return apiClient.post('/register', data)
  },
}
