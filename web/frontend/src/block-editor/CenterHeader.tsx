import { DraftDto } from './contracts'

export function CenterHeader({ draft }: { draft: DraftDto }) {
  return (
    <div className="flex justify-center items-center h-full">
      <div className="flex justify-center items-center bg-surface p-1 gap-1 rounded-xl">
        <a
          href={`/book/${draft.book.id}/chapters/${draft.chapter.id}`}
          target="_blank"
          className="btn btn--icon btn--primary"
          rel="noreferrer"
        >
          <i className="fa-solid fa-up-right-from-square" />
        </a>
      </div>
    </div>
  )
}
