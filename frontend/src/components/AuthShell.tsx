import type { ReactNode } from 'react'
import { Logo } from './Logo'

interface AuthShellProps {
  title: string
  subtitle: string
  children: ReactNode
  footer: ReactNode
}

// Centered card layout shared by the Login and Register screens.
export function AuthShell({ title, subtitle, children, footer }: AuthShellProps) {
  return (
    <div className="flex min-h-screen flex-col items-center justify-center px-4 py-10">
      <div className="mb-8">
        <Logo />
      </div>
      <div className="card w-full max-w-md p-8">
        <h4 className="mb-1">{title}</h4>
        <p className="mb-6 text-text/60">{subtitle}</p>
        {children}
      </div>
      <p className="mt-6 text-sm text-text/60">{footer}</p>
    </div>
  )
}
