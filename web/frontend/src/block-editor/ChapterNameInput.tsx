import { useBEState } from './state'

export function ChapterNameInput() {
  const chapterName = useBEState((s) => s.chapterName)

  return (
    <div class="my-4">
      <span class="text-muted-foreground">Chapter name</span>
      <input
        name="chapterName"
        class="be-chapter-name-input"
        value={chapterName}
        onChange={(e) => {
          useBEState.getState().setChapterName((e.target as HTMLInputElement).value)
        }}
      />
    </div>
  )
}
