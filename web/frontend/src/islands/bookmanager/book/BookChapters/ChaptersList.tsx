import { ManagerBookChapterDto } from '@/api/bm'
import SanitizeHTML from '@/common/SanitizeHTML'
import { useWindowVirtualizer } from '@tanstack/react-virtual'
import { useRef } from 'preact/hooks'
import ReorderChapterButton from './ReorderChapterButton'
import { useBookChaptersState } from './state'
import clsx from 'clsx'

export default function ChaptersList() {
  const listRef = useRef<HTMLDivElement | null>(null)

  const { chapters, bookId } = useBookChaptersState()

  const virtualizer = useWindowVirtualizer({
    count: chapters.length,
    estimateSize: () => 120,
    overscan: 5,
    scrollMargin: listRef.current?.offsetTop ?? 0,
  })

  return (
    <div ref={listRef} class="bm-chapters-list">
      <div
        style={{
          height: `${virtualizer.getTotalSize()}px`,
          width: '100%',
          position: 'relative',
        }}
      >
        {virtualizer.getVirtualItems().map((item) => (
          <div
            key={item.key}
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '100%',
              height: `${item.size}px`,
              transform: `translateY(${item.start - virtualizer.options.scrollMargin}px)`,
            }}
          >
            <ChapterCard
              isLast={item.index === chapters.length - 1}
              chapter={chapters[item.index]}
              bookId={bookId}
              height={120}
            />
          </div>
        ))}
      </div>
    </div>
  )
}

function ChapterCard({
  chapter,
  bookId,
  height,
  isLast,
}: {
  chapter: ManagerBookChapterDto
  bookId: string
  height: number
  isLast: boolean
}) {
  return (
    <div
      class={clsx('bm-chapters-list-item', {
        'bm-chapters-list-item--last': isLast,
      })}
      style={{ height }}
    >
      <div class="bm-chapters-list-item__head">
        <span>{chapter.name}</span>
      </div>

      {chapter.summary && (
        <div class="bm-chapters-list-item__summary">
          <p class="user-content mb-2">
            <SanitizeHTML value={chapter.summary} />
          </p>
        </div>
      )}

      <div class="btn-group btn-group--rounded-md border bm-chapters-list-item__actions">
        <a class="btn btn--ghost" href={`/books-manager/book/${bookId}/chapter/${chapter.id}`}>
          <i class="fa-solid fa-pen" />
          &nbsp;
          {window._('common.edit')}
        </a>
        <a
          class="btn btn--ghost destructive"
          href={`/books-manager/book/${bookId}/chapter/${chapter.id}`}
        >
          <i class="fa-solid fa-trash" />
          &nbsp;
          {window._('common.delete')}
        </a>
        <ReorderChapterButton chapter={chapter} />
      </div>
    </div>
  )
}
