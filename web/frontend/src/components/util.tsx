import clsx from 'clsx'
import type { ClassValue } from 'clsx'
import { ComponentType, HTMLAttributes } from 'react'
import { forwardRef } from 'react'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function createClassComponent(
  classes: string,
): ComponentType<HTMLAttributes<HTMLDivElement>> {
  return forwardRef<HTMLDivElement, HTMLAttributes<HTMLDivElement>>(
    ({ className, ...props }, ref) => (
      <div className={clsx(classes, className)} {...props} ref={ref} />
    ),
  )
}
