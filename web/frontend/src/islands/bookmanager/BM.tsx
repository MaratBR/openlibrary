import { dashboardItem } from '@/components/dashboard-nav-item'
import { ReactIslandProps } from '../common/react-island'
import { createHashRouter, Navigate, Outlet, NavLink, useParams } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import BMLayout from './BMLayout'
import { Books, booksRouteLoader } from './books'
import { Book, bookRouteLoader } from './books/Book'
import BookEdit from './books/BookEdit'
import { usePageTitle } from './ui'
import { ManagerSkeleton } from './loading'
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
        element: <BookRoute />,
        loader: bookRouteLoader,
        handle: dashboardItem<Awaited<ReturnType<typeof bookRouteLoader>>>(
          (data) => data.bookResponse.data.name,
        ),
        errorElement: <RouteError />,
        hydrateFallbackElement: <Loading />,
      },
      {
        path: '/books/:bookId/edit',
        element: <BookEdit />,
        loader: bookRouteLoader,
        handle: dashboardItem<Awaited<ReturnType<typeof bookRouteLoader>>>(
          (data) => data.bookResponse.data.name,
        ),
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

// A new book must not inherit the previous book's open chapter panel or local state.
function BookRoute() {
  const { bookId } = useParams()
  return <Book key={bookId} />
}

function Loading() {
  return <ManagerSkeleton />
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
