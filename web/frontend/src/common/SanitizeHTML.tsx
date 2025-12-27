import { useMemo } from 'preact/hooks'
import sanitizeHtml from 'sanitize-html'

export default function SanitizeHTML({ value }: { value: string }) {
  const sanitized = useMemo(() => sanitizeHtml(value), [value])

  // eslint-disable-next-line react/no-danger
  return <div dangerouslySetInnerHTML={{ __html: sanitized }} />
}
