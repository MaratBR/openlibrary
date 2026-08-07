import type { HTMLAttributes, ReactNode } from 'react'
import { cn } from './util'

function TimelineRoot({ className, ...props }: HTMLAttributes<HTMLOListElement>) {
  return <ol className={cn('timeline grid gap-0', className)} {...props} />
}

function TimelineItem({
  marker,
  className,
  children,
  ...props
}: HTMLAttributes<HTMLLIElement> & { marker: ReactNode }) {
  return (
    <li className={cn('timeline__item grid grid-cols-[2rem_1fr] gap-3', className)} {...props}>
      <div className="flex flex-col items-center">
        <span className="w-8 h-8 rounded-full bg-secondary grid place-items-center text-secondary-foreground">
          {marker}
        </span>
        <span className="timeline__connector w-px flex-1 bg-border min-h-8" aria-hidden="true" />
      </div>
      <div className="pb-5">{children}</div>
    </li>
  )
}

export const Timeline = {
  Root: TimelineRoot,
  Item: TimelineItem,
}
