import { useQuery } from '@tanstack/react-query'
import { useParams } from 'react-router'
import { httpManagerGetBook } from '../../api'
import { BookManagerContext } from './book-manager-context'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'

export default function BookManagerLayout({ children }: React.PropsWithChildren) {
  const { bookId } = useParams<{ bookId: string }>()

  const { data, status } = useQuery({
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
    <BookManagerContext.Provider value={{ book: data }}>
      <section className="bg-muted">
        <header className="page-header container-default">
          <Breadcrumb className="mb-2">
            <BreadcrumbList>
              <BreadcrumbItem>Manager</BreadcrumbItem>
              <BreadcrumbSeparator />
              <BreadcrumbItem>{data.name}</BreadcrumbItem>
            </BreadcrumbList>
          </Breadcrumb>

          <h1 className="page-header-text">{data.name}</h1>
        </header>
      </section>

      <main className="container-default pt-2">
        <BookManagerContext.Provider value={{ book: data }}>{children}</BookManagerContext.Provider>
      </main>
    </BookManagerContext.Provider>
  )
}
