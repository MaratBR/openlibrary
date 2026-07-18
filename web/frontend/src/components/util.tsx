import clsx from 'clsx'
import { ComponentType, HTMLAttributes } from 'react'
import { forwardRef } from 'react'

export function createClassComponent(
  classes: string,
): ComponentType<HTMLAttributes<HTMLDivElement>> {
  return forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
    ({ className, ...props }, ref) => (
      <div className={clsx(classes, className)} {...props} ref={ref} />
    ),
  )
}
