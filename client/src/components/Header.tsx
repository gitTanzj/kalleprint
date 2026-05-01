import { useHealth } from '#/lib/queries'

export default function Header() {
  const { data: healthy } = useHealth()

  return (
    <header className="sticky top-0 z-50 border-b-2 border-black bg-white px-4">
      <div className="page-wrap flex items-center justify-between py-4">
        <span className="text-lg font-bold tracking-tight">Kalleprint</span>
        <div className="flex items-center gap-2 text-sm">
          <span
            className={`h-2.5 w-2.5 rounded-full ${healthy === true ? 'bg-green-500' : healthy === false ? 'bg-red-500' : 'bg-gray-300'}`}
          />
          <span className="text-gray-500">
            {healthy === true ? 'Online' : healthy === false ? 'Offline' : 'Connecting…'}
          </span>
        </div>
      </div>
    </header>
  )
}
