import { HTMLAttributes } from 'preact'
import { useMemo } from 'preact/hooks'
import sanitizeHtml from 'sanitize-html'

export default function SanitizeHTML({
  value,
  ...props
}: { value: string } & HTMLAttributes<HTMLDivElement>) {
  const sanitized = useMemo(() => sanitizeHtml(value), [value])

  // eslint-disable-next-line react/no-danger
  return <div {...props} dangerouslySetInnerHTML={{ __html: sanitized }} />
}
