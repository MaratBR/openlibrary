import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { ComponentChildren } from 'preact'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 0,
      gcTime: 50000,
    },
  },
})

export function PreactIslandSetup({ children }: { children?: ComponentChildren }) {
  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
}
