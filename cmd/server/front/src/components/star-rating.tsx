import React from 'react'
import './star-rating.css'
import { cn } from '@/lib/utils'

export type StarRatingProps = React.HTMLAttributes<HTMLDivElement> & {
  value: number
  size?: string | number
}

export default function StarRating({
  value,
  size = '2.5em',
  className,
  style,
  ...props
}: StarRatingProps) {
  return (
    <div
      style={
        {
          '--star-size': typeof size === 'number' ? `${size}px` : size,
          ...style,
        } as React.CSSProperties
      }
      className={cn('star-rating', className)}
      {...props}
    >
      <FilledStar fillPercentage={Math.min(value, 1) * 100} />
      <FilledStar fillPercentage={Math.min(value - 1, 1) * 100} />
      <FilledStar fillPercentage={Math.min(value - 2, 1) * 100} />
      <FilledStar fillPercentage={Math.min(value - 3, 1) * 100} />
      <FilledStar fillPercentage={Math.min(value - 4, 1) * 100} />
    </div>
  )
}

const Star: React.FC<{ className?: string }> = ({ className }) => (
  <svg className={className} viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
    <polygon
      points="50,5 35,40 0,40 30,60 20,95 50,75 80,95 70,60 100,40 65,40"
      fill="currentColor"
      stroke="currentColor"
      stroke-width="2"
    />
  </svg>
)

const FilledStar = ({ fillPercentage = 0 }) => {
  const percentage = Math.max(0, Math.min(100, fillPercentage))

  return (
    <div className="relative inline-block m-[calc(var(--star-size)*0.1)]">
      {/* Background star (gray) */}
      <Star className="text-gray-300 w-[var(--star-size)] dark:text-gray-700" />

      {/* Filled star (yellow) with clip path */}
      <div
        className="absolute top-0 left-0 text-yellow-400 overflow-hidden"
        style={{ width: `${percentage}%` }}
      >
        <Star className="w-[var(--star-size)]" />
      </div>
    </div>
  )
}
