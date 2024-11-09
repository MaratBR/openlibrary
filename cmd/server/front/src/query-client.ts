import { QueryClient, useIsFetching } from '@tanstack/react-query'
import { useHeaderLoading } from './components/loading-bar-context'
import { HTTPError } from 'ky'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: (failureCount, err) => {
        if (failureCount > 5) return false

        if (err instanceof HTTPError) {
          if (err.response.status === 429) return false
          if (err.response.status === 401) return false
          if (err.response.status === 403) return false
        }

        return true
      },
    },
  },
})

export default queryClient

export function QueryClientLoadingBar() {
  const fetchingCount = useIsFetching()
  useHeaderLoading(fetchingCount > 0)
  return null
}
