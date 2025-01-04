import { useTranslation } from 'react-i18next'
import './BookComments.css'
import VisibilityTrigger from '@/modules/common/components/visibility-trigger'
import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import Spinner from '@/components/spinner'
import { getPreloadedReviews, httpGetReviews } from '../../api'
import BookReview from './BookReview'
import WriteReview from './WriteReview'

export default function BookComments({ bookId, authorId }: { bookId: string; authorId: string }) {
  const { t } = useTranslation()

  const [visible, setVisible] = useState(false)

  const { data, isLoading } = useQuery({
    enabled: visible,
    queryKey: ['book', bookId, 'comments'],
    meta: { disableLoader: true },
    queryFn: () => httpGetReviews(bookId),
    initialData: getPreloadedReviews(bookId),
    staleTime: 100,
  })

  return (
    <section className="book-comments">
      <h2 id="reviews" className="font-title text-2xl">
        {t('book.comments')}
      </h2>

      <WriteReview bookId={bookId} />

      <VisibilityTrigger onVisibilityChange={setVisible} className="mb-10">
        {isLoading && (
          <span className="inline-block pt-4">
            <Spinner thickness={3} size={50} />
          </span>
        )}
        {data &&
          data.reviews.map((review) => (
            <BookReview
              key={review.user.id}
              isAuthor={review.user.id === authorId}
              review={review}
            />
          ))}
      </VisibilityTrigger>
    </section>
  )
}
