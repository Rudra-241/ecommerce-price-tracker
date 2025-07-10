import { clearLoggedIn } from './auth'
import type { Product, PriceHistoryResponse } from '../types'

export class ApiError extends Error {
  status: number
  constructor(message: string, status: number) {
    super(message)
    this.status = status
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const res = await fetch(`/api${path}`, {
    credentials: 'include',
    headers: { 'Content-Type': 'application/json', ...(options.headers ?? {}) },
    ...options,
  })

  // Auth expired / missing on a protected route: drop the hint and bounce to login.
  if (res.status === 401) {
    clearLoggedIn()
    if (window.location.pathname !== '/login') {
      window.location.assign('/login')
    }
    throw new ApiError('Session expired. Please log in again.', 401)
  }

  // Some endpoints (logout) may return no body.
  const text = await res.text()
  const data = text ? JSON.parse(text) : {}

  if (!res.ok) {
    throw new ApiError(data?.error ?? `Request failed (${res.status})`, res.status)
  }
  return data as T
}

export const api = {
  login: (email: string, password: string) =>
    request<{ message: string }>('/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  register: (email: string, password: string) =>
    request<{ message: string }>('/register', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    }),

  // Endpoint is added in the backend phase; tolerate it not existing yet.
  logout: () => request<unknown>('/logout', { method: 'POST' }).catch(() => undefined),

  getProducts: () =>
    request<{ tracked_products: Product[] | null }>('/products').then(
      (r) => r.tracked_products ?? [],
    ),

  addProduct: (url: string) =>
    request<{ message: string }>('/product', {
      method: 'POST',
      body: JSON.stringify({ url }),
    }),

  getProduct: (id: number | string) => request<PriceHistoryResponse>(`/product/${id}`),

  deleteProduct: (id: number | string) =>
    request<{ message: string }>(`/product/${id}`, { method: 'DELETE' }),
}
