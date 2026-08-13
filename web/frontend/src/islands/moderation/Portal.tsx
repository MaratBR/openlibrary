import { createHashRouter, Navigate, Outlet } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { ReactIslandProps } from '../common/react-island'
import { useMemo } from 'react'
import ModerationLayout from './ModerationLayout'
import PlaceholderPage from './PlaceholderPage'
import UserModerationPage, {
  UserModerationErrorPage,
  userModerationRouteLoader,
} from './UserModerationPage'
import ModerationUsers, { ModerationUsersErrorPage, moderationUsersLoader } from './ModerationUsers'
import ReportPage, { ReportErrorPage, reportRouteLoader } from './ReportPage'
import ReportsPage, { ReportsErrorPage, reportsSearchLoader } from './ReportsPage'
import BookModerationPage, {
  BookModerationErrorPage,
  bookModerationRouteLoader,
} from './BookModerationPage'
import AuditLogPage, { AuditLogErrorPage, auditLogLoader } from './AuditLogPage'

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

function createModerationRouter(roles: string[]) {
  return createHashRouter([
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
        ...routes
          .filter(([path]) => path !== 'users' && path !== 'audit-log')
          .map(([path, translationKey]) => ({
            path: `/${path}`,
            element: <PlaceholderPage title={window._(translationKey)} />,
          })),
        {
          path: '/users',
          element: <ModerationUsers roles={roles} />,
          errorElement: <ModerationUsersErrorPage />,
          loader: moderationUsersLoader,
        },
        {
          path: '/audit-log',
          element: <AuditLogPage />,
          errorElement: <AuditLogErrorPage />,
          loader: auditLogLoader,
        },
        {
          path: '/users/:userId',
          element: <UserModerationPage />,
          errorElement: <UserModerationErrorPage />,
          loader: userModerationRouteLoader,
        },
        {
          path: '/reports',
          element: <ReportsPage view="overview" />,
        },
        {
          path: '/reports/search',
          element: <ReportsPage view="search" />,
          errorElement: <ReportsErrorPage />,
          loader: reportsSearchLoader,
        },
        {
          path: '/reports/:reportId',
          element: <ReportPage />,
          errorElement: <ReportErrorPage />,
          loader: reportRouteLoader,
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
          path: '/books/:bookId',
          element: <BookModerationPage />,
          errorElement: <BookModerationErrorPage />,
          loader: bookModerationRouteLoader,
        },
        ...['actions', 'chapters', 'activity'].map((resource) => ({
          path: `/books/:bookId/${resource}`,
          element: <BookModerationPage />,
          errorElement: <BookModerationErrorPage />,
          loader: bookModerationRouteLoader,
        })),
        // TODO: Replace this placeholder when comment moderation is implemented.
        {
          path: '/comments/:commentId',
          element: <PlaceholderPage title={window._('moderationPortal.comments')} />,
        },
        {
          path: '*',
          element: <div>404</div>,
        },
      ],
    },
  ])
}

type ModerationPortalData = { roles: string[] }

export default function Portal({ data }: ReactIslandProps) {
  const roles = (data as ModerationPortalData | undefined)?.roles ?? []
  const router = useMemo(() => createModerationRouter(roles), [roles])
  return <RouterProvider router={router} />
}
