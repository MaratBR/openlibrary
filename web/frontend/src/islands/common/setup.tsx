import Wrapper from '@/preact/wrapper'
import { ComponentChildren } from 'preact'

export function PreactIslandSetup({ children }: { children?: ComponentChildren }) {
  return <Wrapper>{children}</Wrapper>
}
