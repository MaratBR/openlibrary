import { getErrorMessage } from '@/common/error'

export type ErrorDisplayProps = {
  error: unknown
}

export function ErrorDisplay({ error }: ErrorDisplayProps) {
  return (
    <div className="error" role="alert">
      <svg
        className="error__icon"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
        aria-hidden="true"
      >
        <circle cx="12" cy="12" r="9" />
        <path d="M12 8v5" />
        <path d="M12 16.5h.01" />
      </svg>
      <p className="error__message">{getErrorMessage(error)}</p>
    </div>
  )
}
