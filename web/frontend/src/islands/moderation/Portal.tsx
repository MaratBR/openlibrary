import { createHashRouter, Navigate, Outlet } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { ReactIslandProps } from '../common/react-island'
import ModerationLayout from './ModerationLayout'
import PlaceholderPage from './PlaceholderPage'

const routes = [
  ['overview', 'moderationPortal.overview'],
  ['books', 'moderationPortal.books'],
  ['search', 'moderationPortal.search'],
  ['chapters', 'moderationPortal.chapters'],
  ['comments', 'moderationPortal.comments'],
  ['users', 'moderationPortal.users'],
  ['login-history', 'moderationPortal.loginHistory'],
  ['audit-log', 'moderationPortal.auditLog'],
] as const

const router = createHashRouter([
  {
    path: '/',
    element: (
      <ModerationLayout>
        <Outlet />
      </ModerationLayout>
    ),
    children: [
      {
        path: '/',
        element: <Navigate to="/overview" replace />,
      },
      ...routes.map(([path, translationKey]) => ({
        path: `/${path}`,
        element: <PlaceholderPage title={window._(translationKey)} />,
      })),
      {
        path: '*',
        element: <div>404</div>,
      },
    ],
  },
])

export default function Portal(_props: ReactIslandProps) {
  return <RouterProvider router={router} />
}
