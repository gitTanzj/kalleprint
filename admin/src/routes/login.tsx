import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { useLogin } from '../lib/queries'

export const Route = createFileRoute('/login')({ component: LoginPage })

function LoginPage() {
  const navigate = useNavigate()
  const loginMutation = useLogin()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    loginMutation.mutate(
      { email, password },
      { onSuccess: () => navigate({ to: '/' }) },
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-white">
      <div className="w-full max-w-sm border-2 border-black p-8">
        <h1 className="mb-6 text-2xl font-bold">Kalleprint Admin</h1>
        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
          <div className="flex flex-col gap-1">
            <label className="text-sm font-semibold">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              className="border-2 border-black bg-white px-3 py-2 focus:outline-none"
            />
          </div>
          <div className="flex flex-col gap-1">
            <label className="text-sm font-semibold">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="border-2 border-black bg-white px-3 py-2 focus:outline-none"
            />
          </div>
          {loginMutation.isError && (
            <p className="border-2 border-black bg-red-300 px-3 py-2 text-sm font-semibold">
              Invalid credentials
            </p>
          )}
          <button
            type="submit"
            disabled={loginMutation.isPending}
            className="hover-black border-2 border-black bg-white px-6 py-2 font-semibold disabled:opacity-50"
          >
            {loginMutation.isPending ? 'Logging in…' : 'Login'}
          </button>
        </form>
      </div>
    </div>
  )
}
