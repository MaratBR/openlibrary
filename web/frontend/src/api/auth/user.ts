import { useMemo } from 'react'
import { Schema } from 'effect'

export const UserRoleSchema = Schema.Literals(['user', 'admin', 'system', 'moderator'])

export const SelfUserDtoSchema = Schema.Struct({
  id: Schema.String,
  name: Schema.String,
  email: Schema.String,
  role: UserRoleSchema,
  avatar: Schema.Struct({
    lg: Schema.String,
    md: Schema.String,
  }),
  joinedAt: Schema.String,
  isBanned: Schema.Boolean,
  isEmailVerified: Schema.Boolean,
  preferredTheme: Schema.String,
})

export function useUserSelfData() {
  return useMemo(() => Schema.decodeUnknownSync(SelfUserDtoSchema)(window.__server__.user), [])
}
