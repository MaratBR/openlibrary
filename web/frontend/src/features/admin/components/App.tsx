import { useLayoutEffect } from 'react'
import { createHashRouter, Link, NavLink, Outlet, useLocation, useNavigation } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { Users, usersLoader, UserEdit, userLoader } from '@/features/admin/pages/users'
import { Tags, tagsLoader, TagDetails, tagLoader, TagEdit } from '@/features/admin/pages/tags'
import { Home, Books, Debug, debugLoader } from '@/features/admin/pages/screens'
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
    <div className="Admin">
      <header className="Admin-masthead">
        <Link to="/" className="Admin-brand">
          <img src="/_/embed-assets/logo.svg" className="dark:hidden" alt="" />
          <img src="/_/embed-assets/logo-dark.svg" className="hidden dark:block" alt="" />
          <span>
            <strong>{t('admin.ui.title')}</strong>
            <small>{t('admin.ui.subtitle')}</small>
          </span>
        </Link>
        <div className="Admin-actions">
          <button
            type="button"
            className="Btn Btn--ghost Btn--icon"
            aria-label={t('admin.ui.theme')}
            onClick={() => window.OLTheme.toggle()}
          >
            <i className="fa-solid fa-circle-half-stroke" aria-hidden="true" />
          </button>
          <a className="Btn Btn--ghost" data-admin-external href="/">
            {t('admin.goBackToSite')}
          </a>
          <a className="Btn Btn--outline" data-admin-external href="/logout">
            {t('admin.sidebar.logout')}
          </a>
        </div>
      </header>
      <div className="Admin-frame">
        <nav className="Admin-nav" aria-label={t('admin.ui.section')}>
          {[
            ['/', 'home', 'fa-house'],
            ['/users', 'users', 'fa-users'],
            ['/books', 'books', 'fa-book'],
            ['/tags', 'tags', 'fa-tags'],
            ['/debug', 'debugActions', 'fa-wrench'],
          ].map(([path, key, icon]) => (
            <NavLink end={path === '/'} key={path} to={path}>
              <i className={`fa-solid ${icon}`} aria-hidden="true" />
              {t(`admin.sidebar.${key}`)}
            </NavLink>
          ))}
        </nav>
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
        errorElement: <RouteError />,
        hydrateFallbackElement: <Loading />,
      },
      { path: 'books', element: <Books /> },
      { path: 'tags', element: <Tags />, loader: tagsLoader, errorElement: <RouteError /> },
      {
        path: 'tags/:id',
        element: <TagDetails />,
        loader: tagLoader,
        errorElement: <RouteError />,
        hydrateFallbackElement: <Loading />,
      },
      {
        path: 'tags/:id/edit',
        element: <TagEdit />,
        loader: tagLoader,
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
