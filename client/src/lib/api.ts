import axios from 'axios'

const client = axios.create({ baseURL: `${import.meta.env.VITE_API_URL ?? 'http://127.0.0.1:8080'}/api/v1` })

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
  filament_id: string
  customer: {
    name: string
    email: string
    address: string
  }
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
  const { data } = await client.post<OrderResponse>('/orders', req)
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
