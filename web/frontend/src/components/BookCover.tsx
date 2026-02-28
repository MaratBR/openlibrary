import { BookCover as BookCoverDto } from '@/api/common'

export function BookCover({ cover }: { cover: BookCoverDto }) {
  return (
    <div class="book-cover">
      <img src={cover.url} />
    </div>
  )
}
