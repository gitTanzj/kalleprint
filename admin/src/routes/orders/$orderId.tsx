import { createFileRoute, Link } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { StatusBadge } from '../../components/StatusBadge'
import { useOrder, useUpdateOrderStatus } from '../../lib/queries'

export const Route = createFileRoute('/orders/$orderId')({
  component: OrderDetailPage,
})

const VALID_STATUSES = [
  'pending',
  'printing',
  'waiting-for-shipment',
  'shipped',
  'done',
]

function OrderDetailPage() {
  const { orderId } = Route.useParams()
  const { data: order, isLoading } = useOrder(orderId)
  const updateStatus = useUpdateOrderStatus()
  const [status, setStatus] = useState('')
  const [saved, setSaved] = useState(false)

  useEffect(() => {
    if (order) setStatus(order.status)
  }, [order])

  async function handleSave() {
    if (!order) return
    updateStatus.mutate(
      { id: order.id, status },
      {
        onSuccess: () => {
          setSaved(true)
          setTimeout(() => setSaved(false), 2000)
        },
      },
    )
  }

  if (isLoading) return <p>Loading…</p>
  if (!order) return <p>Order not found</p>

  return (
    <div className="max-w-2xl">
      <Link to="/orders" className="mb-6 inline-block font-semibold underline">
        ← Back to orders
      </Link>
      <h1 className="mb-6 text-2xl font-bold">Order Detail</h1>

      <div className="mb-6 border-2 border-black p-6">
        <dl className="grid grid-cols-2 gap-x-8 gap-y-3">
          <dt className="font-semibold">Name</dt>
          <dd>{order.name}</dd>
          <dt className="font-semibold">Email</dt>
          <dd>{order.email}</dd>
          <dt className="font-semibold">Shipping Address</dt>
          <dd className="whitespace-pre-line">{order.shipping_address}</dd>
          {order.notes && (
            <>
              <dt className="font-semibold">Notes</dt>
              <dd>{order.notes}</dd>
            </>
          )}
          <dt className="font-semibold">Total</dt>
          <dd>€{order.total.toFixed(2)}</dd>
          <dt className="font-semibold">Status</dt>
          <dd>
            <StatusBadge status={order.status} />
          </dd>
          <dt className="font-semibold">Date</dt>
          <dd>{new Date(order.created_at).toLocaleString()}</dd>
        </dl>
      </div>

      <div className="border-2 border-black p-6">
        <h2 className="mb-4 font-bold">Update Status</h2>
        <div className="flex items-center gap-4">
          <select
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            className="border-2 border-black bg-white px-3 py-2 focus:outline-none"
          >
            {VALID_STATUSES.map((s) => (
              <option key={s} value={s}>
                {s}
              </option>
            ))}
          </select>
          <button
            onClick={handleSave}
            disabled={updateStatus.isPending || status === order.status}
            className="hover-black border-2 border-black bg-white px-6 py-2 font-semibold disabled:opacity-50"
          >
            {updateStatus.isPending ? 'Saving…' : saved ? 'Saved!' : 'Save'}
          </button>
        </div>
      </div>
    </div>
  )
}
