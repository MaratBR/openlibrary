import { Wrapper } from '@/preact'
import { ReactNode, StrictMode } from 'react'

export function ReactIslandSetup({ children }: { children?: ReactNode }) {
  return (
    <StrictMode>
      <Wrapper>{children}</Wrapper>
    </StrictMode>
  )
}
