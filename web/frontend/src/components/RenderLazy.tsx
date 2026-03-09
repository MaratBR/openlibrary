import { ComponentChildren } from 'preact'
import { useRef } from 'preact/hooks'

export function RenderLazy({ show, children }: { show: boolean; children: ComponentChildren }) {
  const shown = useRef(false)

  if (show) {
    shown.current = true
  }

  if (shown.current) return <div style={{ display: show ? 'contents' : 'none' }}>{children}</div>
  return null
}
