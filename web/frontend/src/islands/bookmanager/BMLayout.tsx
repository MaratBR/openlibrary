import { JSX, ReactNode } from 'react'
import { NavLink, To } from 'react-router'

export default function BMLayout({ children }: { children: ReactNode }) {
  return (
    <div className="dashboard-layout">
      <aside className="dashboard-layout__sidebar">
        <Logo />
        <ul className="dashboard-sidebar-list">
          <SidebarItem
            icon={<i className="fa-solid fa-book" />}
            label={window._('bookManager.books.title')}
            to="/books"
          />
        </ul>
      </aside>

      <div className="dashboard-layout__body">{children}</div>
    </div>
  )
}

function Logo() {
  return (
    <div className="flex justify-center my-4">
      <img className="h-20 dark:hidden" src="/_/embed-assets/logo.svg" />
      <img className="h-20 hidden dark:block" src="/_/embed-assets/logo-dark.svg" />
    </div>
  )
}

function SidebarItem({ icon, label, to }: { icon?: JSX.Element; label: string; to: To }) {
  return (
    <li className="dashboard-sidebar-item">
      <NavLink to={to}>
        <div className="dashboard-sidebar-item__icon">{icon}</div>
        <div className="dashboard-sidebar-item__label">{label}</div>
      </NavLink>
    </li>
  )
}
