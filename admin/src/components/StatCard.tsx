export function StatCard({
  label,
  value,
}: {
  label: string
  value: string | number
}) {
  return (
    <div className="border-2 border-black p-6">
      <div className="text-sm font-semibold uppercase tracking-wide text-gray-500">
        {label}
      </div>
      <div className="mt-2 text-3xl font-bold">{value}</div>
    </div>
  )
}
