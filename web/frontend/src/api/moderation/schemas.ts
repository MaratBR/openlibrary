import z from 'zod'
import type {
  BookModerationLogResponse,
  LoginHistoryEntryResponse,
  ModerationBookResponse,
  UserLoginHistoryResponse,
} from '@/backend-types'

export const ModerationBookSchema: z.ZodType<ModerationBookResponse> = z.object({
  id: z.string(),
  name: z.string(),
  summary: z.string(),
  isBanned: z.boolean(),
  isShadowBanned: z.boolean(),
  isPermanentlyRemoved: z.boolean(),
})

export const BookModerationLogEntrySchema = z.object({
  time: z.string(),
  action: z.string(),
  payload: z.union([
    z.null(),
    z.boolean(),
    z.number(),
    z.string(),
    z.array(z.any()),
    z.record(z.any()),
  ]),
  reason: z.string(),
  actorUserId: z.string(),
  actorUserName: z.string(),
})

export const BookModerationLogSchema: z.ZodType<BookModerationLogResponse> = z.object({
  entries: z.array(BookModerationLogEntrySchema),
  page: z.number(),
  pageSize: z.number(),
  hasNextPage: z.boolean(),
  hasPreviousPage: z.boolean(),
  totalPages: z.number(),
})

export const LoginHistoryEntrySchema: z.ZodType<LoginHistoryEntryResponse> = z.object({
  ipAddress: z.string(),
  userAgent: z.string(),
  loggedInAt: z.string(),
})

export const UserLoginHistorySchema: z.ZodType<UserLoginHistoryResponse> =
  z.array(LoginHistoryEntrySchema)

export type ModerationBook = ModerationBookResponse
export type BookModerationLog = BookModerationLogResponse
export type LoginHistoryEntry = LoginHistoryEntryResponse
