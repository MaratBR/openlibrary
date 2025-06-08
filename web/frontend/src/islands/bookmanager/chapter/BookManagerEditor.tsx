import BookContentEditor from './BookContentEditor'
import { useMemo } from 'preact/hooks'
import { PreactIslandProps } from '@/islands/common'
import { DraftDtoSchema } from '../contracts'

const dataSchema = DraftDtoSchema

export default function BookManagerEditor({ data }: PreactIslandProps) {
  const draft = useMemo(() => dataSchema.parse(data), [data])

  return <BookContentEditor content={draft.content} />
}
