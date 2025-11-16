import { OLAPIResponse } from '@/http-client'

export function getErrorMessage(error: unknown): string {
  if (typeof error === 'string') return error
  if (error instanceof Error) return error.message
  if (error instanceof OLAPIResponse) {
    return error.error?.message ?? 'Unknown API error'
  }
  return 'Unknown error'
}

declare global {
  interface Window {
    getErrorMessage: typeof getErrorMessage
  }
}

window.getErrorMessage = getErrorMessage
