import { PreactIslandProps } from '../common/preact-island'
import { createHashRouter, Navigate, Outlet } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import BMLayout from './BMLayout'
import { Books, booksRouteLoader } from './books'
import { Book, bookRouteLoader } from './books/book'

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
      },
      {
        path: '/books/:bookId',
        element: <Book />,
        loader: bookRouteLoader,
      },
      {
        path: '*',
        element: <div>404</div>,
      },
    ],
  },
])

export default function BM(_props: PreactIslandProps) {
  return <RouterProvider router={router} />
}
