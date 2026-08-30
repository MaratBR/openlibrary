import { ReactNode } from 'react'

export function FormControl({
  label,
  htmlFor,
  children,
  description,
  error,
}: {
  label: string
  htmlFor?: string
  children?: ReactNode
  description?: string
  error?: string
}) {
  return (
    <div className="FormControl">
      <div className="FormControl-label ">
        {htmlFor ? (
          <label className="label" htmlFor={htmlFor}>
            {label}
          </label>
        ) : (
          <span className="label">{label}</span>
        )}

        {description && <p className="FormControl-hint">{description}</p>}
      </div>

      <div className="FormControl-value">
        {children}
        {error && (
          <p
            id={htmlFor ? `${htmlFor}-error` : undefined}
            className="FormControl-error"
            role="Alert"
          >
            <i className="fa-solid fa-circle-exclamation" aria-hidden="true" />
            {error}
          </p>
        )}
      </div>
    </div>
  )
}
