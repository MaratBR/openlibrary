import { getCsrfToken } from '@/http-client/client'
import { useMemo } from 'preact/hooks'

export default function CSRFInput() {
  const value = useMemo(() => getCsrfToken(), [])

  return <input hidden value={value} />
}
