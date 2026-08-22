import { useBEState } from './state'

export function ChapterNameInput() {
  const chapterName = useBEState((s) => s.chapterName)
  const valid = chapterName.trim().length > 0 && Array.from(chapterName.trim()).length <= 70

  return (
    <div className="my-4">
      <label className="text-secondary-foreground" htmlFor="chapter-name-input">
        {window._('editor.chapterName')}
      </label>
      <input
        id="chapter-name-input"
        name="chapterName"
        className="be-chapter-name-input"
        value={chapterName}
        aria-invalid={!valid}
        aria-describedby={!valid ? 'chapter-name-error' : undefined}
        onChange={(e) => {
          useBEState.getState().setChapterName((e.target as HTMLInputElement).value)
        }}
      />
      {!valid && (
        <p id="chapter-name-error" className="form-control__error" role="alert">
          <i className="fa-solid fa-circle-exclamation" aria-hidden="true" />
          {window._('editor.chapterNameInvalid')}
        </p>
      )}
    </div>
  )
}
