import { useState, useRef } from 'react'
import { Loader2 } from 'lucide-react'
import { useFilaments, usePlaceOrder } from '#/lib/queries'
import type { Quote, OrderResponse } from '#/lib/api'

interface Props {
  quote: Quote
  file: File
  filename: string
  filamentId: string
  isRequoting?: boolean
  requoteError?: string
  onFilamentChange: (filamentId: string) => void
  onRequote: () => void
  onConfirmed: (result: OrderResponse, customerName: string) => void
}

export default function QuoteSection({ quote, file, filename, filamentId, isRequoting, requoteError, onFilamentChange, onRequote, onConfirmed }: Props) {
  const { data: filaments = [], isLoading: filamentsLoading } = useFilaments()
  const placeOrder = usePlaceOrder()

  const [name, setName] = useState('')
  const [email, setEmail] = useState('')
  const [address, setAddress] = useState('')
  const [notes, setNotes] = useState('')
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  function handleFilamentChange(id: string) {
    if (debounceRef.current) clearTimeout(debounceRef.current)
    debounceRef.current = setTimeout(() => onFilamentChange(id), 500)
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const result = await placeOrder.mutateAsync({
      file,
      filament_id: filamentId,
      name,
      email,
      shipping_address: address,
      notes,
    })
    onConfirmed(result, name)
  }

  const price = new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'EUR',
  }).format(quote.price_eur)

  return (
    <section className="bg-gray-700 px-4 py-20 text-white">
      <div className="page-wrap flex flex-col gap-12 md:flex-row">
        {/* Price card */}
        <div className="md:w-72 flex-shrink-0">
          <div className="border-2 border-white p-8">
            <p className="mb-1 text-sm uppercase tracking-widest text-gray-400">Your quote</p>
            {isRequoting ? (
              <div className="mb-4 flex items-center gap-3">
                <Loader2 className="animate-spin text-gray-400" size={32} />
                <span className="text-lg text-gray-400">Recalculating…</span>
              </div>
            ) : (
              <p className="mb-4 text-5xl font-bold">{price}</p>
            )}
            <p className="mb-6 break-all text-sm text-gray-300">{filename}</p>
            {requoteError && (
              <p className="mb-4 text-sm text-red-400">{requoteError}</p>
            )}
            <button
              onClick={onRequote}
              className="text-sm underline underline-offset-2 hover:text-gray-300"
            >
              Upload a different file
            </button>
          </div>
        </div>

        {/* Order form */}
        <form onSubmit={handleSubmit} className="flex-1 flex flex-col gap-5">
          <h2 className="text-3xl font-bold tracking-tight">Place Your Order</h2>

          <div className="grid gap-5 sm:grid-cols-2">
            <label className="flex flex-col gap-1.5 text-sm font-semibold">
              Name
              <input
                required
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="border-2 border-white bg-transparent px-3 py-2 font-normal text-white placeholder:text-gray-400 focus:outline-none focus:border-gray-300"
                placeholder="Jane Smith"
              />
            </label>
            <label className="flex flex-col gap-1.5 text-sm font-semibold">
              Email
              <input
                required
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="border-2 border-white bg-transparent px-3 py-2 font-normal text-white placeholder:text-gray-400 focus:outline-none focus:border-gray-300"
                placeholder="jane@example.com"
              />
            </label>
          </div>

          <label className="flex flex-col gap-1.5 text-sm font-semibold">
            Shipping Address
            <textarea
              required
              rows={3}
              value={address}
              onChange={(e) => setAddress(e.target.value)}
              className="border-2 border-white bg-transparent px-3 py-2 font-normal text-white placeholder:text-gray-400 focus:outline-none focus:border-gray-300 resize-none"
              placeholder="123 Main St, City, Country"
            />
          </label>

          <label className="flex flex-col gap-1.5 text-sm font-semibold">
            Filament
            <select
              required
              value={filamentId}
              onChange={(e) => handleFilamentChange(e.target.value)}
              disabled={filamentsLoading || isRequoting}
              className="border-2 border-white bg-gray-700 px-3 py-2 font-normal text-white focus:outline-none focus:border-gray-300 disabled:opacity-50"
            >
              <option value="" disabled>
                {filamentsLoading ? 'Loading…' : 'Select a filament'}
              </option>
              {filaments.map((f) => (
                <option key={f.id} value={f.id}>
                  {f.filament_name}
                </option>
              ))}
            </select>
          </label>

          <label className="flex flex-col gap-1.5 text-sm font-semibold">
            Notes
            <textarea
              rows={2}
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              className="border-2 border-white bg-transparent px-3 py-2 font-normal text-white placeholder:text-gray-400 focus:outline-none focus:border-gray-300 resize-none"
              placeholder="Any special requests?"
            />
          </label>

          <div>
            <button
              type="submit"
              disabled={placeOrder.isPending || !name || !email || !address || !filamentId}
              className="hover-black border-2 border-white bg-white px-8 py-3 font-semibold text-black disabled:cursor-not-allowed hover:cursor-pointer disabled:opacity-40"
            >
              {placeOrder.isPending ? 'Placing order…' : 'Place Order'}
            </button>
            {placeOrder.error && (
              <p className="mt-3 text-sm text-red-400">
                {placeOrder.error.message}
              </p>
            )}
          </div>
        </form>
      </div>
    </section>
  )
}
