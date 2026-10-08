import { dashboardItem } from '@/components/dashboard-nav-item'
import { createHashRouter, Navigate, Outlet, redirect } from 'react-router'
import { RouterProvider } from 'react-router/dom'
import { ReactIslandProps } from '@/islands/common/react-island'
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
import ModerationBooks, { ModerationBooksErrorPage, moderationBooksLoader } from './ModerationBooks'
import LoginHistoryPage, { LoginHistoryErrorPage, loginHistoryLoader } from './LoginHistoryPage'

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

function createModerationRouter(roles: string[], isAdmin: boolean) {
  return createHashRouter([
    {
      path: '/',
      element: (
        <ModerationLayout isAdmin={isAdmin}>
          <Outlet />
        </ModerationLayout>
      ),
      children: [
        {
          path: '/',
          element: <Navigate to="/overview" replace />,
        },
        ...routes
          .filter(
            ([path]) =>
              path !== 'users' &&
              path !== 'books' &&
              path !== 'login-history' &&
              path !== 'audit-log',
          )
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
          path: '/books',
          element: <ModerationBooks />,
          errorElement: <ModerationBooksErrorPage />,
          loader: moderationBooksLoader,
        },
        {
          path: '/login-history',
          element: <LoginHistoryPage />,
          errorElement: <LoginHistoryErrorPage />,
          loader: loginHistoryLoader,
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
          handle: dashboardItem<Awaited<ReturnType<typeof userModerationRouteLoader>>>(
            (data) => data.user.name,
          ),
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
          handle: dashboardItem<Awaited<ReturnType<typeof reportRouteLoader>>>(
            (data) => `#${data.number} · ${data.reason}`,
          ),
        },
        {
          path: '/users/:userId/activity',
          element: <UserModerationPage />,
          errorElement: <UserModerationErrorPage />,
          loader: userModerationRouteLoader,
          handle: dashboardItem<Awaited<ReturnType<typeof userModerationRouteLoader>>>(
            (data) => data.user.name,
          ),
        },
        {
          path: '/users/:userId/actions',
          element: <UserModerationPage />,
          errorElement: <UserModerationErrorPage />,
          loader: userModerationRouteLoader,
          handle: dashboardItem<Awaited<ReturnType<typeof userModerationRouteLoader>>>(
            (data) => data.user.name,
          ),
        },
        ...['history', 'reports', 'books', 'comments'].map((resource) => ({
          path: `/users/:userId/${resource}`,
          element: <UserModerationPage />,
          errorElement: <UserModerationErrorPage />,
          loader: userModerationRouteLoader,
          handle: dashboardItem<Awaited<ReturnType<typeof userModerationRouteLoader>>>(
            (data) => data.user.name,
          ),
        })),
        {
          path: '/users/:userId/login-history',
          loader: ({ params }) =>
            redirect(`/login-history?users=${encodeURIComponent(params.userId ?? '')}`),
          element: null,
        },
        {
          path: '/books/:bookId',
          element: <BookModerationPage />,
          errorElement: <BookModerationErrorPage />,
          loader: bookModerationRouteLoader,
          handle: dashboardItem<Awaited<ReturnType<typeof bookModerationRouteLoader>>>(
            (data) => data.book.name,
          ),
        },
        ...['actions', 'chapters', 'activity'].map((resource) => ({
          path: `/books/:bookId/${resource}`,
          element: <BookModerationPage />,
          errorElement: <BookModerationErrorPage />,
          loader: bookModerationRouteLoader,
          handle: dashboardItem<Awaited<ReturnType<typeof bookModerationRouteLoader>>>(
            (data) => data.book.name,
          ),
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

type ModerationPortalData = { roles: string[]; isAdmin: boolean }

export default function Portal({ data }: ReactIslandProps) {
  const roles = (data as ModerationPortalData | undefined)?.roles ?? []
  const isAdmin = (data as ModerationPortalData | undefined)?.isAdmin ?? false
  const router = useMemo(() => createModerationRouter(roles, isAdmin), [roles, isAdmin])
  return <RouterProvider router={router} />
}
