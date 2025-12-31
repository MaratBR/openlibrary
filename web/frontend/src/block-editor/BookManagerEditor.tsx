import { useLayoutEffect, useMemo, useState } from 'preact/hooks'
import { PreactIslandProps } from '@/islands/common/preact-island'
import { z } from 'zod'
import './BookManagerEditor.scss'
import { DraftDtoSchema } from './contracts'
import { EditorIframe } from './wysiwyg'
import { useWYSIWYG, useWYSIWYGHasChanges } from './wysiwyg/state'
import { useBEState } from './state'
import { createPortal } from 'preact/compat'

const dataSchema = z.object({
  bookId: z.string(),
  draft: DraftDtoSchema,
})

export default function BookManagerEditor({ data }: PreactIslandProps) {
  const { draft } = useMemo(() => dataSchema.parse(data), [data])

  useLayoutEffect(() => {
    useBEState.getState().init(draft)
  }, [draft])

  return (
    <div class="be-layout">
      <div class="be-layout__header">
        <header class="be-header">
          <div />
          <div class="be-header__left">Left</div>
          <div class="be-header__center">Center</div>
          <div class="be-header__right">
            <SaveButton />
          </div>
          <div />
        </header>
      </div>
      <div class="be-layout__body">
        <div class="be-layout__left">Left</div>
        <div class="be-layout__center">
          <ChapterNameInput />
          <EditorIframe />
        </div>
        <div class="be-layout__right">Right</div>
      </div>
    </div>
  )
}

function SaveButton() {
  const wasChangedFirstTime = useWYSIWYGHasChanges()
  const saving = useBEState((s) => s.saving)

  function handleClick() {
    useBEState.getState().saveDraft()
  }

  return (
    <button
      onClick={handleClick}
      disabled={!wasChangedFirstTime || saving}
      class="btn btn--lg btn--sq w-30 flex justify-center items-center"
    >
      {saving ? <span class="loader loader--dark" /> : window._('common.save')}
    </button>
  )
}

function ChapterNameInput() {
  const chapterName = useBEState((s) => s.chapterName)
  const elements = useWYSIWYG((s) => s.initData?.elements)
  const [container, setContainer] = useState<HTMLElement | null>(null)

  useLayoutEffect(() => {
    if (!elements) return
    const container = document.createElement('div')
    container.classList.add('contents')

    elements.contentWrapper.prepend(container)
    setContainer(container)
    return () => {
      setContainer(null)
    }
  }, [elements])

  if (!elements || !container) return null

  return createPortal(
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
    </div>,
    container,
  )
}
