import { httpClient, OLAPIResponse } from '@/http-client'
import type {
  ModerationUserBooksPageResponse,
  ModerationUserCommentsPageResponse,
  ModerationUserHistoryPageResponse,
  ModerationUserReportsPageResponse,
  ModerationUserResponse,
  ModerationUsersPageResponse,
  ModerationBanRequest,
  ModerationReasonRequest,
  ModerationValueRequest,
  UserLoginHistoryResponse,
  LoginLocationResponse,
  ModerationReportDetailResponse,
  ModerationReportsSearchResponse,
} from '@/backend-types'

export function searchModerationUsers(
  search = '',
  banned = '',
  role = '',
  page = 1,
  pageSize = 20,
) {
  return httpClient
    .get('/_api/moderation/users', { searchParams: { search, banned, role, page, pageSize } })
    .then((response) => OLAPIResponse.create<ModerationUsersPageResponse>(response))
}

export function getModerationUser(userId: string) {
  return httpClient
    .get(`/_api/moderation/users/${encodeURIComponent(userId)}`)
    .then((response) => OLAPIResponse.create<ModerationUserResponse>(response))
}

export function banUser(userId: string, reason: string, until: Date | string) {
  const expiresAt = until instanceof Date ? until.toISOString() : until
  const body: ModerationBanRequest = { reason, until: expiresAt }
  return userAction(userId, 'ban', body)
}

export function permanentlyBanUser(userId: string, reason: string) {
  return userAction(userId, 'permanent-ban', { reason } satisfies ModerationReasonRequest)
}

export function unbanUser(userId: string, reason: string) {
  return userAction(userId, 'unban', { reason } satisfies ModerationReasonRequest)
}

export function renameUser(userId: string, value: string, reason: string) {
  return userAction(userId, 'rename', { reason, value } satisfies ModerationValueRequest)
}

export function changeUserAbout(userId: string, value: string, reason: string) {
  return userAction(userId, 'change-about', { reason, value } satisfies ModerationValueRequest)
}

export function getUserLoginHistory(userId: string, page = 1, pageSize = 20) {
  return httpClient
    .get(`/_api/moderation/users/${encodeURIComponent(userId)}/login-history`, {
      searchParams: { page, pageSize },
    })
    .then((response) => OLAPIResponse.create<UserLoginHistoryResponse>(response))
}

export function getUserLoginLocations(userId: string) {
  return httpClient
    .get(`/_api/moderation/users/${encodeURIComponent(userId)}/login-locations`)
    .then((response) => OLAPIResponse.create<LoginLocationResponse[]>(response))
}

export function getModerationReport(reportId: string) {
  return httpClient
    .get(`/_api/moderation/reports/${encodeURIComponent(reportId)}`)
    .then((response) => OLAPIResponse.create<ModerationReportDetailResponse>(response))
}

export function searchModerationReports(search = '', targetType = '', page = 1, pageSize = 20) {
  return httpClient
    .get('/_api/moderation/reports', { searchParams: { search, targetType, page, pageSize } })
    .then((response) => OLAPIResponse.create<ModerationReportsSearchResponse>(response))
}

function getUserPage<T>(userId: string, resource: string, page = 1, pageSize = 20) {
  return httpClient
    .get(`/_api/moderation/users/${encodeURIComponent(userId)}/${resource}`, {
      searchParams: { page, pageSize },
    })
    .then((response) => OLAPIResponse.create<T>(response))
}
export function getModerationUserBooks(userId: string, page = 1, pageSize = 20) {
  return getUserPage<ModerationUserBooksPageResponse>(userId, 'books', page, pageSize)
}
export function getModerationUserComments(userId: string, page = 1, pageSize = 20) {
  return getUserPage<ModerationUserCommentsPageResponse>(userId, 'comments', page, pageSize)
}
export function getModerationUserHistory(userId: string, page = 1, pageSize = 20) {
  return getUserPage<ModerationUserHistoryPageResponse>(userId, 'history', page, pageSize)
}
export function getModerationUserReports(userId: string, page = 1, pageSize = 20) {
  return getUserPage<ModerationUserReportsPageResponse>(userId, 'reports', page, pageSize)
}

function userAction(
  userId: string,
  action: string,
  body: ModerationBanRequest | ModerationReasonRequest | ModerationValueRequest,
) {
  return httpClient
    .post(`/_api/moderation/users/${encodeURIComponent(userId)}/actions/${action}`, { json: body })
    .then((response) => OLAPIResponse.createNoBody(response))
}
