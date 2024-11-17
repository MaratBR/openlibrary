import React, { useMemo } from 'react'
import sanitize from 'sanitize-html'

const SanitizeHtml = React.forwardRef(
  ({ html }: { html: string }, ref: React.ForwardedRef<HTMLDivElement>) => {
    const safeHtml = useMemo(() => sanitize(html), [html])

    return <div ref={ref} dangerouslySetInnerHTML={{ __html: safeHtml }} />
  },
)

export default SanitizeHtml
