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
    <div className="form-control">
      <div className="form-control__label ">
        {htmlFor ? (
          <label className="label" htmlFor={htmlFor}>
            {label}
          </label>
        ) : (
          <span className="label">{label}</span>
        )}

        {description && <p className="form-control__hint">{description}</p>}
      </div>

      <div className="form-control__value">
        {children}
        {error && (
          <p
            id={htmlFor ? `${htmlFor}-error` : undefined}
            className="form-control__error"
            role="alert"
          >
            <i className="fa-solid fa-circle-exclamation" aria-hidden="true" />
            {error}
          </p>
        )}
      </div>
    </div>
  )
}
