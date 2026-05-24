import { ComponentChild } from 'preact'

export function FormControl({
  label,
  htmlFor,
  children,
  description,
}: {
  label: string
  htmlFor?: string
  children?: ComponentChild
  description?: string
}) {
  return (
    <div class="form-control">
      <div class="form-control__label ">
        {htmlFor ? (
          <label class="label" htmlFor={htmlFor}>
            {label}
          </label>
        ) : (
          <span class="label">{label}</span>
        )}

        {description && <p class="form-control__hint">{description}</p>}
      </div>

      <div className="form-control__value">{children}</div>
    </div>
  )
}
