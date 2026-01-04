export function getCookie(name: string): string | undefined {
  const value = `; ${document.cookie}`
  const parts = value.split(`; ${name}=`)
  if (parts.length === 2) return parts.pop()!.split(';').shift()
}

export function setCookie(
  name: string,
  value: string,
  options?: {
    days?: number // expiration in days (default = 7300 ≈ 20 years)
    path?: string
    domain?: string
    secure?: boolean
    sameSite?: 'Strict' | 'Lax' | 'None'
  },
): void {
  if (typeof document === 'undefined') return // SSR/worker safety

  const { days = 365 * 20, path = '/', domain, secure = false, sameSite = 'Lax' } = options || {}

  const expires = new Date()
  expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000)

  let cookie = `${encodeURIComponent(name)}=${encodeURIComponent(value)}`
  cookie += `; expires=${expires.toUTCString()}`
  cookie += `; path=${path}`
  if (domain) cookie += `; domain=${domain}`
  if (secure) cookie += `; secure`
  if (sameSite) cookie += `; samesite=${sameSite}`

  document.cookie = cookie
}

declare global {
  interface Window {
    getCookie: typeof getCookie
    setCookie: typeof setCookie
  }
}

window.setCookie = setCookie
window.getCookie = getCookie
