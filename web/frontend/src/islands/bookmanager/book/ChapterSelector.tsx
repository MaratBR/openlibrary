import { ManagerBookChapterDto } from '@/api/bm'
import { useMemo, useRef, useState } from 'preact/hooks'
import { ComponentType, CSSProperties, MouseEventHandler } from 'preact'
import './ChapterSelector.scss'
import { useVirtualizer } from '@tanstack/react-virtual'
import clsx from 'clsx'
import { textSearch } from '@/lib/search'

export type ChapterSelectorDescriptionProps = {
  selectedChapter: ManagerBookChapterDto | null
}

export type ChapterSelectorProps = {
  chapters: ManagerBookChapterDto[]
  onSelected: (chapter: ManagerBookChapterDto) => void | Promise<void>
  ActionDescriptionComponent?: ComponentType<ChapterSelectorDescriptionProps>
  shouldDisableItem?: (chapter: ManagerBookChapterDto) => boolean
}

export default function ChapterSelector({
  chapters,
  onSelected,
  ActionDescriptionComponent,
  shouldDisableItem,
}: ChapterSelectorProps) {
  const [search, setSearch] = useState('')
  const [selectedChapter, setSelectedChapter] = useState<ManagerBookChapterDto | null>(null)
  const [loading, setLoading] = useState(false)

  const listRef = useRef<HTMLDivElement | null>(null)

  const filteredChapters = useMemo(() => {
    if (search.trim().length === 0) return chapters

    return chapters.filter((c) => textSearch(search, c.name))
  }, [chapters, search])

  const virtualizer = useVirtualizer({
    count: filteredChapters.length,
    estimateSize: () => 48,
    overscan: 5,
    getScrollElement: () => listRef.current || null,
  })

  function handleConfirm() {
    if (!selectedChapter) return
    const result = onSelected(selectedChapter)
    if (result instanceof Promise) {
      setLoading(true)
      result.finally(() => {
        setLoading(false)
      })
    }
  }

  return (
    <div class="ChapterSelector">
      {ActionDescriptionComponent && (
        <div class="ChapterSelector__actionDescription">
          <ActionDescriptionComponent selectedChapter={selectedChapter} />
        </div>
      )}
      <div class="ChapterSelector__inputContainer">
        <input
          class="ChapterSelector__input"
          value={search}
          placeholder="Search..."
          onChange={(e) => setSearch((e.target as HTMLInputElement).value)}
        />
      </div>
      <div ref={listRef} class="ChapterSelector__listContainer">
        <div
          role="list"
          style={{
            height: `${virtualizer.getTotalSize()}px`,
            width: '100%',
            position: 'relative',
          }}
        >
          {virtualizer.getVirtualItems().map((item) => (
            <ChapterItem
              key={filteredChapters[item.index].id}
              chapter={filteredChapters[item.index]}
              disabled={shouldDisableItem ? shouldDisableItem(filteredChapters[item.index]) : false}
              onClick={() => setSelectedChapter(filteredChapters[item.index])}
              isSelected={filteredChapters[item.index].id === selectedChapter?.id}
              style={{
                position: 'absolute',
                top: 0,
                left: 0,
                width: '100%',
                height: `${item.size}px`,
                transform: `translateY(${item.start - virtualizer.options.scrollMargin}px)`,
              }}
            />
          ))}
        </div>
      </div>
      <div class="ChapterSelector__confirm">
        <button disabled={!selectedChapter || loading} class="btn w-32" onClick={handleConfirm}>
          {loading ? <span class="loader loader--dark" /> : window._('common.confirm')}
        </button>
      </div>
    </div>
  )
}

function ChapterItem({
  chapter,
  style,
  isSelected,
  onClick,
  disabled,
}: {
  chapter: ManagerBookChapterDto
  isSelected: boolean
  style: CSSProperties
  onClick: MouseEventHandler<HTMLButtonElement>
  disabled: boolean
}) {
  return (
    <button
      disabled={disabled}
      aria-disabled={disabled ? 'true' : 'false'}
      onClick={onClick}
      class={clsx('listitem ChapterSelector__item', {
        'ChapterSelector__item--selected': isSelected,
      })}
      role="listitem"
      style={style}
    >
      {isSelected && (
        <>
          <i class="fa-solid fa-circle-check" />
          &nbsp;
        </>
      )}
      {chapter.name}
    </button>
  )
}
