import { dashboardItem } from '@/components/dashboard-nav-item'
import { DashboardNavItem } from '@/components/dashboard-nav-item'
import { useLayoutEffect } from 'react'
import { createHashRouter, Link, Outlet, useLocation, useNavigation } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { Users, usersLoader, UserEdit, userLoader } from '@/features/admin/pages/users'
import { Tags, tagsLoader, TagDetails, tagLoader, TagEdit } from '@/features/admin/pages/tags'
import { Home, Debug, debugLoader } from '@/features/admin/pages/screens'
import { Books, BookDetails, booksLoader, bookLoader } from '@/features/admin/pages/books'
import { Card, Loading, PageHeader, RouteError, t } from '@/features/admin/components/ui'

function Layout() {
  const navigation = useNavigation()
  const location = useLocation()
  const loading = navigation.state === 'loading'
  useLayoutEffect(() => {
    document.querySelector<HTMLElement>('.Admin-main h1')?.focus({ preventScroll: true })
    window.scrollTo({ top: 0 })
  }, [location.pathname])
  return (
    <div className="Admin DashboardShell">
      <aside className="DashboardShell-sidebar">
        <header className="DashboardShell-brand">
          <Link to="/" className="Admin-brand">
            <img src="/_/embed-assets/logo.svg" className="dark:hidden" alt="" />
            <img src="/_/embed-assets/logo-dark.svg" className="hidden dark:block" alt="" />
            <span>
              <strong>{t('admin.ui.title')}</strong>
              <small>{t('admin.ui.subtitle')}</small>
            </span>
          </Link>
        </header>
        <nav className="DashboardShell-nav" aria-label={t('admin.ui.section')}>
          {[
            ['/', 'home', 'fa-house'],
            ['/users', 'users', 'fa-users'],
            ['/books', 'books', 'fa-book'],
            ['/tags', 'tags', 'fa-tags'],
            ['/debug', 'debugActions', 'fa-wrench'],
          ].map(([path, key, icon]) => (
            <DashboardNavItem dashboard="admin" end={path === '/'} key={path} to={path}>
              <i className={`fa-solid ${icon}`} aria-hidden="true" />
              {t(`admin.sidebar.${key}`)}
            </DashboardNavItem>
          ))}
        </nav>
        <div className="DashboardShell-utilities">
          <div className="Admin-actions">
            <button
              type="button"
              className="Btn Btn--ghost"
              aria-label={t('admin.ui.theme')}
              onClick={() => window.OLTheme.toggle()}
            >
              <i className="fa-solid fa-circle-half-stroke" aria-hidden="true" />
              {t('admin.ui.theme')}
            </button>
            <a className="Btn Btn--ghost" data-admin-external href="/">
              <i className="fa-solid fa-arrow-left" aria-hidden="true" />
              {t('admin.goBackToSite')}
            </a>
            <a className="Btn Btn--outline" data-admin-external href="/logout">
              <i className="fa-solid fa-arrow-right-from-bracket" aria-hidden="true" />
              {t('admin.sidebar.logout')}
            </a>
          </div>
        </div>
      </aside>
      <div className="Admin-frame DashboardShell-content">
        <main className="Admin-main" id="admin-main" aria-busy={loading}>
          {loading ? <Loading /> : <Outlet />}
        </main>
        <footer className="Admin-footer">
          OpenLibrary <span>·</span> {t('admin.ui.title')}
        </footer>
      </div>
    </div>
  )
}
function NotFound() {
  return (
    <>
      <PageHeader title={t('admin.ui.notFound')} />
      <Card>
        <p>{t('admin.ui.notFoundDescription')}</p>
        <Link className="Btn Btn--primary mt-4" to="/">
          {t('admin.ui.backHome')}
        </Link>
      </Card>
    </>
  )
}
const router = createHashRouter([
  {
    path: '/',
    element: <Layout />,
    children: [
      { index: true, element: <Home /> },
      { path: 'users', element: <Users />, loader: usersLoader, errorElement: <RouteError /> },
      {
        path: 'users/:id',
        element: <UserEdit />,
        loader: userLoader,
        handle: dashboardItem<Awaited<ReturnType<typeof userLoader>>>((data) => data.name),
        errorElement: <RouteError />,
        hydrateFallbackElement: <Loading />,
      },
      { path: 'books', element: <Books />, loader: booksLoader, errorElement: <RouteError /> },
      {
        path: 'books/:id',
        element: <BookDetails />,
        loader: bookLoader,
        handle: dashboardItem<Awaited<ReturnType<typeof bookLoader>>>((data) => data.name),
        errorElement: <RouteError />,
        hydrateFallbackElement: <Loading />,
      },
      { path: 'tags', element: <Tags />, loader: tagsLoader, errorElement: <RouteError /> },
      {
        path: 'tags/:id',
        element: <TagDetails />,
        loader: tagLoader,
        handle: dashboardItem<Awaited<ReturnType<typeof tagLoader>>>((data) => data.name),
        errorElement: <RouteError />,
        hydrateFallbackElement: <Loading />,
      },
      {
        path: 'tags/:id/edit',
        element: <TagEdit />,
        loader: tagLoader,
        handle: dashboardItem<Awaited<ReturnType<typeof tagLoader>>>((data) => data.name),
        errorElement: <RouteError />,
        hydrateFallbackElement: <Loading />,
      },
      { path: 'debug', element: <Debug />, loader: debugLoader, errorElement: <RouteError /> },
      { path: '*', element: <NotFound /> },
    ],
  },
])
export default function App() {
  return <RouterProvider router={router} />
}
