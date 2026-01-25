import { DraftDto } from './contracts'

export function CenterHeader({ draft }: { draft: DraftDto }) {
  return (
    <div class="flex justify-center items-center h-full">
      <div class="flex justify-center items-center bg-secondary p-1 gap-1 rounded-xl">
        <a
          href={`/book/${draft.book.id}/chapters/${draft.chapterId}`}
          target="_blank"
          class="btn btn--icon btn--ghost"
          rel="noreferrer"
        >
          <i class="fa-solid fa-up-right-from-square" />
        </a>
      </div>
    </div>
  )
}
