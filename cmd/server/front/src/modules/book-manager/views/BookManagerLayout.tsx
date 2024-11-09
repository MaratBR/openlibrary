import { useQuery } from '@tanstack/react-query'
import { useParams } from 'react-router'
import { httpManagerGetBook } from '../api'
import { BookManagerContext } from './book-manager-context'
import { NavLink } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'
import BookIsHiddenIndicator from './common/book-is-hidden-indicator'
import BookIsBannedIndicator from './common/book-is-banned-indicator'
import GoToBookPage from './common/go-to-book-page'

export default function BookManagerLayout({ children }: React.PropsWithChildren) {
  const { bookId } = useParams<{ bookId: string }>()

  const { data, status, refetch } = useQuery({
    queryKey: ['manager', 'book', bookId],
    enabled: !!bookId,
    queryFn: () => httpManagerGetBook(bookId!),
    refetchInterval: false,
    retry: false,
  })

  if (status === 'error') {
    return (
      <div className="container-default pt-10">
        <h2 className="text-xl">404: Could not find the book you requested</h2>
      </div>
    )
  }

  if (!data) return null

  return (
    <>
      <section className="bg-muted">
        <header className="page-header container-default">
          <NavLink to="/manager/books" className="link-default inline-flex gap-2 p-1">
            <ArrowLeft /> Back to list of your books
          </NavLink>
          <h1 className="page-header-text">{data.name}</h1>

          <div className="flex gap-2 mt-5">
            <GoToBookPage bookId={data.id} />
            {!data.isPubliclyVisible && <BookIsHiddenIndicator />}
            {data.isBanned && <BookIsBannedIndicator bookId={data.id} />}
          </div>
        </header>
      </section>

      <main className="container-default pt-2">
        <BookManagerContext.Provider value={{ book: data, refetch }}>
          {children}
        </BookManagerContext.Provider>
      </main>
    </>
  )
}
