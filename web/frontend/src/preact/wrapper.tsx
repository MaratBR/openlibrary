import { QueryClientProvider } from '@tanstack/react-query'
import { ComponentChildren } from 'preact'
import { preactQueryCache } from './queryCache'

export default function Wrapper({ children }: { children: ComponentChildren }) {
  return <QueryClientProvider client={preactQueryCache}>{children}</QueryClientProvider>
}
