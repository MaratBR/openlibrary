import { useAtom, useAtomValue } from 'jotai'
import { chapterNameAtom, chapterNameIsValidAtom, maxChapterNameLengthAtom } from './state'

export function ChapterNameInput() {
  const [chapterName, setChapterName] = useAtom(chapterNameAtom)
  const valid = useAtomValue(chapterNameIsValidAtom)
  const maxChapterNameLength = useAtomValue(maxChapterNameLengthAtom)

  return (
    <div className="be-chapter-name">
      <label className="sr-only" htmlFor="chapter-name-input">
        {window._('editor.chapterName')}
      </label>
      <textarea
        id="chapter-name-input"
        name="chapterName"
        className="be-chapter-name-input"
        rows={1}
        value={chapterName}
        aria-invalid={!valid}
        aria-describedby={!valid ? 'chapter-name-error' : undefined}
        onKeyDown={(e) => {
          if (e.key === 'Enter') e.preventDefault()
        }}
        onChange={(e) => {
          setChapterName(e.currentTarget.value.replace(/[\r\n]+/g, ' '))
        }}
      />
      {!valid && (
        <p id="chapter-name-error" className="FormControl-error" role="Alert">
          <i className="fa-solid fa-circle-exclamation" aria-hidden="true" />
          {window._('editor.chapterNameInvalid', { max: String(maxChapterNameLength) })}
        </p>
      )}
    </div>
  )
}
