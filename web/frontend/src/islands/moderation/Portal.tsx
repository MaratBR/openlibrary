import { createHashRouter, Navigate, Outlet } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { ReactIslandProps } from '../common/react-island'
import ModerationLayout from './ModerationLayout'
import PlaceholderPage from './PlaceholderPage'
import UserModerationPage, {
  UserModerationErrorPage,
  userModerationRouteLoader,
} from './UserModerationPage'

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
        path: '/users/:userId',
        element: <UserModerationPage />,
        errorElement: <UserModerationErrorPage />,
        loader: userModerationRouteLoader,
      },
      {
        path: '/users/:userId/activity',
        element: <UserModerationPage />,
        errorElement: <UserModerationErrorPage />,
        loader: userModerationRouteLoader,
      },
      {
        path: '/users/:userId/actions',
        element: <UserModerationPage />,
        errorElement: <UserModerationErrorPage />,
        loader: userModerationRouteLoader,
      },
      ...['history', 'reports', 'login-history', 'books', 'comments'].map((resource) => ({
        path: `/users/:userId/${resource}`,
        element: <UserModerationPage />,
        errorElement: <UserModerationErrorPage />,
        loader: userModerationRouteLoader,
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
