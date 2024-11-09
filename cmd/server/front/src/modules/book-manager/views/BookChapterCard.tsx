import { useChapterName } from '@/modules/book/utils'
import { ManagerBookChapterDto } from '../api'

import './BookChapterCard.css'
import { NavLink } from 'react-router-dom'
import { useBookManager } from './book-manager-context'
import Timestamp from '@/components/timestamp'

export default function BookChapterCard({ chapter }: { chapter: ManagerBookChapterDto }) {
  const chapterName = useChapterName(chapter.name, chapter.order)
  const { book } = useBookManager()

  return (
    <div className="manager-chapter-card">
      <NavLink
        to={`/manager/book/${book.id}/chapters/${chapter.id}`}
        className="link-default manager-chapter-card__name"
      >
        {chapterName}
      </NavLink>
      <span> &bull; {chapter.words} words</span>

      <p className="manager-chapter-card__summary">{chapter.summary}</p>

      <p className="manager-chapter-card__info">
        <span className="font-semibold">Created:</span> <Timestamp value={chapter.createdAt} />
      </p>
    </div>
  )
}
