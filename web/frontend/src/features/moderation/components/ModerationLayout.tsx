import { JSX, ReactNode } from 'react'
import { NavLink, To } from 'react-router'

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

export default function ModerationLayout({ children }: { children: ReactNode }) {
  return (
    <div className="dashboard-layout">
      <aside className="dashboard-layout__sidebar">
        <Logo />
        <ul className="dashboard-sidebar-list">
          {navigation.map((item) => (
            <SidebarItem
              key={item.to}
              icon={<i className={item.icon} />}
              label={window._(item.label)}
              to={item.to}
            />
          ))}
        </ul>
      </aside>

      <div className="dashboard-layout__body">{children}</div>
    </div>
  )
}

function Logo() {
  return (
    <div className="flex justify-center my-4">
      <img className="h-20 dark:hidden" src="/_/embed-assets/logo.svg" alt="" />
      <img className="h-20 hidden dark:block" src="/_/embed-assets/logo-dark.svg" alt="" />
    </div>
  )
}

function SidebarItem({ icon, label, to }: { icon?: JSX.Element; label: string; to: To }) {
  return (
    <li className="dashboard-sidebar-item">
      <NavLink to={to} className="dashboard-sidebar-item__container">
        <div className="dashboard-sidebar-item__icon">{icon}</div>
        <div className="dashboard-sidebar-item__label">{label}</div>
      </NavLink>
    </li>
  )
}
