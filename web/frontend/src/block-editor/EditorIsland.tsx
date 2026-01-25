import { useLayoutEffect, useMemo } from 'preact/hooks'
import { PreactIslandProps } from '@/islands/common/preact-island'
import { z } from 'zod'
import './BookManagerEditor.scss'
import { DraftDtoSchema } from './contracts'
import { useBEState } from './state'
import { EditorIframe } from './EditorIframe'
import { SaveButton } from './SaveButton'
import { CenterHeader } from './CenterHeader'

const dataSchema = z.object({
  bookId: z.string(),
  draft: DraftDtoSchema,
})

export default function EditorIslandComponent({ data }: PreactIslandProps) {
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
          <div class="be-header__center">
            <CenterHeader draft={draft} />
          </div>
          <div class="be-header__right">
            <SaveButton />
          </div>
          <div />
        </header>
      </div>
      <div class="be-layout__body">
        <div class="be-layout__left">Left</div>
        <div class="be-layout__center">
          <EditorIframe />
        </div>
        <div class="be-layout__right">Right</div>
      </div>
    </div>
  )
}
