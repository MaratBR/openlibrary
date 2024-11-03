import { cn } from '@/lib/utils'
import React from 'react'

export type AdultIndicatorProps = React.HTMLAttributes<HTMLDivElement>

const AdultIndicator = React.forwardRef(
  ({ className, ...props }: AdultIndicatorProps, ref: React.ForwardedRef<HTMLDivElement>) => {
    return (
      <div
        ref={ref}
        className={cn(
          'font-bold px-[0.25em] py-[0.125em] bg-red-800 text-white rounded-sm text-[1.2rem] h-7 inline-flex items-center justify-center',
          className,
        )}
        {...props}
      >
        Adult
      </div>
    )
  },
)

export default AdultIndicator
