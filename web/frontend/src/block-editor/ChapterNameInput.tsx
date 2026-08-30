import { useAtom, useAtomValue } from 'jotai'
import { chapterNameAtom, chapterNameIsValidAtom } from './state'

export function ChapterNameInput() {
  const [chapterName, setChapterName] = useAtom(chapterNameAtom)
  const valid = useAtomValue(chapterNameIsValidAtom)

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
          setChapterName(e.currentTarget.value)
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
