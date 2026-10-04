import { ReactIslandProps } from '../common/react-island'
import { createHashRouter, Navigate, Outlet, NavLink } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import BMLayout from './BMLayout'
import { Books, booksRouteLoader } from './books'
import { Book, bookRouteLoader } from './books/Book'
import BookEdit from './books/BookEdit'
import { usePageTitle } from './ui'
import NewBookForm from './new-book/NewBookForm'

const router = createHashRouter([
  {
    path: '/',
    element: (
      <BMLayout>
        <Outlet />
      </BMLayout>
    ),
    children: [
      {
        path: '/',
        element: <Navigate to="/books" replace />,
      },
      {
        path: '/books',
        element: <Books />,
        loader: booksRouteLoader,
        errorElement: <RouteError />,
        hydrateFallbackElement: <Loading />,
      },
      {
        path: '/books/new',
        element: <NewBookForm />,
      },
      {
        path: '/books/:bookId',
        element: <Book />,
        loader: bookRouteLoader,
        errorElement: <RouteError />,
        hydrateFallbackElement: <Loading />,
      },
      {
        path: '/books/:bookId/edit',
        element: <BookEdit />,
        loader: bookRouteLoader,
        errorElement: <RouteError />,
        hydrateFallbackElement: <Loading />,
      },
      {
        path: '*',
        element: <div>404</div>,
      },
    ],
  },
])

export default function BM(_props: ReactIslandProps) {
  return <RouterProvider router={router} />
}

function Loading() {
  return <p role="status">{window._('bookManager.ui.loading')}</p>
}
function RouteError() {
  usePageTitle(window._('bookManager.ui.unavailable'))
  return (
    <section className="Card my-6 space-y-4">
      <h1>{window._('bookManager.ui.unavailable')}</h1>
      <p role="alert">{window._('bookManager.ui.loadError')}</p>
      <div className="flex flex-wrap gap-3">
        <button className="Btn Btn--primary" onClick={() => window.location.reload()}>
          {window._('bookManager.ui.retry')}
        </button>
        <NavLink to="/books" className="Btn Btn--outline">
          {window._('bookManager.ui.back')}
        </NavLink>
      </div>
    </section>
  )
}
