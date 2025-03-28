function getErrorMessage(error: unknown): string {
  if (typeof error === 'string') return error
  if (error instanceof Error) return error.message
  return 'Unknown error'
}

export type ErrorDisplayProps = {
  error: unknown
}

export function ErrorDisplay({ error }: ErrorDisplayProps) {
  return (
    <div class="ol-error">
      <p class="ol-error__message">{getErrorMessage(error)}</p>
    </div>
  )
}
