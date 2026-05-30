import {
  HeadContent,
  Outlet,
  Scripts,
  createRootRouteWithContext,
  redirect,
  useRouterState,
} from '@tanstack/react-router'
import type { QueryClient } from '@tanstack/react-query'
import { Sidebar } from '../components/Sidebar'
import { isTokenValid } from '../lib/api'
import appCss from '../styles.css?url'

interface MyRouterContext {
  queryClient: QueryClient
}

export const Route = createRootRouteWithContext<MyRouterContext>()({
  head: () => ({
    meta: [
      { charSet: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { title: 'Kalleprint Admin' },
    ],
    links: [{ rel: 'stylesheet', href: appCss }],
  }),
  beforeLoad: ({ location }) => {
    if (typeof window === 'undefined') return
    const isLoginPage = location.pathname === '/login'
    if (!isLoginPage && !isTokenValid()) {
      throw redirect({ to: '/login' })
    }
    if (isLoginPage && isTokenValid()) {
      throw redirect({ to: '/' })
    }
  },
  shellComponent: RootDocument,
  component: RootLayout,
})

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body>
        {children}
        <Scripts />
      </body>
    </html>
  )
}

function RootLayout() {
  const pathname = useRouterState({ select: (s) => s.location.pathname })
  if (pathname === '/login') return <Outlet />
  return (
    <div className="flex h-screen">
      <Sidebar />
      <main className="flex-1 overflow-y-auto p-8">
        <Outlet />
      </main>
    </div>
  )
}
