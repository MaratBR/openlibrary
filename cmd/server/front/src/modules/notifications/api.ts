import { z } from 'zod'

export const genericNotificationSchema = z.object({
  id: z.string(),
  text: z.string(),
})

export type GenericNotification = z.infer<typeof genericNotificationSchema>
