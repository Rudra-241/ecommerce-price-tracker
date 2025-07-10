import { useState } from 'react'
import { Link } from 'react-router-dom'
import type { Product } from '../types'
import { directionMeta, formatPrice } from '../lib/format'
import { Badge, Spinner } from './ui'

interface ProductCardProps {
  product: Product
  onUntrack: (id: number) => Promise<void>
}

export function ProductCard({ product, onUntrack }: ProductCardProps) {
  const [removing, setRemoving] = useState(false)
  const meta = directionMeta(product.Direction)

  async function handleUntrack() {
    setRemoving(true)
    try {
      await onUntrack(product.ID)
    } finally {
      setRemoving(false)
    }
  }

  return (
    <div className="card flex flex-col overflow-hidden">
      <Link to={`/product/${product.ID}`} className="block">
        <div className="grid aspect-video place-items-center border-b-2 border-text bg-white p-4">
          {product.ImgLink ? (
            <img
              src={product.ImgLink}
              alt={product.Name}
              className="max-h-full max-w-full object-contain"
              loading="lazy"
            />
          ) : (
            <span className="text-5xl">📦</span>
          )}
        </div>
      </Link>

      <div className="flex flex-1 flex-col gap-3 p-4">
        <Link
          to={`/product/${product.ID}`}
          className="line-clamp-2 font-bold leading-snug hover:underline"
          title={product.Name}
        >
          {product.Name}
        </Link>

        <div className="mt-auto flex items-end justify-between gap-2">
          <span className="text-2xl font-bold">{formatPrice(product.Price)}</span>
          <Badge className={meta.badge}>{meta.label}</Badge>
        </div>

        <div className="flex items-center justify-between gap-2 border-t-2 border-text/10 pt-3">
          <Link to={`/product/${product.ID}`} className="text-sm font-bold text-secondary underline">
            View history
          </Link>
          <button
            onClick={handleUntrack}
            disabled={removing}
            className="inline-flex items-center gap-1.5 text-sm font-bold text-text/50 transition hover:text-accent disabled:opacity-50"
          >
            {removing ? <Spinner /> : null}
            Untrack
          </button>
        </div>
      </div>
    </div>
  )
}
