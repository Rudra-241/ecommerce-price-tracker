import { Link } from 'react-router-dom'
import { Logo } from '../components/Logo'

// Placeholder landing page — to be designed later.
export function Landing() {
  return (
    <div className="flex min-h-screen flex-col">
      <header className="mx-auto flex w-full max-w-5xl items-center justify-between px-4 py-5">
        <Logo />
        <div className="flex items-center gap-3">
          <Link to="/login" className="btn btn-ghost">
            Log in
          </Link>
          <Link to="/register" className="btn btn-primary">
            Sign up
          </Link>
        </div>
      </header>

      <main className="mx-auto flex w-full max-w-3xl flex-1 flex-col items-center justify-center px-4 text-center">
        <span className="mb-6 inline-flex items-center rounded-full border-2 border-text bg-secondary-soft px-4 py-1 text-sm font-bold">
          🚧 Landing page coming soon
        </span>
        <h1 className="mb-4">Track prices. Catch the drop.</h1>
        <p className="mb-8 max-w-xl text-lg text-text/70">
          PriceWatch keeps an eye on product prices for you and pings you the moment they fall.
          The full landing experience is still being built — jump straight in for now.
        </p>
        <div className="flex flex-wrap items-center justify-center gap-4">
          <Link to="/register" className="btn btn-primary">
            Get started
          </Link>
          <Link to="/login" className="btn btn-secondary">
            I already have an account
          </Link>
        </div>
      </main>

      <footer className="mx-auto w-full max-w-5xl px-4 py-6 text-center text-sm text-text/40">
        PriceWatch
      </footer>
    </div>
  )
}
