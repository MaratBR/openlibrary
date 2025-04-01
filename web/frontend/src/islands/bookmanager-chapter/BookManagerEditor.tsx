import BookContentEditor from './BookContentEditor'
import { PreactIslandProps } from '../common'
import { z } from 'zod'
import { useMemo } from 'preact/hooks'

const managerBookChapterDetailsDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  createdAt: z.string(),
  words: z.number(),
  summary: z.string(),
  order: z.number().int(),
  isAdultOverride: z.boolean(),
  content: z.string(),
  isPubliclyVisible: z.boolean(),
})

export default function BookManagerEditor({ data }: PreactIslandProps) {
  const chapter = useMemo(() => managerBookChapterDetailsDtoSchema.parse(data), [data])

  return <BookContentEditor content={chapter.content} />
}
