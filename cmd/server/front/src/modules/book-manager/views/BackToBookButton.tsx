import { ArrowLeftIcon } from 'lucide-react'
import { useBookManager } from './book-manager-context'
import { NavLink } from 'react-router-dom'

export default function BackToBookButton() {
  const { book } = useBookManager()

  return (
    <NavLink className="link-default inline-flex gap-1 mb-5" to={`/manager/book/${book.id}`}>
      <ArrowLeftIcon />
      Back to {book.name}
    </NavLink>
  )
}
