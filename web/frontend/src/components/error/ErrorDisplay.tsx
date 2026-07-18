import { getErrorMessage } from '@/common/error'

export type ErrorDisplayProps = {
  error: unknown
}

export function ErrorDisplay({ error }: ErrorDisplayProps) {
  return (
    <div className="error">
      <p className="error__message">{getErrorMessage(error)}</p>
    </div>
  )
}
