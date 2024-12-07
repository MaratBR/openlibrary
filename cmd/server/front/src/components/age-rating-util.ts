import { AgeRating } from '@/modules/book/api/api'

type AgeRatingInfo = {
  title: string
  summary: string
}

export function useAgeRatingsInfo(): Record<AgeRating, AgeRatingInfo> {
  return RATINGS
}

const RATINGS: Record<AgeRating, { title: string; summary: string }> = {
  '?': {
    title: 'Unknown rating',
    summary: 'The author did not specify the rating of this book, anything is possible.',
  },
  G: { title: 'G', summary: 'Suitable for all ages.' },
  PG: {
    title: 'Parental Guidance Suggested',
    summary:
      'Some material may not be suitable for children. Parents urged to give "parental guidance". May contain some material parents might not like for their young children.',
  },
  'PG-13': {
    title: 'PG-13',
    summary:
      'Some material may be inappropriate for children under 13. Parents are urged to be cautious. Some material may be inappropriate for pre-teenagers.',
  },
  R: {
    title: 'R',
    summary:
      'Contains some adult material. Parents are urged to learn more about the film before taking their young children with them.',
  },
  'NC-17': {
    title: 'NC-17',
    summary:
      'Contains adult content not suitable for children or teens. This may include: explicit sexual content, excessive violence.',
  },
}
