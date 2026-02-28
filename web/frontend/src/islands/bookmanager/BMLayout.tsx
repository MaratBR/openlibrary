import { ComponentChild, JSX } from 'preact'
import { NavLink, To } from 'react-router'

export default function BMLayout({ children }: { children: ComponentChild }) {
  return (
    <div class="dashboard-layout">
      <aside class="dashboard-layout__sidebar">
        <Logo />
        <hr class="dashboard-sidebar-hr my-3" />
        <ul class="dashboard-sidebar-list">
          <SidebarItem
            icon={<i class="fa-solid fa-book" />}
            label={window._('bookManager.books.title')}
            to="/books"
          />
        </ul>
      </aside>

      <div class="dashboard-layout__body">{children}</div>
    </div>
  )
}

function Logo() {
  return (
    <div class="flex justify-center my-4">
      <img class="h-20" src="/_/embed-assets/logo.svg" />
    </div>
  )
}

function SidebarItem({ icon, label, to }: { icon?: JSX.Element; label: string; to: To }) {
  return (
    <li class="dashboard-sidebar-item">
      <NavLink to={to}>
        <div class="dashboard-sidebar-item__icon">{icon}</div>
        <div class="dashboard-sidebar-item__label">{label}</div>
      </NavLink>
    </li>
  )
}
