import { RouterProvider } from 'react-router-dom'
import router from './router'
import { QueryClientProvider } from '@tanstack/react-query'
import queryClient, { QueryClientLoadingBar } from './query-client'
import { TooltipProvider } from './components/ui/tooltip'
import { LoadingBarProvider } from './components/loading-bar-context'
import React, { Suspense } from 'react'
import { initIframeAgent } from './lib/iframe-auto-resize'
import { initScrollbarWidth } from './lib/scrollbar-width'

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <LoadingBarProvider>
          <QueryClientLoadingBar />
          <Suspense>
            <Initialization>
              <RouterProvider router={router} />
            </Initialization>
          </Suspense>
        </LoadingBarProvider>
      </TooltipProvider>
    </QueryClientProvider>
  )
}

function Initialization({ children }: React.PropsWithChildren) {
  return <>{children}</>
}

export default App

export function staticInitApp() {
  document.addEventListener('DOMContentLoaded', initIframeAgent)

  initScrollbarWidth()
}
