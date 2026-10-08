import { DashboardNavItem } from '@/components/dashboard-nav-item'
import { ReactNode } from 'react'
import { NavLink } from 'react-router'

const navigation = [
  { to: '/overview', icon: 'fa-solid fa-gauge-high', label: 'moderationPortal.overview' },
  { to: '/search', icon: 'fa-solid fa-search', label: 'moderationPortal.search' },
  { to: '/books', icon: 'fa-solid fa-book', label: 'moderationPortal.books' },
  { to: '/chapters', icon: 'fa-solid fa-file-lines', label: 'moderationPortal.chapters' },
  { to: '/comments', icon: 'fa-solid fa-comments', label: 'moderationPortal.comments' },
  { to: '/users', icon: 'fa-solid fa-users', label: 'moderationPortal.users' },
  { to: '/reports', icon: 'fa-solid fa-flag', label: 'moderationPortal.reports' },
  {
    to: '/login-history',
    icon: 'fa-solid fa-clock-rotate-left',
    label: 'moderationPortal.loginHistory',
  },
  { to: '/audit-log', icon: 'fa-solid fa-clipboard-list', label: 'moderationPortal.auditLog' },
] as const

export default function ModerationLayout({
  children,
  isAdmin,
}: {
  children: ReactNode
  isAdmin: boolean
}) {
  return (
    <div className="DashboardShell">
      <aside className="DashboardShell-sidebar">
        <header className="DashboardShell-brand">
          <NavLink className="DashboardShell-brandLink" to="/overview">
            <img className="dark:hidden" src="/_/embed-assets/logo.svg" alt="" />
            <img className="hidden dark:block" src="/_/embed-assets/logo-dark.svg" alt="" />
            <strong>{window._('moderationPortal.title')}</strong>
          </NavLink>
        </header>
        <nav className="DashboardShell-nav" aria-label={window._('moderationPortal.title')}>
          {navigation.map((item) => (
            <DashboardNavItem dashboard="moderation" key={item.to} to={item.to}>
              <i className={item.icon} aria-hidden="true" />
              {window._(item.label)}
            </DashboardNavItem>
          ))}
        </nav>
        <div className="DashboardShell-utilities flex flex-col gap-2">
          {isAdmin && (
            <a className="Btn Btn--outline" href="/admin">
              {window._('moderationPortal.adminDashboard')}
            </a>
          )}
          <a className="Btn Btn--ghost" href="/">
            {window._('moderationPortal.backToSite')}
          </a>
          <button className="Btn Btn--ghost" onClick={() => window.OLTheme.toggle()}>
            {window._('moderationPortal.theme')}
          </button>
        </div>
      </aside>
      <main className="DashboardShell-content p-4 md:p-8">{children}</main>
    </div>
  )
}
