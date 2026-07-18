import { Wrapper } from '@/preact'
import { ReactNode, StrictMode, useEffect } from 'react';

export function ReactIslandSetup({ children }: { children?: ReactNode }) {
  return <StrictMode>
    <Wrapper>{children}</Wrapper>
  </StrictMode>
}
