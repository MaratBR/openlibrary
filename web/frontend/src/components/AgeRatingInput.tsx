import { AgeRating } from '@/api/common'
import { Fragment } from 'react'

export type AgeRatingProps = {
  name?: string
  value: AgeRating | null

  onChange: (rating: string) => void
}

export default function AgeRatingInput({ value, onChange, name }: AgeRatingProps) {
  return (
    <div className="flex gap-2 flex-wrap">
      {window.__server__.ageRatings.map((ageRating) => {
        const id = `new-book-${ageRating}`
        return (
          <Fragment key={ageRating}>
            <input
              key={ageRating}
              id={id}
              className="age-rating-input"
              name={name}
              value={ageRating}
              type="radio"
              checked={ageRating === value}
              onChange={() => onChange(ageRating)}
            />
            <label data-rating={ageRating} className="age-rating" htmlFor={id}>
              {ageRating}
            </label>
          </Fragment>
        )
      })}
    </div>
  )
}
