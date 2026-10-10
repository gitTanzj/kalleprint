import axios from 'axios'

const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

export const api = axios.create({
  baseURL: `${BASE_URL}/api/v1/admin`,
})

api.interceptors.request.use((config) => {
  const token =
    typeof window !== 'undefined' ? localStorage.getItem('admin_token') : null
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

export interface SlicerField {
  key: string
  label: string
  /** Server-side default, shown as the input placeholder. */
  placeholder: string
  step?: string
  group: 'Temperature' | 'Flow' | 'Cooling' | 'Retraction'
}

// Per-filament PrusaSlicer settings editable in the admin panel. Printer and
// print settings come from the server's static Creality K1 profiles.
export const SLICER_FIELDS = [
  {
    key: 'temperature',
    label: 'Nozzle Temp (°C)',
    placeholder: '210',
    group: 'Temperature',
  },
  {
    key: 'first_layer_temperature',
    label: 'First Layer Nozzle (°C)',
    placeholder: '215',
    group: 'Temperature',
  },
  {
    key: 'bed_temperature',
    label: 'Bed Temp (°C)',
    placeholder: '60',
    group: 'Temperature',
  },
  {
    key: 'first_layer_bed_temperature',
    label: 'First Layer Bed (°C)',
    placeholder: '60',
    group: 'Temperature',
  },
  {
    key: 'filament_max_volumetric_speed',
    label: 'Max Flow (mm³/s)',
    placeholder: '18',
    step: '0.5',
    group: 'Flow',
  },
  {
    key: 'extrusion_multiplier',
    label: 'Extrusion Mult.',
    placeholder: '0.98',
    step: '0.01',
    group: 'Flow',
  },
  {
    key: 'pressure_advance',
    label: 'Pressure Advance',
    placeholder: '0.04',
    step: '0.001',
    group: 'Flow',
  },
  {
    key: 'filament_diameter',
    label: 'Diameter (mm)',
    placeholder: '1.75',
    step: '0.01',
    group: 'Flow',
  },
  {
    key: 'filament_density',
    label: 'Density (g/cm³)',
    placeholder: '1.24',
    step: '0.01',
    group: 'Flow',
  },
  {
    key: 'min_fan_speed',
    label: 'Min Fan (%)',
    placeholder: '100',
    group: 'Cooling',
  },
  {
    key: 'max_fan_speed',
    label: 'Max Fan (%)',
    placeholder: '100',
    group: 'Cooling',
  },
  {
    key: 'bridge_fan_speed',
    label: 'Bridge Fan (%)',
    placeholder: '100',
    group: 'Cooling',
  },
  {
    key: 'disable_fan_first_layers',
    label: 'Fan Off Layers',
    placeholder: '1',
    group: 'Cooling',
  },
  {
    key: 'full_fan_speed_layer',
    label: 'Full Fan From Layer',
    placeholder: '3',
    group: 'Cooling',
  },
  {
    key: 'slowdown_below_layer_time',
    label: 'Slow Down Below (s)',
    placeholder: '8',
    step: '0.5',
    group: 'Cooling',
  },
  {
    key: 'min_print_speed',
    label: 'Min Print Speed (mm/s)',
    placeholder: '15',
    group: 'Cooling',
  },
  {
    key: 'filament_retract_length',
    label: 'Retract Length (mm)',
    placeholder: '0.8',
    step: '0.05',
    group: 'Retraction',
  },
  {
    key: 'filament_retract_speed',
    label: 'Retract Speed (mm/s)',
    placeholder: '40',
    group: 'Retraction',
  },
] as const satisfies ReadonlyArray<SlicerField>

export type SlicerKey = (typeof SLICER_FIELDS)[number]['key']
export type SlicerValues = Partial<Record<SlicerKey, string>>

function appendSlicerValues(fd: FormData, values: SlicerValues) {
  for (const { key } of SLICER_FIELDS) {
    const v = values[key]
    if (v) fd.append(key, v)
  }
}

export interface Filament extends Record<SlicerKey, string> {
  id: string
  filament_name: string
  current_amount: number
  cost_per_gram: number
  color_hex: string
  ini_file_path: string
  filament_type: string
  filament_cost: string
}

export interface CreateFilamentInput {
  name: string
  amount_grams: number
  total_price: number
  color_hex: string
  filament_type: string
  slicer: SlicerValues
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
  fd.append(
    'filament_cost_per_gram',
    String(data.total_price / data.amount_grams),
  )
  fd.append('filament_color_hex', data.color_hex)
  fd.append('filament_type', data.filament_type)
  fd.append('filament_cost', String(data.total_price))
  appendSlicerValues(fd, data.slicer)

  const res = await api.post<Filament>('/filaments', fd)
  return res.data
}

export interface UpdateFilamentInput {
  name?: string
  amount_grams?: number
  total_price?: number
  color_hex?: string
  filament_type?: string
  slicer?: SlicerValues
}

export async function updateFilament(id: string, data: UpdateFilamentInput) {
  const fd = new FormData()
  if (data.name) fd.append('filament_name', data.name)
  if (data.amount_grams != null)
    fd.append('filament_stock', String(data.amount_grams))
  if (data.amount_grams != null && data.total_price != null)
    fd.append(
      'filament_cost_per_gram',
      String(data.total_price / data.amount_grams),
    )
  if (data.total_price != null)
    fd.append('filament_cost', String(data.total_price))
  if (data.color_hex) fd.append('filament_color_hex', data.color_hex)

  if (data.filament_type) fd.append('filament_type', data.filament_type)
  if (data.slicer) appendSlicerValues(fd, data.slicer)

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
  const res = await api.patch<{ status: string }>(`/jobs/${id}/status`, {
    status,
  })
  return res.data
}
