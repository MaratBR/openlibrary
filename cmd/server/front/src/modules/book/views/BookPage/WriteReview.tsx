import './WriteReview.css'
import { Button } from '@/components/ui/button'
import { useTranslation } from 'react-i18next'
import { useState } from 'react'
import {
  BookDetailsDto,
  CreateReviewRequest,
  httpGetMyReview,
  httpUpdateReview,
  ReviewDto,
} from '../../api'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import BookReview from './BookReview'
import ReviewEditor, { ReviewData } from './ReviewEditor'
import { useCurrentUser } from '@/modules/auth/state'
import { pullPreloadedData } from '@/modules/common/api'

export default function WriteReview({ book }: { book: BookDetailsDto }) {
  const { t } = useTranslation()
  const [active, setActive] = useState(false)
  const user = useCurrentUser()

  const updateReview = useUpdateReviewMutation(book.id)

  const { data: myReview } = useQuery({
    queryKey: ['book', book.id, 'reviews', 'my'],
    queryFn: () => httpGetMyReview(book.id),
    initialData: pullPreloadedData<ReviewDto>(`/api/reviews/${book.id}/my`),
    staleTime: 100,
    enabled: !!user,
  })

  async function onUpdated(reviewData: ReviewData, _review: ReviewDto | null) {
    await updateReview.mutateAsync({
      content: reviewData.content,
      rating: reviewData.rating,
    })
    setActive(false)
  }

  return (
    <section className="book-write-review">
      {myReview && (
        <div className="mb-8">
          <h2 className="font-title text-2xl mb-3">{t('book.review.yourReview')}</h2>
          <BookReview bookId={book.id} review={myReview} isAuthor={false} />
        </div>
      )}

      <h2 className="font-title text-2xl">{t('book.reviews')}</h2>

      {!active && !myReview && (
        <Button
          onClick={() => setActive(true)}
          variant="outline3"
          className="rounded-full mt-6 mb-2"
          size="lg"
        >
          {t('book.review.writeReview')}
        </Button>
      )}

      {active && myReview && (
        <div className="book-write-review__form">
          <ReviewEditor review={myReview} onUpdated={onUpdated} onClose={() => setActive(false)} />
        </div>
      )}
    </section>
  )
}

export function useUpdateReviewMutation(bookId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async (request: CreateReviewRequest) => {
      const result = await httpUpdateReview(bookId, request)
      queryClient.setQueryData(['book', bookId, 'reviews', 'my'], result)
      return result
    },
  })
}
