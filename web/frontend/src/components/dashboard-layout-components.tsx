import clsx from 'clsx'
import { HTMLAttributes, ReactNode } from 'react'

function DashboardContent_Root({ children }: { children: ReactNode }) {
  return (
    <section className="DashboardContent" data-testid="DashboardContent_Root">
      {children}
    </section>
  )
}

function DashboardContent_Card({ children, className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div className={clsx('Card', className)} data-testid="DashboardContent_Card" {...props}>
      {children}
    </div>
  )
}

function DashboardContent_Header({
  title,
  children,
}: {
  title: ReactNode
  children?: ReactNode
}) {
  return (
    <header className="PageHeader" data-testid="DashboardContent_Header">
      <div className="PageHeader-body">
        <h1 className="PageHeader-title">{title}</h1>
      </div>
      {children && <div className="PageHeader-actions">{children}</div>}
    </header>
  )
}

export const DashboardContent = {
  Root: DashboardContent_Root,
  Card: DashboardContent_Card,
  Header: DashboardContent_Header,
}
