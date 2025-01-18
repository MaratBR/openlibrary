import ky from 'ky'

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
