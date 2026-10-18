import { type ReactNode } from 'react'
import { HeaderActions } from '@/components/header-actions'

/** Shared shell placement keeps account controls outside page-specific padding. */
export function DashboardHeader({ children }: { children: ReactNode }) {
  return (
    <header className="DashboardShell-header">
      {children}
      <HeaderActions />
    </header>
  )
}
