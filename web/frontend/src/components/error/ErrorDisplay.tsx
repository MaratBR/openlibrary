import { getErrorMessage } from '@/common/error'

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
