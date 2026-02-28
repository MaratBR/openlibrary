import { PreactIslandProps } from '../common/preact-island'
import { createHashRouter, Outlet } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import BMLayout from './BMLayout'
import { Books, booksRouteLoader } from './books'

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
        path: '/books',
        element: <Books />,
        loader: booksRouteLoader,
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
