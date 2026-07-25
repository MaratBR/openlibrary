import { httpClient, OLAPIResponse } from '@/http-client'
import { UserLoginHistorySchema } from './schemas'

export function banUser(userId: string, reason: string, until: Date | string) {
  const expiresAt = until instanceof Date ? until.toISOString() : until
  return userAction(userId, 'ban', { reason, until: expiresAt })
}

export function permanentlyBanUser(userId: string, reason: string) {
  return userAction(userId, 'permanent-ban', { reason })
}

export function unbanUser(userId: string, reason: string) {
  return userAction(userId, 'unban', { reason })
}

export function renameUser(userId: string, value: string, reason: string) {
  return userAction(userId, 'rename', { reason, value })
}

export function changeUserAbout(userId: string, value: string, reason: string) {
  return userAction(userId, 'change-about', { reason, value })
}

export function getUserLoginHistory(userId: string) {
  return httpClient
    .get(`/_api/moderation/users/${encodeURIComponent(userId)}/login-history`)
    .then((response) => OLAPIResponse.create(response, UserLoginHistorySchema))
}

function userAction(userId: string, action: string, body: Record<string, string>) {
  return httpClient
    .post(`/_api/moderation/users/${encodeURIComponent(userId)}/actions/${action}`, { json: body })
    .then((response) => OLAPIResponse.createNoBody(response))
}
