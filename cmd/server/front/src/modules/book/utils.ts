import { AgeRating } from './api'

export function useChapterName({ name, order }: { name: string; order: number }): string {
  let finalName = `Chapter ${order}`
  if (name) {
    finalName += ': ' + name
  }
  return finalName
}

export function isAgeRatingAdult(ageRating: AgeRating): boolean {
  return ageRating === 'R' || ageRating === 'NC-17'
}
