import axios from 'axios'

const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

export const api = axios.create({
  baseURL: `${BASE_URL}/api/v1/admin`,
})

api.interceptors.request.use((config) => {
  const token = typeof window !== 'undefined' ? localStorage.getItem('admin_token') : null
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export function getToken() {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('admin_token')
}

export function isTokenValid() {
  if (typeof window === 'undefined') return false
  const token = localStorage.getItem('admin_token')
  const expiresAt = localStorage.getItem('admin_token_expires_at')
  if (!token || !expiresAt) return false
  return new Date(expiresAt) > new Date()
}

export function saveToken(token: string, expiresAt: string) {
  if (typeof window === 'undefined') return
  localStorage.setItem('admin_token', token)
  localStorage.setItem('admin_token_expires_at', expiresAt)
}

export function clearToken() {
  if (typeof window === 'undefined') return
  localStorage.removeItem('admin_token')
  localStorage.removeItem('admin_token_expires_at')
}

export interface Order {
  id: string
  name: string
  email: string
  shipping_address: string
  notes: string
  print_id: string
  total: number
  status: string
  created_at: string
}

export interface Filament {
  id: string
  filament_name: string
  current_amount: number
  cost_per_gram: number
}

export interface DashboardStats {
  total_orders: number
  pending: number
  in_progress: number
  total_revenue: number
}

export async function login(email: string, password: string) {
  const res = await api.post<{ token: string; expires_at: string }>('/login', {
    email,
    password,
  })
  return res.data
}

export async function fetchDashboard() {
  const res = await api.get<DashboardStats>('/dashboard')
  return res.data
}

export async function fetchOrders(status?: string) {
  const res = await api.get<Order[]>('/orders', {
    params: status ? { status } : undefined,
  })
  return res.data
}

export async function fetchOrder(id: string) {
  const res = await api.get<Order[]>('/orders')
  const order = res.data.find((o) => o.id === id)
  if (!order) throw new Error('Order not found')
  return order
}

export async function updateOrderStatus(id: string, status: string) {
  const res = await api.patch<{ status: string }>(`/orders/${id}/status`, {
    status,
  })
  return res.data
}

export async function fetchFilaments() {
  const res = await api.get<Filament[]>('/filaments')
  return res.data
}

export async function createFilament(data: {
  name: string
  amount_grams: number
  total_price: number
}) {
  const res = await api.post<Filament>('/filaments', data)
  return res.data
}

export async function updateFilament(
  id: string,
  data: { name: string; amount_grams: number; total_price: number },
) {
  const res = await api.put<Filament>(`/filaments/${id}`, data)
  return res.data
}

export async function deleteFilament(id: string) {
  await api.delete(`/filaments/${id}`)
}
