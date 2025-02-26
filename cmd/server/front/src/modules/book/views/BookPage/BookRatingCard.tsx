import StarRating from '@/components/star-rating'
import './BookRatingCard.css'
import { useTranslation } from 'react-i18next'

export default function BookRatingCard({
  bookId,
  rating,
  votes,
  reviews,
}: {
  bookId: string
  rating: number | null
  votes: number
  reviews: number
}) {
  const { t } = useTranslation()

  return (
    <a href="#reviews" className="book-rating-card focus-outline">
      {rating ? (
        <>
          <StarRating value={rating} />

          <div className="book-rating-card__info">
            {votes} votes <br />
            {reviews} reviews
          </div>
        </>
      ) : (
        <>
          <StarRating value={0} />
          <p className="ml-2">{t('book.noRatingsYet')}</p>
        </>
      )}
    </a>
  )
}
