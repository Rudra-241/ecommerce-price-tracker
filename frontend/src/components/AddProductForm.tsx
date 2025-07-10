import { useState, type FormEvent } from 'react'
import { api, ApiError } from '../lib/api'
import { Alert, Button } from './ui'

export function AddProductForm({ onAdded }: { onAdded: () => void }) {
  const [url, setUrl] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await api.addProduct(url.trim())
      setUrl('')
      onAdded()
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Could not track that product.')
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="card p-5">
      <label htmlFor="product-url" className="mb-2 block font-bold">
        Track a new product
      </label>
      <div className="flex flex-col gap-3 sm:flex-row">
        <input
          id="product-url"
          type="url"
          required
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="Paste a product URL…"
          className="field flex-1"
        />
        <Button type="submit" loading={loading} className="shrink-0">
          {loading ? 'Adding…' : 'Track it'}
        </Button>
      </div>
      {error && (
        <div className="mt-3">
          <Alert>{error}</Alert>
        </div>
      )}
    </form>
  )
}
