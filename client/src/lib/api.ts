import axios from 'axios'

const client = axios.create({ baseURL: `${import.meta.env.VITE_API_URL ?? 'http://localhost:8080'}/api/v1` })

client.interceptors.response.use(
  (res) => res,
  (err) => {
    const message = err.response?.data?.error ?? err.message
    return Promise.reject(new Error(message))
  },
)

export interface Quote {
  price_eur: number
}

export interface Filament {
  id: string
  filament_name: string
  current_amount: number
  cost_per_gram: number
}

export interface OrderRequest {
  file: File
  filament_id: string
  name: string
  email: string
  shipping_address: string
  notes: string
}

export interface OrderResponse {
  id: string
}

export async function postQuote(file: File, filamentId: string): Promise<Quote> {
  const form = new FormData()
  form.append('printable', file)
  form.append('filament_id', filamentId)
  const { data } = await client.post<Quote>('/quote', form)
  return data
}

export async function getFilaments(): Promise<Filament[]> {
  const { data } = await client.get<Filament[]>('/filaments')
  return data
}

export async function postOrder(req: OrderRequest): Promise<OrderResponse> {
  const form = new FormData()
  form.append('printable', req.file)
  form.append('filament_id', req.filament_id)
  form.append('name', req.name)
  form.append('email', req.email)
  form.append('shipping_address', req.shipping_address)
  form.append('notes', req.notes)
  const { data } = await client.post<OrderResponse>('/orders', form)
  return data
}

export async function getHealth(): Promise<boolean> {
  try {
    await client.get('/health')
    return true
  } catch {
    return false
  }
}
