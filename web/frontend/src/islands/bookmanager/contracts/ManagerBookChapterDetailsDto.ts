import { z } from 'zod'

export const managerBookChapterDetailsDtoSchema = z.object({
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
