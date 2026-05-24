import { Wrapper } from '@/preact'
import { ComponentChildren } from 'preact'

export function PreactIslandSetup({ children }: { children?: ComponentChildren }) {
  return <Wrapper>{children}</Wrapper>
}
