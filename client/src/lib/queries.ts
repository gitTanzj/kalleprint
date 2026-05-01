import { useQuery, useMutation } from '@tanstack/react-query'
import { getFilaments, getHealth, postOrder, postQuote } from '#/lib/api'
import type { OrderRequest } from '#/lib/api'

export function useHealth() {
  return useQuery({
    queryKey: ['health'],
    queryFn: getHealth,
    refetchInterval: 30_000,
  })
}

export function useFilaments() {
  return useQuery({
    queryKey: ['filaments'],
    queryFn: getFilaments,
  })
}

export function useGetQuote() {
  return useMutation({
    mutationFn: ({ file, filamentId }: { file: File; filamentId: string }) =>
      postQuote(file, filamentId),
  })
}

export function usePlaceOrder() {
  return useMutation({
    mutationFn: (req: OrderRequest) => postOrder(req),
  })
}
