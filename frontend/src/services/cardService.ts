import apiClient from './api'

export interface MembershipCard {
  id: number
  card_no: string
  user_id: number
  card_type_id: number
  status: number
  start_date: string
  end_date: string
  remaining_times?: number
  total_times?: number
  freeze_times: number
  freeze_days: number
  is_frozen: number
  purchase_price: number
  remark?: string
  created_at: string
}

export interface CardListResponse {
  code: number
  message: string
  data: {
    list: MembershipCard[]
    total: number
    page: number
    page_size: number
  }
}

export const cardService = {
  list: async (params?: any): Promise<CardListResponse> => {
    return apiClient.get('/cards', { params })
  },

  get: async (id: number): Promise<MembershipCard> => {
    return apiClient.get(`/cards/${id}`)
  },

  create: async (data: Partial<MembershipCard>): Promise<MembershipCard> => {
    return apiClient.post('/cards', data)
  },

  update: async (id: number, data: Partial<MembershipCard>): Promise<any> => {
    return apiClient.put(`/cards/${id}`, data)
  },

  delete: async (id: number): Promise<any> => {
    return apiClient.delete(`/cards/${id}`)
  },
}
