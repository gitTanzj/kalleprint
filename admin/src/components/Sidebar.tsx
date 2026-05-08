import { Link } from '@tanstack/react-router'
import { clearToken } from '../lib/api'

export function Sidebar() {
  function handleLogout() {
    clearToken()
    window.location.href = '/login'
  }

  return (
    <aside className="flex w-56 shrink-0 flex-col border-r-2 border-black bg-white">
      <div className="border-b-2 border-black px-6 py-4">
        <span className="text-lg font-bold">Kalleprint Admin</span>
      </div>
      <nav className="flex flex-1 flex-col gap-1 p-4">
        <Link
          to="/"
          className="border-2 border-transparent px-3 py-2 font-semibold hover:border-black"
          activeProps={{ className: 'border-2 border-black px-3 py-2 font-semibold bg-gray-100' }}
          activeOptions={{ exact: true }}
        >
          Dashboard
        </Link>
        <Link
          to="/orders"
          className="border-2 border-transparent px-3 py-2 font-semibold hover:border-black"
          activeProps={{ className: 'border-2 border-black px-3 py-2 font-semibold bg-gray-100' }}
        >
          Orders
        </Link>
        <Link
          to="/filaments"
          className="border-2 border-transparent px-3 py-2 font-semibold hover:border-black"
          activeProps={{ className: 'border-2 border-black px-3 py-2 font-semibold bg-gray-100' }}
        >
          Filaments
        </Link>
      </nav>
      <div className="border-t-2 border-black p-4">
        <button
          onClick={handleLogout}
          className="hover-black w-full border-2 border-black px-4 py-2 font-semibold bg-red-200 hover:bg-red-400 cursor-pointer"
        >
          Logout
        </button>
      </div>
    </aside>
  )
}
