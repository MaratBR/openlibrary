import { useBEState } from './state'

export function ChapterNameInput() {
  const chapterName = useBEState((s) => s.chapterName)

  return (
    <div className="my-4">
      <span className="text-muted-foreground">Chapter name</span>
      <input
        name="chapterName"
        className="be-chapter-name-input"
        value={chapterName}
        onChange={(e) => {
          useBEState.getState().setChapterName((e.target as HTMLInputElement).value)
        }}
      />
    </div>
  )
}
