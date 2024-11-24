import ky, { KyResponse } from 'ky'
import { ZodSchema } from 'zod'

export const httpClient = ky.create({
  timeout: 60000,
  hooks: {
    beforeRequest: [
      (req) => {
        if (!['GET', 'HEAD', 'OPTIONS'].includes(req.method)) {
          const csrfToken = getCsrfToken()
          if (csrfToken) {
            req.headers.set('x-csrf-token', csrfToken)
          }
        }
      },
    ],
  },
})

// setTimeout(() => {
//   httpClient.post('/api/auth/csrf-check')
// }, 200)

function getCsrfToken() {
  try {
    return getCookie('csrf')
  } catch {
    /* empty */
  }
}

function refreshCsrfToken() {
  fetch('/api/auth/csrf', { method: 'GET' })
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

export function getPreloadedData<T>(key: string): T | undefined {
  const value = __server__._preload?.[key]
  if (value === undefined) return undefined
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

export function stringArray(arr: string[]): string | undefined {
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
