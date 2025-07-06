import { Fragment } from 'preact/jsx-runtime'

export type AgeRatingProps = {
  name?: string
  value: string | null

  onChange: (rating: string) => void
}

export default function AgeRatingInput({ value, onChange, name }: AgeRatingProps) {
  return (
    <div class="flex gap-2 flex-wrap">
      {window.__server__.ageRatings.map((ageRating) => {
        const id = `new-book-${ageRating}`
        return (
          <Fragment key={ageRating}>
            <input
              key={ageRating}
              id={id}
              class="age-rating-input"
              name={name}
              value={ageRating}
              type="radio"
              checked={ageRating === value}
              onChange={() => onChange(ageRating)}
            />
            <label data-rating={ageRating} class="age-rating" for={id}>
              {ageRating}
            </label>
          </Fragment>
        )
      })}
    </div>
  )
}
