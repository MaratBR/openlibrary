import { useMemo } from 'preact/hooks'
import z from 'zod'

export const UserRoleSchema = z.enum(['user', 'admin', 'system', 'moderator'])

export const SelfUserDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  email: z.string(),
  role: UserRoleSchema,
  avatar: z.object({
    lg: z.string(),
    md: z.string(),
  }),
  joinedAt: z.string(),
  isBanned: z.boolean(),
  isEmailVerified: z.boolean(),
  preferredTheme: z.string(),
})

export function useUserSelfData() {
  return useMemo(() => SelfUserDtoSchema.parse(window.__server__.user), [])
}
