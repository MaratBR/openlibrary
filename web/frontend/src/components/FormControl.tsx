import { Children, cloneElement, isValidElement, ReactElement, ReactNode, useId } from 'react'
import { SelectTrigger } from './Select'

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
  const generatedId = useId()
  const controlId = htmlFor ?? generatedId
  const labelId = `${controlId}-label`
  const associate = (nodes: ReactNode): ReactNode =>
    Children.map(nodes, (child) => {
      if (!isValidElement(child)) return child
      const element = child as ReactElement<{
        children?: ReactNode
        id?: string
        'aria-labelledby'?: string
      }>
      if (
        ['input', 'textarea', 'select'].includes(String(element.type)) ||
        element.type === SelectTrigger
      ) {
        return cloneElement(element, {
          id: element.props.id ?? controlId,
          'aria-labelledby': element.props['aria-labelledby'] ?? labelId,
        })
      }
      return element.props.children
        ? cloneElement(element, {}, associate(element.props.children))
        : element
    })
  return (
    <div className="FormControl">
      <div className="FormControl-label ">
        <label id={labelId} className="label" htmlFor={controlId}>
          {label}
        </label>

        {description && <p className="FormControl-hint">{description}</p>}
      </div>

      <div className="FormControl-value">
        {associate(children)}
        {error && (
          <p
            id={htmlFor ? `${htmlFor}-error` : undefined}
            className="FormControl-error"
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
