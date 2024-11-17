import React from 'react'
import ReactDOM from 'react-dom'

export const SlotContext = React.createContext<HTMLElement | null>(null)

export function SlotContent({ children }: React.PropsWithChildren) {
  const slot = React.useContext(SlotContext)

  if (!slot) return null

  return ReactDOM.createPortal(<>{children}</>, slot)
}
