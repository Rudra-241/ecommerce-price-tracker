import { Link } from 'react-router-dom'

export function Logo({ to = '/' }: { to?: string }) {
  return (
    <Link to={to} className="inline-flex items-center gap-2 text-text">
      <span className="grid size-9 place-items-center rounded-xl border-2 border-text bg-primary text-lg font-bold">
        ₹
      </span>
      <span className="text-xl font-bold tracking-tight">PriceWatch</span>
    </Link>
  )
}
