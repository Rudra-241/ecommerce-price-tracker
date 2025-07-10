import type { Direction } from '../types'

const inr = new Intl.NumberFormat('en-IN', {
  style: 'currency',
  currency: 'INR',
  maximumFractionDigits: 2,
})

export function formatPrice(value: number): string {
  return inr.format(value)
}

export function formatDate(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleDateString('en-IN', { day: 'numeric', month: 'short', year: 'numeric' })
}

export function formatDateTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return iso
  return d.toLocaleString('en-IN', {
    day: 'numeric',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
  })
}

interface DirectionMeta {
  label: string
  /** Tailwind classes for the badge (bg + text). */
  badge: string
}

// A price drop is good news (greens), a rise is bad (warm), flat is neutral.
export function directionMeta(direction: Direction): DirectionMeta {
  switch (direction) {
    case 'BelowStart':
      return { label: 'Below start price', badge: 'bg-primary text-text' }
    case 'Decreased':
      return { label: 'Price dropped', badge: 'bg-primary-soft text-text' }
    case 'Increased':
      return { label: 'Price rose', badge: 'bg-accent-soft text-accent' }
    default:
      return { label: 'Unchanged', badge: 'bg-text/5 text-text/60' }
  }
}
