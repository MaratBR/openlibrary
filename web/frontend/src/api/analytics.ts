import z from 'zod'

export const ViewsSchema = z.object({
  total: z.number().int(),
  year: z.number().int(),
  month: z.number().int(),
  week: z.number().int(),
  day: z.number().int(),
  hour: z.number().int(),
})
