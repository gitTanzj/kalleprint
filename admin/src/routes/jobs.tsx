import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useJobs, useUpdateJobStatus } from '../lib/queries'

export const Route = createFileRoute('/jobs')({ component: JobsPage })

const STATUSES = [
  { label: 'All', value: '' },
  { label: 'Pending', value: 'pending' },
  { label: 'Printing', value: 'printing' },
  { label: 'Done', value: 'done' },
]

const JOB_STATUSES = ['pending', 'printing', 'done']

function JobsPage() {
  const navigate = useNavigate()
  const { data: jobs, isLoading } = useJobs()
  const { mutate: updateStatus } = useUpdateJobStatus()

  const [statusFilter, setStatusFilter] = useState('')

  const filtered = jobs?.filter((j) => !statusFilter || j.status === statusFilter)

  return (
    <div>
      <h1 className="mb-6 text-2xl font-bold">Jobs</h1>

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
              <th className="px-4 py-2 font-semibold">ID</th>
              <th className="px-4 py-2 font-semibold">Order</th>
              <th className="px-4 py-2 font-semibold">Filament</th>
              <th className="px-4 py-2 font-semibold">Status</th>
              <th className="px-4 py-2 font-semibold">Created</th>
            </tr>
          </thead>
          <tbody>
            {filtered?.map((job) => (
              <tr key={job.id} className="border-b border-black last:border-b-0">
                <td className="px-4 py-2 font-mono text-sm text-gray-500">
                  {job.id.slice(0, 8)}
                </td>
                <td className="px-4 py-2">
                  <button
                    onClick={() =>
                      navigate({ to: '/orders/$orderId', params: { orderId: job.order_id } })
                    }
                    className="font-mono text-sm underline hover:no-underline"
                  >
                    {job.order_id.slice(0, 8)}
                  </button>
                </td>
                <td className="px-4 py-2">{job.filament_name || job.filament_id.slice(0, 8)}</td>
                <td className="px-4 py-2">
                  <select
                    value={job.status}
                    onChange={(e) => updateStatus({ id: job.id, status: e.target.value })}
                    onClick={(e) => e.stopPropagation()}
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
            {filtered?.length === 0 && (
              <tr>
                <td colSpan={5} className="px-4 py-6 text-center text-gray-500">
                  No jobs
                </td>
              </tr>
            )}
          </tbody>
        </table>
      )}
    </div>
  )
}
