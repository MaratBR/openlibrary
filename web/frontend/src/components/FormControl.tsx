import { ReactNode } from 'react'

export function FormControl({
  label,
  htmlFor,
  children,
  description,
}: {
  label: string
  htmlFor?: string
  children?: ReactNode
  description?: string
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

      <div className="form-control__value">{children}</div>
    </div>
  )
}
