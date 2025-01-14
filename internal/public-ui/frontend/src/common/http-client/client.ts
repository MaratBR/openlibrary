import ky from 'ky';

const delayedFetch: typeof window.fetch = async (...args): Promise<Response> => {
  const delay = 0 // delay in milliseconds
  await new Promise((resolve) => setTimeout(resolve, delay))
  return window.fetch(...args)
}

export const httpClient = ky.create({
  timeout: 60000,
  fetch: delayedFetch,
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
