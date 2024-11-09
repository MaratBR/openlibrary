import React from 'react'
import sanitize from 'sanitize-html'

export function NotificationContent({ content }: { content: string }) {
  const html = React.useMemo(
    () =>
      sanitize(replaceLinks(content), { nestingLimit: 5, allowedClasses: { a: ['link-default'] } }),
    [content],
  )

  return <div dangerouslySetInnerHTML={{ __html: html }} />
}

/**
 * Replaces markdown-style links with a tags
 */
function replaceLinks(content: string): string {
  return content.replace(/\[([^\]]+)\]\((\/[^\s)]+)\)/g, '<a class="link-default" href="$2">$1</a>')
}
