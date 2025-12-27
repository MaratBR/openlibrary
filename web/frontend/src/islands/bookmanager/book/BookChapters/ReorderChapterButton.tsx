import { ManagerBookChapterDto } from '@/api/bm'
import Popper from '@/components/Popper'
import { useRef, useState } from 'preact/hooks'
import { useChapterSelectorPopup } from './ChapterSelectorPopup'
import { useBookChaptersState } from './state'
import { ChapterSelectorDescriptionProps } from '../ChapterSelector'

export default function ReorderChapterButton({
  chapter,
  bookId,
}: {
  bookId: string
  chapter: ManagerBookChapterDto
}) {
  const ref = useRef<HTMLButtonElement | null>(null)
  const [reorderActionsOpen, setReorderActionsOpen] = useState(false)
  const chapterSelectorPopup = useChapterSelectorPopup()

  function swapChapter() {
    useBookChaptersState.getState().setReorderActiveChapter(chapter)

    const { current: element } = ref
    if (!element) return
    chapterSelectorPopup.open({
      element,
      onSelected: async (swapWith) => {
        await useBookChaptersState.getState().swapChapters(chapter.id, swapWith.id)
        chapterSelectorPopup.close()
      },
      ActionDescriptionComponent: ActionDescriptionSwap,
    })
  }

  function insertAfter() {
    useBookChaptersState.getState().setReorderActiveChapter(chapter)

    const { current: element } = ref
    if (!element) return
    chapterSelectorPopup.open({
      element,
      onSelected: async (insertAfter) => {
        await useBookChaptersState.getState().insertAfter(chapter.id, insertAfter.id)
        chapterSelectorPopup.close()
      },
      ActionDescriptionComponent: ActionDescriptionInsertAfter,
    })
  }

  function insertBefore() {
    useBookChaptersState.getState().setReorderActiveChapter(chapter)

    const { current: element } = ref
    if (!element) return
    chapterSelectorPopup.open({
      element,
      onSelected: async (insertBefore) => {
        await useBookChaptersState.getState().insertBefore(chapter.id, insertBefore.id)
        chapterSelectorPopup.close()
      },
      ActionDescriptionComponent: ActionDescriptionInsertBefore,
    })
  }

  return (
    <>
      <button ref={ref} class="btn btn--ghost" onClick={() => setReorderActionsOpen(true)}>
        <i class="fa-solid fa-up-down" />
        {window._('bookManager.edit.reorder')}
      </button>

      <Popper
        placement="left"
        anchorEl={ref}
        open={reorderActionsOpen}
        onClose={() => setReorderActionsOpen(false)}
      >
        <div class="card px-0">
          <ul class="btn-list">
            <li>
              <button class="btn btn--ghost" onClick={() => swapChapter()}>
                <i class="fa-solid fa-retweet" /> Swap with another
              </button>
            </li>
            <li>
              <button class="btn btn--ghost" onClick={() => insertBefore()}>
                <i class="fa-solid fa-diagram-predecessor" /> Insert before...
              </button>
            </li>
            <li>
              <button class="btn btn--ghost" onClick={() => insertAfter()}>
                <i class="fa-solid fa-diagram-predecessor" /> Insert after...
              </button>
            </li>
          </ul>
        </div>
      </Popper>
    </>
  )
}

function ActionDescriptionSwap({ selectedChapter }: ChapterSelectorDescriptionProps) {
  return (
    <ActionDescription
      translationKey="bookManager.edit.swapChapterWith"
      chapter={selectedChapter}
    />
  )
}

function ActionDescriptionInsertAfter({ selectedChapter }: ChapterSelectorDescriptionProps) {
  return (
    <ActionDescription
      translationKey="bookManager.edit.insertChapterAfter"
      chapter={selectedChapter}
    />
  )
}

function ActionDescriptionInsertBefore({ selectedChapter }: ChapterSelectorDescriptionProps) {
  return (
    <ActionDescription
      translationKey="bookManager.edit.insertChapterBefore"
      chapter={selectedChapter}
    />
  )
}

function ActionDescription({
  chapter,
  translationKey,
}: {
  chapter: ManagerBookChapterDto | null
  translationKey: string
}) {
  const activeChapter = useBookChaptersState((s) => s.reorderActiveChapter)
  if (!activeChapter) return null

  return (
    <span>
      {window._(translationKey, {
        ChapterName: <strong>{activeChapter.name}</strong>,
        Another: chapter ? <strong>{chapter.name}</strong> : '...',
      })}
    </span>
  )
}
