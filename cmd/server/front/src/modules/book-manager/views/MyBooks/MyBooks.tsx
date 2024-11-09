import { useQuery } from '@tanstack/react-query'
import Spinner from '@/components/spinner'
import { Button } from '@/components/ui/button'
import { PenIcon } from 'lucide-react'
import { useNavigate } from 'react-router'
import { httpGetMyBooks } from '../../api'
import AuthorBookCard from './AuthorBookCard'

export default function MyBooks() {
  return (
    <main className="container-default">
      <header className="page-header">
        <h1 className="page-header-text">My books</h1>
      </header>
      <BooksList />
    </main>
  )
}

function BooksList() {
  const navigate = useNavigate()

  const { data, isFetching } = useQuery({
    queryKey: ['my-books'],
    queryFn: () => httpGetMyBooks().then((r) => r.books),
    initialData: [],
  })

  return (
    <div>
      {isFetching && <Spinner />}
      {data.length === 0 && !isFetching && (
        <div className="mb-5">
          <p className="text-muted-foreground">You have not written any books yet.</p>
        </div>
      )}
      <Button variant="outline" onClick={handleStartNewBookClick}>
        <PenIcon />
        Start a new book
      </Button>
      <div className="space-y-2 mt-6">
        {data.map((book) => (
          <AuthorBookCard key={book.id} book={book} />
        ))}
      </div>
    </div>
  )

  function handleStartNewBookClick() {
    navigate('/new-book')
  }
}
