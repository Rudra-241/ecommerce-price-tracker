import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { api, ApiError } from '../lib/api'
import type { PriceHistoryResponse } from '../types'
import { formatDateTime, formatPrice } from '../lib/format'
import { PriceChart } from '../components/PriceChart'
import { Alert, Spinner } from '../components/ui'

export function ProductDetail() {
  const { id } = useParams<{ id: string }>()
  const [data, setData] = useState<PriceHistoryResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    if (!id) return
    let active = true
    setLoading(true)
    api
      .getProduct(id)
      .then((res) => active && setData(res))
      .catch((err) => {
        if (active && err instanceof ApiError && err.status !== 401) setError(err.message)
      })
      .finally(() => active && setLoading(false))
    return () => {
      active = false
    }
  }, [id])

  const prices = data?.price_history.map((s) => s.Price) ?? []
  const lowest = prices.length ? Math.min(...prices) : null
  const highest = prices.length ? Math.max(...prices) : null
  const lastChecked = data?.price_history.at(-1)?.ChangedAt

  return (
    <div className="space-y-6">
      <Link to="/dashboard" className="inline-flex items-center gap-1 font-bold text-secondary">
        ← Back to dashboard
      </Link>

      {loading ? (
        <div className="flex items-center gap-3 text-text/60">
          <Spinner /> Loading price history…
        </div>
      ) : error ? (
        <Alert>{error}</Alert>
      ) : data ? (
        <>
          <div className="card p-6">
            <h4 className="mb-4 break-words">{data.product_name}</h4>
            <div className="grid grid-cols-2 gap-4 sm:grid-cols-3">
              <Stat label="Current price" value={formatPrice(data.current_price)} highlight />
              {lowest !== null && <Stat label="Lowest seen" value={formatPrice(lowest)} />}
              {highest !== null && <Stat label="Highest seen" value={formatPrice(highest)} />}
            </div>
            {lastChecked && (
              <p className="mt-4 text-sm text-text/50">
                Last checked {formatDateTime(lastChecked)}
              </p>
            )}
          </div>

          <div className="card p-6">
            <h5 className="mb-4">Price history</h5>
            <PriceChart history={data.price_history} />
          </div>
        </>
      ) : null}
    </div>
  )
}

function Stat({
  label,
  value,
  highlight = false,
}: {
  label: string
  value: string
  highlight?: boolean
}) {
  return (
    <div
      className={`rounded-xl border-2 border-text p-4 ${highlight ? 'bg-primary' : 'bg-white'}`}
    >
      <div className="text-sm font-bold text-text/60">{label}</div>
      <div className="text-2xl font-bold">{value}</div>
    </div>
  )
}
