import { NavLink } from 'react-router-dom'

export default function GoToBookPage({ bookId }: { bookId: string }) {
  return (
    <NavLink to={`/book/${bookId}`} className="badge-alt">
      Go to book's page
    </NavLink>
  )
}
