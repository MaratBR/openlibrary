import { PreactIslandProps } from '../common/preact-island'
import { createHashRouter, Outlet } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import BMSidebar from './BMSidebar'

const router = createHashRouter([
  {
    path: '/',
    element: (
      <BMSidebar>
        <Outlet />
      </BMSidebar>
    ),
    children: [{ path: '*', element: <div>404</div> }],
  },
])

export default function BM(_props: PreactIslandProps) {
  return <RouterProvider router={router} />
}
