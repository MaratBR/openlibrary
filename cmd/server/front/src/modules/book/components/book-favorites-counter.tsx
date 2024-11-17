import FavoritesCounter from '@/components/favorites-counter'
import React from 'react'
import { httpFavoriteBook } from '../api'

export type BookFavoritesCounterProps = {
  count: number
  isLiked: boolean
  bookId: string
}

export default function BookFavoritesCounter({
  count,
  isLiked,
  bookId,
}: BookFavoritesCounterProps) {
  const [isLikedState, setIsLikedState] = React.useState(isLiked)
  const [countWithoutSelf, _setCountWithoutSelf] = React.useState(count + (isLiked ? -1 : 0))

  function handleToggle() {
    setIsLikedState(!isLikedState)

    httpFavoriteBook(bookId, !isLikedState)
  }

  return (
    <FavoritesCounter
      count={countWithoutSelf + (isLikedState ? 1 : 0)}
      isLiked={isLikedState}
      onClick={handleToggle}
    />
  )
}
