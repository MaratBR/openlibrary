import { ReactNode, useRef } from 'react'

export function RenderLazy({ show, children }: { show: boolean; children: ReactNode }) {
  const shown = useRef(false)

  if (show) {
    shown.current = true
  }

  if (shown.current) return <div style={{ display: show ? 'contents' : 'none' }}>{children}</div>
  return null
}
