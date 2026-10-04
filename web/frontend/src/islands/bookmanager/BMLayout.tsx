import { ReactNode } from 'react'
import { NavLink, useNavigation } from 'react-router'

export default function BMLayout({ children }: { children: ReactNode }) {
  const navigation = useNavigation()
  return (
    <div className="BM">
      <header className="BM-masthead">
        <NavLink to="/books" aria-label={window._('bookManager.title')}>
          <img className="h-12 dark:hidden" src="/_/embed-assets/logo.svg" alt="" />
          <img className="h-12 hidden dark:block" src="/_/embed-assets/logo-dark.svg" alt="" />
        </NavLink>
        <NavLink to="/books" className="Btn Btn--ghost">
          {window._('bookManager.books.title')}
        </NavLink>
      </header>
      <main className="BM-main" aria-busy={navigation.state === 'loading'}>
        {navigation.state === 'loading' && (
          <p role="status">{window._('bookManager.ui.loading')}</p>
        )}
        {children}
      </main>
    </div>
  )
}
