import { z } from 'zod'

export const DraftDtoSchema = z.object({
  id: z.string(),
  chapterName: z.string(),
  content: z.string(),
  createdAt: z.string(),
  updatedAt: z.string().nullable(),
  chapterId: z.string(),
  createdBy: z.object({
    id: z.string().uuid(),
    name: z.string(),
  }),
})

export type DraftDto = z.infer<typeof DraftDtoSchema>
