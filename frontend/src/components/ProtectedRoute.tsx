import { Navigate, Outlet } from 'react-router-dom'
import { isLoggedIn } from '../lib/auth'
import { Nav } from './Nav'

// Gate for authenticated pages. The localStorage flag gives an instant routing
// decision; the server is still the source of truth (an /api 401 redirects via
// the api client). Renders the shared app chrome (Nav) around the page.
export function ProtectedRoute() {
  if (!isLoggedIn()) {
    return <Navigate to="/login" replace />
  }
  return (
    <div className="min-h-screen">
      <Nav />
      <main className="mx-auto max-w-5xl px-4 py-8">
        <Outlet />
      </main>
    </div>
  )
}
