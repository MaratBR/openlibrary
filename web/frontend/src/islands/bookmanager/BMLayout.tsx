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
    <div className="BM">
      <header className="BM-masthead">
        <NavLink to="/books" aria-label={window._('bookManager.title')}>
          <img className="h-12 dark:hidden" src="/_/embed-assets/logo.svg" alt="" />
          <img className="h-12 hidden dark:block" src="/_/embed-assets/logo-dark.svg" alt="" />
        </NavLink>
        <nav className="flex flex-wrap gap-2" aria-label={window._('bookManager.title')}>
          <a href="/" className="Btn Btn--ghost" data-manager-site-link>
            {window._('bookManager.ui.mainSite')}
          </a>
          <NavLink to="/books" className="Btn Btn--ghost">
            {window._('bookManager.books.title')}
          </NavLink>
        </nav>
      </header>
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
  )
}
