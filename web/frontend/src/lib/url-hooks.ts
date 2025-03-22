import { useMemo, useState } from 'preact/hooks'

function subscribeToUrlHashChange(callback: (hash: string) => void): () => void {
  const fn = () => {
    callback(window.location.hash)
  }
  window.addEventListener('hashchange', fn)
  return () => window.removeEventListener('hashchange', fn)
}

export function useHash(): string {
  const [hash, setHash] = useState(window.location.hash)

  subscribeToUrlHashChange(setHash)

  return hash
}

export function useHashQuery(): URLSearchParams {
  const hash = useHash()

  return useMemo(() => {
    if (hash === '') return new URLSearchParams()
    try {
      return new URLSearchParams(hash.slice(1))
    } catch {
      return new URLSearchParams()
    }
  }, [hash])
}

function setHashQuery(key: string, value: string) {
  let hash = window.location.hash
  if (hash.startsWith('#')) hash = hash.slice(1)
  const searchParams = new URLSearchParams(hash)
  searchParams.delete(key)
  searchParams.set(key, value)
  const queryString = searchParams.toString()
  window.location.hash = queryString
}

export function useHashQueryValue(key: string): [string | null, (value: string) => void] {
  const searchParams = useHashQuery()

  const value = searchParams.get(key)

  const setValue = (value: string) => {
    setHashQuery(key, value)
  }

  return [value, setValue]
}
