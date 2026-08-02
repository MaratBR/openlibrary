import { QueryClientProvider } from '@tanstack/react-query'
import { Component, ReactNode } from 'react'
import { queryClient } from './queryCache'
import React from 'react'
import { ErrorDisplay } from '@/components/error'

export default function Wrapper({ children }: { children: ReactNode }) {
  return (
    <ErrorBoundary>
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </ErrorBoundary>
  )
}

type ErrorBoundaryProps = { children: ReactNode }

class ErrorBoundary extends Component<ErrorBoundaryProps, { error?: unknown }> {
  constructor(props: Readonly<ErrorBoundaryProps>) {
    super(props)
    this.state = {}
  }

  render(): ReactNode {
    if (this.state.error) {
      return <ErrorDisplay error={this.state.error} />
    }

    return this.props.children
  }

  componentDidCatch(error: unknown, errorInfo: React.ErrorInfo): void {
    console.error('[ErrorBoundary] component caught an error', error, errorInfo)
    this.setState({
      error,
    })
  }
}
