import { createFileRoute, Link } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { StatusBadge } from '../../components/StatusBadge'
import { useOrder, useUpdateOrderStatus, useJobs, useUpdateJobStatus } from '../../lib/queries'

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

const JOB_STATUSES = ['pending', 'printing', 'done']

function OrderDetailPage() {
  const { orderId } = Route.useParams()
  const { data: order, isLoading } = useOrder(orderId)
  const { data: allJobs } = useJobs()
  const updateStatus = useUpdateOrderStatus()
  const { mutate: updateJobStatus } = useUpdateJobStatus()
  const [status, setStatus] = useState('')
  const [saved, setSaved] = useState(false)

  const jobs = allJobs?.filter((j) => j.order_id === orderId) ?? []

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

      <div className="mb-6 border-2 border-black p-6">
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

      <div className="border-2 border-black p-6">
        <h2 className="mb-4 font-bold">Jobs</h2>
        {jobs.length === 0 ? (
          <p className="text-gray-500">No jobs for this order</p>
        ) : (
          <table className="w-full border-2 border-black text-left">
            <thead className="border-b-2 border-black bg-gray-100">
              <tr>
                <th className="px-4 py-2 font-semibold">ID</th>
                <th className="px-4 py-2 font-semibold">Filament</th>
                <th className="px-4 py-2 font-semibold">Status</th>
                <th className="px-4 py-2 font-semibold">Created</th>
              </tr>
            </thead>
            <tbody>
              {jobs.map((job) => (
                <tr key={job.id} className="border-b border-black last:border-b-0">
                  <td className="px-4 py-2 font-mono text-sm text-gray-500">
                    {job.id.slice(0, 8)}
                  </td>
                  <td className="px-4 py-2">{job.filament_name || job.filament_id.slice(0, 8)}</td>
                  <td className="px-4 py-2">
                    <select
                      value={job.status}
                      onChange={(e) => updateJobStatus({ id: job.id, status: e.target.value })}
                      className="border-2 border-black bg-white px-2 py-1 text-sm font-semibold"
                    >
                      {JOB_STATUSES.map((s) => (
                        <option key={s} value={s}>
                          {s}
                        </option>
                      ))}
                    </select>
                  </td>
                  <td className="px-4 py-2 text-gray-500">
                    {new Date(job.created_at).toLocaleDateString()}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
