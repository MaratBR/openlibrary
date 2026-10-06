import { Schema } from 'effect'
import { httpClient } from '@/features/http-client'

export const roles = ['user', 'admin', 'system', 'moderator'] as const
export const categories = ['other', 'warning', 'fandom', 'rel', 'reltype', 'genre'] as const
const userSchema = Schema.Struct({
  id: Schema.String,
  name: Schema.String,
  role: Schema.Literals(roles),
  isBanned: Schema.Boolean,
  avatar: Schema.String,
  joinedAt: Schema.String,
  bio: Schema.String,
  gender: Schema.String,
})
const parentSchema = Schema.Struct({ id: Schema.String, name: Schema.String })
const tagSchema = Schema.Struct({
  id: Schema.String,
  name: Schema.String,
  description: Schema.String,
  category: Schema.Literals(categories),
  adult: Schema.Boolean,
  spoiler: Schema.Boolean,
  createdAt: Schema.String,
  isDefault: Schema.Boolean,
  synonymOf: Schema.NullOr(parentSchema),
})
const usersSchema = Schema.Struct({
  users: Schema.Array(userSchema),
  page: Schema.Number,
  totalPages: Schema.Number,
  total: Schema.Number,
})
const tagsSchema = Schema.Struct({
  tags: Schema.Array(tagSchema),
  page: Schema.Number,
  totalPages: Schema.Number,
  total: Schema.Number,
})
const successSchema = Schema.Struct({ ok: Schema.Literal(true) })
export type AdminUser = typeof userSchema.Type
export type AdminTag = typeof tagSchema.Type

export class AdminAPIError extends Error {
  constructor(public status: number) {
    super(`Admin request failed: ${status}`)
  }
}
async function read<S extends Schema.ConstraintDecoder<unknown>>(
  path: string,
  schema: S,
  signal?: AbortSignal,
): Promise<S['Type']> {
  const response = await httpClient.get(`/admin/api${path}`, {
    signal,
    throwHttpErrors: false,
    retry: 0,
  })
  if (!response.ok) throw new AdminAPIError(response.status)
  return Schema.decodeUnknownSync(schema)(await response.json())
}
export async function save(path: string, fields: Record<string, string>) {
  const response = await httpClient.post(`/admin/api${path}`, {
    body: new URLSearchParams(fields),
    throwHttpErrors: false,
    retry: 0,
  })
  if (!response.ok) throw new AdminAPIError(response.status)
  Schema.decodeUnknownSync(successSchema)(await response.json())
}
export const loadUsers = (search: string, signal?: AbortSignal) =>
  read(`/users${search}`, usersSchema, signal)
export const loadTags = (search: string, signal?: AbortSignal) =>
  read(`/tags${search}`, tagsSchema, signal)
export const loadUser = (id: string, signal?: AbortSignal) =>
  read(`/users/${validUserID(id)}`, userSchema, signal)
export const loadTag = (id: string, signal?: AbortSignal) =>
  read(`/tags/${validTagID(id)}`, tagSchema, signal)
export const loadDebug = (signal?: AbortSignal) =>
  read('/debug', Schema.Array(Schema.String), signal)
export function validTagID(id: string) {
  if (!/^[1-9]\d*$/.test(id) || BigInt(id) > 9223372036854775807n) throw new AdminAPIError(400)
  return id
}
export function validUserID(id: string) {
  if (!/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(id))
    throw new AdminAPIError(400)
  return id
}
