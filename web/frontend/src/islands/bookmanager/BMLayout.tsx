import { DashboardNavItem } from '@/components/dashboard-nav-item'
import { DashboardLoader } from '@/components/dashboard-loader'
import { ReactNode } from 'react'
import { NavLink, useNavigation, useRevalidator } from 'react-router'

export default function BMLayout({ children }: { children: ReactNode }) {
  const navigation = useNavigation()
  const revalidator = useRevalidator()
  const busy = navigation.state !== 'idle' || revalidator.state === 'loading'
  return (
    <div className="BM DashboardShell">
      {busy && <DashboardLoader label={window._('bookManager.ui.loading')} />}
      <aside className="DashboardShell-sidebar">
        <header className="DashboardShell-brand">
          <NavLink to="/books" className="DashboardShell-brandLink">
            <img className="dark:hidden" src="/_/embed-assets/logo.svg" alt="" />
            <img className="hidden dark:block" src="/_/embed-assets/logo-dark.svg" alt="" />
            <span>
              <strong>{window._('bookManager.title')}</strong>
              <small>OpenLibrary</small>
            </span>
          </NavLink>
        </header>
        <nav className="DashboardShell-nav" aria-label={window._('bookManager.title')}>
          <DashboardNavItem dashboard="bookmanager" to="/books" exclude="/books/new">
            <i className="fa-solid fa-book" aria-hidden="true" />
            {window._('bookManager.books.title')}
          </DashboardNavItem>
          <NavLink to="/books/new">
            <i className="fa-solid fa-plus" aria-hidden="true" />
            {window._('bookManager.books.addBook')}
          </NavLink>
        </nav>
        <div className="DashboardShell-utilities">
          <a href="/" className="Btn Btn--ghost" data-manager-site-link>
            <i className="fa-solid fa-arrow-left" aria-hidden="true" />
            {window._('bookManager.ui.mainSite')}
          </a>
        </div>
      </aside>
      <div className="DashboardShell-content">
        <main className="BM-main" aria-busy={busy}>
          {children}
        </main>
      </div>
    </div>
  )
}
