import { Wrapper } from '@/react'
import { ReactNode, StrictMode } from 'react'

export function ReactIslandSetup({ children }: { children?: ReactNode }) {
  return (
    <StrictMode>
      <Wrapper>{children}</Wrapper>
    </StrictMode>
  )
}
