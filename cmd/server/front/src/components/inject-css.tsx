import React from 'react'

export type InjectCSSProps = {
  css: string
  document: DocumentOrShadowRoot
}

export function InjectCSS({ css, document: doc }: InjectCSSProps) {
  React.useEffect(() => {
    if (import.meta.env.DEV) {
      const style = document.createElement('style')
      style.textContent = css
      style.setAttribute('data-injected', 'true')
      if (doc instanceof Document || doc instanceof ShadowRoot) {
        doc.appendChild(style)
      }
      return () => style.remove()
    } else {
      const sheet = new CSSStyleSheet()
      sheet.replace(css)
      doc.adoptedStyleSheets.push(sheet)

      return () => {
        doc.adoptedStyleSheets = doc.adoptedStyleSheets.filter((s) => s !== sheet)
      }
    }
  }, [css, doc])

  return null
}
