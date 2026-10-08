import { ReactNode, useLayoutEffect } from 'react'
import { NavLink, useLocation, useNavigation, useRevalidator } from 'react-router'
import { ManagerSkeleton } from './loading'
import './loading.scss'

export default function BMLayout({ children }: { children: ReactNode }) {
  const navigation = useNavigation()
  const location = useLocation()
  const revalidator = useRevalidator()
  const destination = navigation.location
  const changingDestination =
    navigation.state === 'loading' &&
    !!destination &&
    (destination.pathname !== location.pathname || destination.search !== location.search)
  const busy = navigation.state === 'loading' || revalidator.state === 'loading'
  // Dialogs render in document.body, outside the hidden route subtree.
  useLayoutEffect(() => {
    if (!changingDestination) return
    const dialogs = Array.from(document.querySelectorAll<HTMLDialogElement>('.BM-dialog'))
    dialogs.forEach((dialog) => {
      dialog.inert = true
      dialog.classList.add('BM-dialog--loading')
    })
    return () =>
      dialogs.forEach((dialog) => {
        dialog.inert = false
        dialog.classList.remove('BM-dialog--loading')
      })
  }, [changingDestination])
  return (
    <div className="BM DashboardShell">
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
          <NavLink
            to="/books"
            className={({ isActive }) =>
              isActive && location.pathname !== '/books/new' ? 'active' : ''
            }
            aria-current={location.pathname === '/books/new' ? false : undefined}
          >
            <i className="fa-solid fa-book" aria-hidden="true" />
            {window._('bookManager.books.title')}
          </NavLink>
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
          {changingDestination ? (
            <ManagerSkeleton location={destination} />
          ) : (
            busy && <p role="status">{window._('bookManager.ui.loading')}</p>
          )}
          <div hidden={changingDestination} inert={changingDestination}>
            {children}
          </div>
        </main>
      </div>
    </div>
  )
}
