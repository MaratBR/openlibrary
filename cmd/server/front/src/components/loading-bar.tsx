import { cn } from '@/lib/utils'
import React from 'react'
import './loading-bar.css'

const LoadingBar = React.forwardRef(
  (
    { className, ...props }: React.HTMLAttributes<HTMLDivElement>,
    ref: React.ForwardedRef<HTMLDivElement>,
  ) => {
    return (
      <div
        ref={ref}
        style={{ '--loading-bar-duration': '1s' } as React.CSSProperties}
        className={cn('loading-bar h-1 w-full bg-primary/20', className)}
        {...props}
      >
        <div className="loading-bar__inner bg-primary h-full"></div>
      </div>
    )
  },
)

export default LoadingBar
