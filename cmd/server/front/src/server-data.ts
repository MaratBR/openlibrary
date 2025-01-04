import { z } from 'zod'
import { selfUserDtoSchema } from './modules/user/api'

const serverDataSchema = z.object({
  serverPreload: z.boolean(),
  clientPreload: z.boolean(),
  _preload: z.record(z.unknown()).optional(),
  iframeAllowed: z.boolean().optional(),
  session: z
    .object({
      expiresAt: z.string(),
    })
    .optional()
    .nullable(),
  user: selfUserDtoSchema.optional().nullable(),
  /**
   * Base64 encoded protobuf message containing results of the search
   */
  search: z.string().optional(),
})

declare global {
  let __server__: z.infer<typeof serverDataSchema>
}

{
  const result = serverDataSchema.safeParse(__server__)
  if (!result.success) {
    alert('__server__ data is not valid')
    console.error(result.error)
  }
}
