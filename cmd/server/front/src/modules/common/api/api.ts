import ky, { KyResponse } from 'ky'
import { ZodSchema } from 'zod'

const delayedFetch: typeof window.fetch = async (...args): Promise<Response> => {
  const delay = 0 // delay in milliseconds
  await new Promise((resolve) => setTimeout(resolve, delay))
  return window.fetch(...args)
}

const originalFetch = window.fetch

// Override the global fetch
window.fetch = async function (input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
  // Create new init object if it doesn't exist
  const modifiedInit: RequestInit = init ? { ...init } : {}

  // Create headers object if it doesn't exist
  modifiedInit.headers = new Headers(modifiedInit.headers || {})

  // Add CSRF token to headers if not already present
  if (!modifiedInit.headers.has('x-csrf-token')) {
    const token = getCsrfToken()
    if (token) {
      modifiedInit.headers.set('x-csrf-token', token)
    }
  }

  // Call original fetch with modified init object
  return originalFetch.call(window, input, modifiedInit)
}

export const httpClient = ky.create({
  timeout: 60000,
  fetch: delayedFetch,
})

function getCsrfToken() {
  try {
    return getCookie('csrf')
  } catch {
    /* empty */
  }
}

function getCookie(name: string): string | undefined {
  const value = `; ${document.cookie}`
  const parts = value.split(`; ${name}=`)
  if (parts.length === 2) return parts.pop()!.split(';').shift()
}

const disableWithPreloadCache =
  new URLSearchParams(window.location.search).get('debug.disableWithPreloadCache') === 'true' ||
  !__server__.serverPreload

export async function withPreloadCache<T>(key: string, fn: () => Promise<T>): Promise<T> {
  if (disableWithPreloadCache) {
    return await fn()
  }

  if (__server__._preload && __server__._preload[key]) {
    const value = __server__._preload[key] as T
    delete __server__._preload[key]

    return value
  } else {
    console.warn(`Preload cache miss: ${key}`)
    return await fn()
  }
}

export function pullPreloadedData<T>(key: string): T | undefined {
  const value = __server__._preload?.[key]
  if (value === undefined) return undefined
  if (__server__._preload) {
    delete __server__._preload[key]
  }
  return value as T
}

export function parseQueryStringArray(queryString: string | undefined | null): string[] {
  if (!queryString) return []
  return queryString.split(/(?<!\\)\|/g)
}

export function isSearchQueryEqual(a: URLSearchParams, b: URLSearchParams): boolean {
  if (a.size !== b.size) return false
  for (const key of a.keys()) {
    const av = a.getAll(key)
    const bv = b.getAll(key)
    av.sort()
    bv.sort()
    if (av.length !== bv.length) return false
    for (let i = 0; i < av.length; i++) {
      if (av[i] !== bv[i]) return false
    }
  }
  return true
}

export function stringArrayToQueryParameterValue(arr: string[]): string | undefined {
  if (arr.length === 0) return undefined
  const sortedArr = [...arr]
  sortedArr.sort()
  return sortedArr.map((x) => x.replace('|', '\\|')).join('|')
}

export async function getResponse<T>(resp: KyResponse, schema: ZodSchema<T>): Promise<T> {
  const json = await resp.json()
  return schema.parse(json)
}

export function responseSchema<T>(schema: ZodSchema<T>): (resp: KyResponse) => Promise<T> {
  return (resp) => getResponse(resp, schema)
}
