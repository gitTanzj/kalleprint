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
  color_hex: string
  ini_file_path: string
  filament_type: string
  temperature: string
  bed_temperature: string
  first_layer_temperature: string
  first_layer_bed_temperature: string
  filament_diameter: string
  extrusion_multiplier: string
  filament_density: string
  filament_cost: string
}

export interface CreateFilamentInput {
  name: string
  amount_grams: number
  total_price: number
  color_hex: string
  filament_type: string
  temperature?: string
  bed_temperature?: string
  first_layer_temperature?: string
  first_layer_bed_temperature?: string
  filament_diameter?: string
  extrusion_multiplier?: string
  filament_density?: string
}

export interface DashboardStats {
  total_orders: number
  pending: number
  in_progress: number
  total_revenue: number
  realized_revenue: number
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
  const res = await api.get<Order>(`/orders/${id}`)
  return res.data
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

export async function createFilament(data: CreateFilamentInput) {
  const fd = new FormData()
  fd.append('filament_name', data.name)
  fd.append('filament_stock', String(data.amount_grams))
  fd.append('filament_cost_per_gram', String(data.total_price / data.amount_grams))
  fd.append('filament_color_hex', data.color_hex)
  fd.append('filament_type', data.filament_type)
  fd.append('filament_cost', String(data.total_price))

  const optional: Array<[keyof CreateFilamentInput, string]> = [
    ['temperature', 'temperature'],
    ['bed_temperature', 'bed_temperature'],
    ['first_layer_temperature', 'first_layer_temperature'],
    ['first_layer_bed_temperature', 'first_layer_bed_temperature'],
    ['filament_diameter', 'filament_diameter'],
    ['extrusion_multiplier', 'extrusion_multiplier'],
    ['filament_density', 'filament_density'],
  ]
  for (const [key, field] of optional) {
    const v = data[key]
    if (typeof v === 'string' && v !== '') fd.append(field, v)
  }

  const res = await api.post<Filament>('/filaments', fd)
  return res.data
}

export interface UpdateFilamentInput {
  name?: string
  amount_grams?: number
  total_price?: number
  color_hex?: string
  filament_type?: string
  temperature?: string
  bed_temperature?: string
  first_layer_temperature?: string
  first_layer_bed_temperature?: string
  filament_diameter?: string
  extrusion_multiplier?: string
  filament_density?: string
}

export async function updateFilament(id: string, data: UpdateFilamentInput) {
  const fd = new FormData()
  if (data.name) fd.append('filament_name', data.name)
  if (data.amount_grams != null) fd.append('filament_stock', String(data.amount_grams))
  if (data.amount_grams != null && data.total_price != null)
    fd.append('filament_cost_per_gram', String(data.total_price / data.amount_grams))
  if (data.total_price != null) fd.append('filament_cost', String(data.total_price))
  if (data.color_hex) fd.append('filament_color_hex', data.color_hex)

  const iniFields: Array<[keyof UpdateFilamentInput, string]> = [
    ['filament_type', 'filament_type'],
    ['temperature', 'temperature'],
    ['bed_temperature', 'bed_temperature'],
    ['first_layer_temperature', 'first_layer_temperature'],
    ['first_layer_bed_temperature', 'first_layer_bed_temperature'],
    ['filament_diameter', 'filament_diameter'],
    ['extrusion_multiplier', 'extrusion_multiplier'],
    ['filament_density', 'filament_density'],
  ]
  for (const [key, field] of iniFields) {
    const v = data[key]
    if (typeof v === 'string' && v !== '') fd.append(field, v)
  }

  const res = await api.put<Filament>(`/filaments/${id}`, fd)
  return res.data
}

export async function deleteFilament(id: string) {
  await api.delete(`/filaments/${id}`)
}

export interface Job {
  id: string
  order_id: string
  filament_id: string
  filament_name: string
  gcode_path: string
  status: string
  created_at: string
}

export async function fetchJobs(status?: string) {
  const res = await api.get<Job[]>('/jobs', {
    params: status ? { status } : undefined,
  })
  return res.data
}

export async function updateJobStatus(id: string, status: string) {
  const res = await api.patch<{ status: string }>(`/jobs/${id}/status`, { status })
  return res.data
}
