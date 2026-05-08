import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  createFilament,
  deleteFilament,
  fetchDashboard,
  fetchFilaments,
  fetchOrder,
  fetchOrders,
  login,
  saveToken,
  updateFilament,
  updateOrderStatus,
} from './api'

export function useDashboard() {
  return useQuery({ queryKey: ['dashboard'], queryFn: fetchDashboard })
}

export function useOrders(status?: string) {
  return useQuery({
    queryKey: ['orders', status],
    queryFn: () => fetchOrders(status),
  })
}

export function useOrder(id: string) {
  return useQuery({ queryKey: ['orders', id], queryFn: () => fetchOrder(id) })
}

export function useUpdateOrderStatus() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      updateOrderStatus(id, status),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ['orders'] })
    },
  })
}

export function useFilaments() {
  return useQuery({ queryKey: ['filaments'], queryFn: fetchFilaments })
}

export function useCreateFilament() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: createFilament,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['filaments'] }),
  })
}

export function useUpdateFilament() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({
      id,
      ...data
    }: { id: string; name: string; amount_grams: number; total_price: number }) =>
      updateFilament(id, data),
    onSuccess: () => qc.invalidateQueries({ queryKey: ['filaments'] }),
  })
}

export function useDeleteFilament() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: deleteFilament,
    onSuccess: () => qc.invalidateQueries({ queryKey: ['filaments'] }),
  })
}

export function useLogin() {
  return useMutation({
    mutationFn: ({ email, password }: { email: string; password: string }) =>
      login(email, password),
    onSuccess: (data) => {
      saveToken(data.token, data.expires_at)
    },
  })
}
