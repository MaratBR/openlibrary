import { useLayoutEffect, useMemo } from 'preact/hooks'
import { PreactIslandProps } from '@/islands/common/preact-island'
import { z } from 'zod'
import './BookManagerEditor.scss'
import { DraftDtoSchema } from './contracts'
import { EditorIframe, useWYSIWYG } from './wysiwyg'

const dataSchema = z.object({
  bookId: z.string(),
  draft: DraftDtoSchema,
})

export default function BookManagerEditor({ data }: PreactIslandProps) {
  const { draft } = useMemo(() => dataSchema.parse(data), [data])

  useLayoutEffect(() => {
    useWYSIWYG.getState().editor.setContentAndClearHistory(draft.content)
  }, [draft])

  return (
    <div class="be-layout">
      <div class="be-layout__header">
        <header class="be-header">TODO: Editor header content</header>
      </div>
      <div class="be-layout__body">
        <EditorIframe />
      </div>
    </div>
  )
}
