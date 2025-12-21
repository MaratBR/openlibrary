import { QueryClientProvider } from '@tanstack/react-query'
import { Attributes, Component, ComponentChildren } from 'preact'
import { preactQueryCache } from './queryCache'
import React from 'react-dom/src'
import { ErrorDisplay } from '@/components/error'

export default function Wrapper({ children }: { children: ComponentChildren }) {
  return (
    <ErrorBoundary>
      <QueryClientProvider client={preactQueryCache}>{children}</QueryClientProvider>
    </ErrorBoundary>
  )
}

class ErrorBoundary extends Component<object, { error?: unknown }> {
  render(
    props?:
      | Readonly<
          Attributes & { children?: ComponentChildren; ref?: React.Ref<unknown> | undefined }
        >
      | undefined,
    state?: Readonly<{ error?: unknown }> | undefined,
    _context?: unknown,
  ): ComponentChildren {
    if (state && state.error) {
      return <ErrorDisplay error={state.error} />
    }

    return props?.children
  }

  componentDidCatch(error: unknown, _errorInfo: React.ErrorInfo): void {
    this.setState({
      error,
    })
  }
}
