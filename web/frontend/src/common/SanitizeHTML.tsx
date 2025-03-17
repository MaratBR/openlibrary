import { useMemo } from 'preact/hooks'
import sanitizeHtml from 'sanitize-html'

export default function SanitizeHTML({ value }: { value: string }) {
  const sanitized = useMemo(() => sanitizeHtml(value), [value])

  return <div dangerouslySetInnerHTML={{ __html: sanitized }} />
}
