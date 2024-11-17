import { createBrowserRouter, Outlet } from 'react-router-dom'
import LoginPage from './modules/auth/views/LoginPage'
import HomePage from './modules/book/views/HomePage'
import SiteHeader from './components/site-header'
import BookPage from './modules/book/views/BookPage'
import ChapterPage from './modules/book/views/ChapterPage'
import { Suspense } from 'react'
import { componentsChunk } from './lib/utils'
import Spinner from './components/spinner'
import { NotificationRenderer } from './modules/notifications'
import SearchPage from './modules/book/views/SearchPage/SearchPage'
import UserProfile, { UserProfileInner } from './modules/user/views/UserProfile'

const bookManagerChunk = componentsChunk(() => import('./modules/book-manager'))
const BookManagerLayout = bookManagerChunk.componentType('BookManagerLayout')

const router = createBrowserRouter([
  {
    path: 'login',
    element: <LoginPage />,
  },
  {
    path: 'user/:userId/__profile',
    element: <UserProfileInner />,
  },
  {
    path: '*',
    element: (
      <>
        <SiteHeader />
        <div style={{ marginTop: 'var(--site-header-height)' }}>
          <NotificationRenderer />
          <Outlet />
        </div>
      </>
    ),
    children: [
      {
        path: 'home',
        element: <HomePage />,
      },

      {
        path: 'search',
        element: <SearchPage />,
      },
      {
        path: 'book/:id',
        element: <BookPage />,
      },
      {
        path: 'book/:bookId/chapters/:chapterId',
        element: <ChapterPage />,
      },
      {
        path: 'new-book',
        element: bookManagerChunk.element('NewBook'),
      },
      {
        path: 'new-book/import-from-ao3',
        element: bookManagerChunk.element('ImportBookFromAo3'),
      },
      {
        path: 'user/:userId',
        element: <UserProfile />,
      },
      {
        path: 'manager/books',
        element: bookManagerChunk.element('MyBooks'),
      },
      {
        path: 'manager/book/:bookId',
        element: (
          <BookManagerLayout>
            <Suspense fallback={<Spinner />}>
              <Outlet />
            </Suspense>
          </BookManagerLayout>
        ),

        children: [
          {
            path: '',
            element: bookManagerChunk.element('BookManagerHomePage'),
          },
          {
            path: 'new-chapter',
            element: bookManagerChunk.element('CreateChapterPage'),
          },
          {
            path: 'chapters/:chapterId',
            element: bookManagerChunk.element('EditChapterPage'),
          },
          {
            path: 'reorder-chapters',
            element: bookManagerChunk.element('BookChaptersReorder'),
          },
        ],
      },
    ],
  },
])

export default router
