import clsx from 'clsx'
import { ComponentChild, HTMLAttributes } from 'preact'

function DashboardContent_Root({ children }: { children: ComponentChild }) {
  return (
    <section class="dashboard-content" data-testid="DashboardContent_Root">
      {children}
    </section>
  )
}

function DashboardContent_Card({
  children,
  class: clazz,
  className,
  ...props
}: HTMLAttributes<HTMLDivElement>) {
  return (
    <div className={clsx('card', clazz, className)} data-testid="DashboardContent_Card" {...props}>
      {children}
    </div>
  )
}

function DashboardContent_StickyHeader({
  title,
  children,
}: {
  title: ComponentChild
  children?: ComponentChild
}) {
  return (
    <div class="dashboard-content__sticky-header" data-testid="DashboardContent_StickyHeader">
      <header class="page-header-container">
        <h1 class="page-header">{title}</h1>
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
