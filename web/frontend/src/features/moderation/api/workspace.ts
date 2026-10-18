import { httpClient, OLAPIResponse } from '@/features/http-client'
import type {
  ModerationContentPage,
  ModerationContentEntry,
  ModerationOverview,
  ModerationReportNoteRequest,
} from '@/backend-types'

export async function moderationRead<T>(
  path: string,
  filters: Record<string, string> = {},
  signal?: AbortSignal,
): Promise<T> {
  filters = Object.fromEntries(
    Object.entries(filters).map(([key, value]) => [
      key,
      value === 'all' && key !== 'status' ? '' : value,
    ]),
  )
  const response = await httpClient.get(`/_api/moderation/${path}`, {
    searchParams: filters,
    signal,
  })
  const result = await OLAPIResponse.create<T>(response)
  result.throwIfError()
  return result.data
}
export const getWorkspaceOverview = (signal?: AbortSignal) =>
  moderationRead<ModerationOverview>('overview', {}, signal)
export const searchContent = (
  kind: 'chapters' | 'comments',
  filters: Record<string, string>,
  signal?: AbortSignal,
) => moderationRead<ModerationContentPage>(kind, filters, signal)
export const getContent = (kind: 'chapters' | 'comments', id: string, signal?: AbortSignal) =>
  moderationRead<ModerationContentEntry>(`${kind}/${encodeURIComponent(id)}`, {}, signal)
export async function addReportNote(id: string, body: ModerationReportNoteRequest) {
  const response = await httpClient.post(
    `/_api/moderation/reports/${encodeURIComponent(id)}/notes`,
    { json: body },
  )
  const result = await OLAPIResponse.createNoBody(response)
  result.throwIfError()
}
