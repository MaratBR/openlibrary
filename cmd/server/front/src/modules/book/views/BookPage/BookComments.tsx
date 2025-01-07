import { useTranslation } from 'react-i18next'
import './BookComments.css'
import VisibilityTrigger from '@/modules/common/components/visibility-trigger'
import { useMemo, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import Spinner from '@/components/spinner'
import { BookDetailsDto, getPreloadedReviews, httpGetReviews } from '../../api'
import BookReview from './BookReview'
import WriteReview from './WriteReview'
import { useCurrentUser } from '@/modules/auth/state'

export default function BookComments({ book }: { book: BookDetailsDto }) {
  const { t } = useTranslation()

  const user = useCurrentUser()

  const [visible, setVisible] = useState(false)

  const { data, isLoading } = useQuery({
    enabled: visible,
    queryKey: ['book', book.id, 'comments'],
    meta: { disableLoader: true },
    queryFn: () => httpGetReviews(book.id),
    initialData: getPreloadedReviews(book.id),
    staleTime: 100,
  })

  const reviews = useMemo(() => {
    if (!data) return []

    if (user) {
      return data.reviews.filter((x) => x.user.id !== user.id)
    } else {
      return data.reviews
    }
  }, [data, user])

  return (
    <section className="book-comments">
      <div aria-hidden id="reviews" className="relative pointer-events-none bottom-[70px]"></div>

      <WriteReview book={book} />

      <VisibilityTrigger onVisibilityChange={setVisible} className="mb-10">
        {isLoading && (
          <span className="inline-block pt-4">
            <Spinner thickness={3} size={50} />
          </span>
        )}
        <div className="space-y-3">
          {reviews.map((review) => (
            <BookReview
              key={review.user.id}
              isAuthor={review.user.id === book.author.id}
              review={review}
              bookId={book.id}
            />
          ))}
        </div>
      </VisibilityTrigger>
    </section>
  )
}
