import { useNavigate } from 'react-router-dom'
import { api } from '../lib/api'
import { clearLoggedIn } from '../lib/auth'
import { Button } from './ui'
import { Logo } from './Logo'

export function Nav() {
  const navigate = useNavigate()

  async function handleLogout() {
    await api.logout()
    clearLoggedIn()
    navigate('/login', { replace: true })
  }

  return (
    <header className="sticky top-0 z-10 border-b-2 border-text bg-background/90 backdrop-blur">
      <div className="mx-auto flex max-w-5xl items-center justify-between px-4 py-3">
        <Logo to="/dashboard" />
        <Button variant="ghost" onClick={handleLogout}>
          Log out
        </Button>
      </div>
    </header>
  )
}
