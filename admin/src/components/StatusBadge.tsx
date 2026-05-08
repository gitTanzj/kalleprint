const bgByStatus: Record<string, string> = {
  pending: 'bg-gray-100',
  printing: 'bg-yellow-100',
  'waiting-for-shipment': 'bg-blue-100',
  shipped: 'bg-purple-100',
  done: 'bg-green-100',
}

export function StatusBadge({ status }: { status: string }) {
  const bg = bgByStatus[status] ?? 'bg-gray-100'
  return (
    <span
      className={`inline-block border-2 border-black px-2 py-0.5 text-xs font-semibold ${bg}`}
    >
      {status}
    </span>
  )
}
