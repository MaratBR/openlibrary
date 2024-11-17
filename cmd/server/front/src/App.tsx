import { RouterProvider } from 'react-router-dom'
import router from './router'
import { QueryClientProvider } from '@tanstack/react-query'
import queryClient, { QueryClientLoadingBar } from './query-client'
import { TooltipProvider } from './components/ui/tooltip'
import { LoadingBarProvider } from './components/loading-bar-context'
import { Suspense } from 'react'
import { initIframeAgent } from './lib/iframe-auto-resize'

document.addEventListener('DOMContentLoaded', initIframeAgent)

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <LoadingBarProvider>
          <QueryClientLoadingBar />
          <Suspense>
            <RouterProvider router={router} />
          </Suspense>
        </LoadingBarProvider>
      </TooltipProvider>
    </QueryClientProvider>
  )
}

export default App
