import { getCsrfToken } from '@/http-client/client'
import { useMemo } from 'react'

export default function CSRFInput() {
  const value = useMemo(() => getCsrfToken(), [])

  return <input name="__csrf" hidden value={value} />
}
