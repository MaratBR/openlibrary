import { QueryClient } from '@tanstack/react-query'

export const preactQueryCache = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 0,
      gcTime: 50000,
    },
  },
})
