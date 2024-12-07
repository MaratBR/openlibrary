import { LabeledValue, LabeledValueLabel, LabeledValueLayout } from '@/components/labeled-value'
import { Card } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import WordsCount from '@/components/words-count'
import Tag from '../Tag'
import AgeRatingBadge from '@/components/age-rating-badge'
import { BookDetailsDto } from '../../api/api'
import React from 'react'

const BookInfoCard = React.memo(
  ({ book, className }: { book: BookDetailsDto; className?: string }) => {
    return (
      <Card className={className}>
        <div className="space-y-2 py-3">
          <LabeledValueLayout className="px-4 py-1">
            <LabeledValueLabel>Age rating</LabeledValueLabel>
            <LabeledValue>
              <AgeRatingBadge value={book.ageRating} />
            </LabeledValue>
          </LabeledValueLayout>
          <Separator />
          <LabeledValueLayout className="px-4 py-1">
            <LabeledValueLabel>Tags</LabeledValueLabel>
            <LabeledValue>
              <div className="flex flex-wrap gap-2">
                {book.tags.map((tag) => {
                  return <Tag key={tag.id} tag={tag} />
                })}
              </div>
            </LabeledValue>
          </LabeledValueLayout>
          <Separator />
          <LabeledValueLayout className="px-4 py-1">
            <LabeledValueLabel>Words in total</LabeledValueLabel>
            <LabeledValue>
              <WordsCount value={book.words} />
            </LabeledValue>
          </LabeledValueLayout>
          <LabeledValueLayout className="px-4 py-1">
            <LabeledValueLabel>Words per chapter (average)</LabeledValueLabel>
            <LabeledValue>
              <WordsCount value={book.wordsPerChapter} />
            </LabeledValue>
          </LabeledValueLayout>
          <Separator />
          <LabeledValueLayout className="px-4 py-1">
            <LabeledValueLabel>Collections</LabeledValueLabel>
            <LabeledValue>
              <pre>{JSON.stringify(book.collections, null, 2)}</pre>
            </LabeledValue>
          </LabeledValueLayout>
        </div>
      </Card>
    )
  },
)

export default BookInfoCard
