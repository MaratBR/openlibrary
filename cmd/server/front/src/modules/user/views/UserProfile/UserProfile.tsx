import { getIframeUrl } from '@/lib/iframe-auto-resize'
import React, { useMemo } from 'react'
import { useParams } from 'react-router'

export default function UserProfile() {
  const { userId } = useParams<{ userId: string }>()
  const [loaded, setLoaded] = React.useState(false)

  const id = `iframe${userId?.replace(/-/g, '')}`

  const iframeUrl = useMemo(
    () => getIframeUrl(`/users/__profile?userId=${userId}`, false, true),
    [userId],
  )

  return (
    <iframe
      onLoad={() => window.setTimeout(() => setLoaded(true), 0)}
      id={id}
      className="w-full transition-opacity duration-300 ease-in-out"
      style={{ opacity: loaded ? 1 : 0 }}
      src={iframeUrl}
    />
  )
}
