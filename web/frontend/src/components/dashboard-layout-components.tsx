import clsx from 'clsx'
import {  HTMLAttributes, ReactNode } from 'react'

function DashboardContent_Root({ children }: { children: ReactNode }) {
  return (
    <section className="dashboard-content" data-testid="DashboardContent_Root">
      {children}
    </section>
  )
}

function DashboardContent_Card({
  children,
  className,
  ...props
}: HTMLAttributes<HTMLDivElement>) {
  return (
    <div className={clsx('card', className)} data-testid="DashboardContent_Card" {...props}>
      {children}
    </div>
  )
}

function DashboardContent_StickyHeader({
  title,
  children,
}: {
  title: ReactNode
  children?: ReactNode
}) {
  return (
    <div className="dashboard-content__sticky-header" data-testid="DashboardContent_StickyHeader">
      <header className="page-header-container">
        <h1 className="page-header">{title}</h1>
        {children}
      </header>
    </div>
  )
}

export const DashboardContent = {
  Root: DashboardContent_Root,
  Card: DashboardContent_Card,
  StickyHeader: DashboardContent_StickyHeader,
}
