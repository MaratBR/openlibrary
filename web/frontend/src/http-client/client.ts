import ky from 'ky'

import './client-meta'
import { getCookie } from './util'

const originalFetch = window.fetch

// Override the global fetch
window.fetch = async function (
  input: Request | string | URL,
  init?: globalThis.RequestInit,
): Promise<Response> {
  // Create new init object if it doesn't exist
  const modifiedInit: globalThis.RequestInit = init ? { ...init } : {}

  let headers: Headers

  if (input instanceof Request) {
    headers = input.headers
  } else {
    headers = new Headers(modifiedInit.headers)
  }

  // Add CSRF token to headers if not already present
  const token = getCsrfToken()
  if (token) {
    headers.set('x-csrf-token', token)
  }

  modifiedInit.headers = headers

  // Call original fetch with modified init object
  return originalFetch.call(window, input, modifiedInit)
}

export const httpClient = ky.create({
  timeout: 60000,
})

export function getCsrfToken() {
  try {
    return getCookie('csrf')
  } catch {
    /* empty */
  }
}
