import { createFileRoute, Link } from '@tanstack/react-router'
import { StatCard } from '../components/StatCard'
import { StatusBadge } from '../components/StatusBadge'
import { useDashboard, useOrders } from '../lib/queries'

export const Route = createFileRoute('/')({ component: DashboardPage })

function DashboardPage() {
  const { data: stats, isLoading: statsLoading } = useDashboard()
  const { data: orders } = useOrders()

  const recent = orders?.slice(0, 10) ?? []

  return (
    <div>
      <h1 className="mb-6 text-2xl font-bold">Dashboard</h1>

      {statsLoading ? (
        <p>Loading…</p>
      ) : stats ? (
        <div className="mb-8 grid grid-cols-5 gap-4">
          <StatCard label="Total Orders" value={stats.total_orders} />
          <StatCard label="Pending" value={stats.pending} />
          <StatCard label="In Progress" value={stats.in_progress} />
          <StatCard
            label="Total Revenue"
            value={`€${stats.total_revenue.toFixed(2)}`}
          />
          <StatCard
            label="Realized Revenue"
            value={`€${stats.realized_revenue.toFixed(2)}`}
          />
        </div>
      ) : null}

      <h2 className="mb-3 text-lg font-bold">Recent Orders</h2>
      <table className="w-full border-2 border-black text-left">
        <thead className="border-b-2 border-black bg-gray-100">
          <tr>
            <th className="px-4 py-2 font-semibold">Name</th>
            <th className="px-4 py-2 font-semibold">Status</th>
            <th className="px-4 py-2 font-semibold">Total</th>
            <th className="px-4 py-2 font-semibold">Date</th>
          </tr>
        </thead>
        <tbody>
          {recent.map((order) => (
            <tr
              key={order.id}
              className="border-b border-black last:border-b-0 hover:bg-gray-50"
            >
              <td className="px-4 py-2">
                <Link
                  to="/orders/$orderId"
                  params={{ orderId: order.id }}
                  className="font-semibold underline"
                >
                  {order.name}
                </Link>
              </td>
              <td className="px-4 py-2">
                <StatusBadge status={order.status} />
              </td>
              <td className="px-4 py-2">€{order.total.toFixed(2)}</td>
              <td className="px-4 py-2 text-gray-500">
                {new Date(order.created_at).toLocaleDateString()}
              </td>
            </tr>
          ))}
          {recent.length === 0 && (
            <tr>
              <td colSpan={4} className="px-4 py-6 text-center text-gray-500">
                No orders yet
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  )
}
