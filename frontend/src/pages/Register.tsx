import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { api, ApiError } from '../lib/api'
import { setLoggedIn } from '../lib/auth'
import { AuthShell } from '../components/AuthShell'
import { Alert, Button, Input } from '../components/ui'

export function Register() {
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await api.register(email, password)
      // The backend sets auth cookies on successful registration.
      setLoggedIn()
      navigate('/dashboard', { replace: true })
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Something went wrong. Try again.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <AuthShell
      title="Create your account"
      subtitle="Start tracking prices in seconds."
      footer={
        <>
          Already have an account?{' '}
          <Link to="/login" className="font-bold text-accent underline">
            Log in
          </Link>
        </>
      }
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Email"
          type="email"
          autoComplete="email"
          required
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="you@example.com"
        />
        <Input
          label="Password"
          type="password"
          autoComplete="new-password"
          required
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Choose a password"
        />
        <Alert>{error}</Alert>
        <Button type="submit" loading={loading} className="w-full">
          Sign up
        </Button>
      </form>
    </AuthShell>
  )
}
