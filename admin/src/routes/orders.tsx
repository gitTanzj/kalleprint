import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { StatusBadge } from '../components/StatusBadge'
import { useOrders } from '../lib/queries'

export const Route = createFileRoute('/orders')({ component: OrdersPage })

const STATUSES = [
  { label: 'All', value: '' },
  { label: 'Pending', value: 'pending' },
  { label: 'Printing', value: 'printing' },
  { label: 'Waiting for Shipment', value: 'waiting-for-shipment' },
  { label: 'Shipped', value: 'shipped' },
  { label: 'Done', value: 'done' },
]

function OrdersPage() {
  const navigate = useNavigate()
  const [statusFilter, setStatusFilter] = useState('')
  const { data: orders, isLoading } = useOrders(statusFilter || undefined)

  return (
    <div>
      <h1 className="mb-6 text-2xl font-bold">Orders</h1>

      <div className="mb-6 flex gap-2">
        {STATUSES.map((s) => (
          <button
            key={s.value}
            onClick={() => setStatusFilter(s.value)}
            className={`border-2 border-black px-4 py-1.5 text-sm font-semibold ${
              statusFilter === s.value ? 'bg-black text-white' : 'bg-white hover:bg-gray-100'
            }`}
          >
            {s.label}
          </button>
        ))}
      </div>

      {isLoading ? (
        <p>Loading…</p>
      ) : (
        <table className="w-full border-2 border-black text-left">
          <thead className="border-b-2 border-black bg-gray-100">
            <tr>
              <th className="px-4 py-2 font-semibold">Name</th>
              <th className="px-4 py-2 font-semibold">Email</th>
              <th className="px-4 py-2 font-semibold">Total</th>
              <th className="px-4 py-2 font-semibold">Status</th>
              <th className="px-4 py-2 font-semibold">Date</th>
            </tr>
          </thead>
          <tbody>
            {orders?.map((order) => (
              <tr
                key={order.id}
                onClick={() =>
                  navigate({ to: '/orders/$orderId', params: { orderId: order.id } })
                }
                className="cursor-pointer border-b border-black last:border-b-0 hover:bg-gray-50"
              >
                <td className="px-4 py-2 font-semibold">{order.name}</td>
                <td className="px-4 py-2 text-gray-500">{order.email}</td>
                <td className="px-4 py-2">€{order.total.toFixed(2)}</td>
                <td className="px-4 py-2">
                  <StatusBadge status={order.status} />
                </td>
                <td className="px-4 py-2 text-gray-500">
                  {new Date(order.created_at).toLocaleDateString()}
                </td>
              </tr>
            ))}
            {orders?.length === 0 && (
              <tr>
                <td colSpan={5} className="px-4 py-6 text-center text-gray-500">
                  No orders
                </td>
              </tr>
            )}
          </tbody>
        </table>
      )}
    </div>
  )
}
