import { useLayoutEffect, useMemo } from 'react'
import { ReactIslandProps } from '@/islands/common/react-island'
import { z } from 'zod'
import './BookManagerEditor.scss'
import { DraftDtoSchema } from './contracts'
import { useBEState } from './state'
import { EditorIframe } from './EditorIframe'
import { SaveButton } from './SaveButton'
import { CenterHeader } from './CenterHeader'
import { WidgetsMenu, WidgetsService } from './widgets'

const dataSchema = z.object({
  bookId: z.string(),
  draft: DraftDtoSchema,
})

export default function EditorIslandComponent({ data }: ReactIslandProps) {
  const { draft } = useMemo(() => dataSchema.parse(data), [data])

  useLayoutEffect(() => {
    useBEState.getState().init(draft)
  }, [draft])

  return (
    <div className="be-layout">
      <div className="be-layout__header">
        <header className="be-header">
          <div />
          <div className="be-header__left">Left</div>
          <div className="be-header__center">
            <CenterHeader draft={draft} />
          </div>
          <div className="be-header__right">
            <SaveButton />
          </div>
          <div />
        </header>
      </div>
      <div className="be-layout__body">
        <div className="be-layout__left">
          <WidgetsMenu service={WidgetsService.instance()} />
        </div>
        <div className="be-layout__center">
          <EditorIframe initialContent={draft.content} />
        </div>
        <div className="be-layout__right">Right</div>
      </div>
    </div>
  )
}
