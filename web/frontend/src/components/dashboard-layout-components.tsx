import { ComponentChild } from 'preact'

function DashboardContent_Root({ children }: { children: ComponentChild }) {
  return <section class="dashboard-content">{children}</section>
}

function DashboardContent_StickyHeader({ title }: { title: string }) {
  return (
    <div class="dashboard-content__sticky-header">
      <header class="page-header-container">
        <h1 class="page-header">{title}</h1>
      </header>
    </div>
  )
}

export const DashboardContent = {
  Root: DashboardContent_Root,
  StickyHeader: DashboardContent_StickyHeader,
}
