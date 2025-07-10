import { useCallback, useEffect, useState } from 'react'
import { api, ApiError } from '../lib/api'
import type { Product } from '../types'
import { AddProductForm } from '../components/AddProductForm'
import { ProductCard } from '../components/ProductCard'
import { Alert, Spinner } from '../components/ui'

export function Dashboard() {
  const [products, setProducts] = useState<Product[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const load = useCallback(async () => {
    setError('')
    try {
      setProducts(await api.getProducts())
    } catch (err) {
      // 401s are handled globally (redirect); surface anything else.
      if (err instanceof ApiError && err.status !== 401) {
        setError(err.message)
      }
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void load()
  }, [load])

  async function handleUntrack(id: number) {
    // Optimistic removal; reload on failure to resync.
    const previous = products
    setProducts((p) => p.filter((x) => x.ID !== id))
    try {
      await api.deleteProduct(id)
    } catch {
      setProducts(previous)
    }
  }

  return (
    <div className="space-y-8">
      <div>
        <h3 className="mb-1">Your products</h3>
        <p className="text-text/60">Paste a link and we'll watch the price for you.</p>
      </div>

      <AddProductForm onAdded={load} />

      {error && <Alert>{error}</Alert>}

      {loading ? (
        <div className="flex items-center gap-3 text-text/60">
          <Spinner /> Loading your products…
        </div>
      ) : products.length === 0 ? (
        <EmptyState />
      ) : (
        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {products.map((product) => (
            <ProductCard key={product.ID} product={product} onUntrack={handleUntrack} />
          ))}
        </div>
      )}
    </div>
  )
}

function EmptyState() {
  return (
    <div className="card flex flex-col items-center gap-2 p-12 text-center">
      <span className="text-5xl">🛒</span>
      <h5>Nothing tracked yet</h5>
      <p className="max-w-sm text-text/60">
        Add your first product above to start watching its price.
      </p>
    </div>
  )
}
