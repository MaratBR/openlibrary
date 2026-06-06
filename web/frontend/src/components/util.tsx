import clsx from 'clsx'
import { ComponentType, HTMLAttributes } from 'preact'
import { forwardRef } from 'preact/compat'

export function createClassComponent(
  classes: string,
): ComponentType<HTMLAttributes<HTMLDivElement>> {
  return forwardRef(
    ({ class: class_, className, ...props }: HTMLAttributes<HTMLDivElement>, ref) => (
      <div class={clsx(classes, className, class_)} {...props} />
    ),
  )
}
