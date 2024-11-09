import { Card } from '@/components/ui/card'
import { ManagerAuthorBookDto } from '../../api'
import { NavLink } from 'react-router-dom'

export type BookCardProps = {
  book: ManagerAuthorBookDto
}

export default function BookCard({ book }: BookCardProps) {
  return (
    <Card data-test-id="book-card" data-book-id={book.id} className="rounded-[5px] p-2">
      <NavLink to={`/book/${book.id}`}>{book.name}</NavLink>
    </Card>
  )
}
