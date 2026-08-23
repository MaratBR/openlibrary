import { z } from 'zod'

export const managerBookChapterDetailsDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  createdAt: z.string(),
  words: z.number(),
  summary: z.string(),
  order: z.number().int(),
  content: z.string(),
  isPubliclyVisible: z.boolean(),
  fonts: z.array(z.string()),
})
