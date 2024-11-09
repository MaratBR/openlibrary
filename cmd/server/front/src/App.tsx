import { RouterProvider } from 'react-router-dom'
import router from './router'
import { QueryClientProvider } from '@tanstack/react-query'
import queryClient, { QueryClientLoadingBar } from './query-client'
import { TooltipProvider } from './components/ui/tooltip'
import { LoadingBarProvider } from './components/loading-bar-context'
import { ScrollArea } from './components/ui/scroll-area'
import { Suspense } from 'react'

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <TooltipProvider>
        <LoadingBarProvider>
          <QueryClientLoadingBar />
          <ScrollArea className="w-screen h-screen">
            <Suspense>
              <RouterProvider router={router} />
            </Suspense>
          </ScrollArea>
        </LoadingBarProvider>
      </TooltipProvider>
    </QueryClientProvider>
  )
}

export default App
