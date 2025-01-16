import { Button } from '@/components/ui/button'
import { BookDetailsDto, bookDetailsDtoSchema } from '../../api'
import './StartReading.css'
import { NavLink } from 'react-router-dom'
import { Trans, useTranslation } from 'react-i18next'
import { updateReadingListStartReading } from '../../api/reading-list'
import { useMemo } from 'react'
import { ArrowRight } from 'lucide-react'
import { Separator } from '@/components/ui/separator'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { DateTime } from 'luxon'

export default function StartReading({ book }: { book: BookDetailsDto }) {
  const { t } = useTranslation()

  const queryClient = useQueryClient()

  const startReading = useMutation({
    mutationFn: async (chapterId: string) => {
      const data = queryClient.getQueryData(['book', book.id])
      if (data) {
        const parseResult = bookDetailsDtoSchema.safeParse(data)
        if (parseResult.success) {
          const newData: BookDetailsDto = {
            ...parseResult.data,
            readingList: {
              status: 'reading',
              chapterId,
              lastUpdatedAt: DateTime.now().toISO(),
            },
          }
          queryClient.setQueryData(['book', book.id], newData)
        }
      }
      await updateReadingListStartReading(book.id, chapterId)
    },
  })

  const memoizedChildren = useMemo(() => {
    if (book.chapters.length === 0) return undefined

    if (!book.readingList || book.readingList.status === 'want_to_read') {
      return (
        <>
          <NavLink to={`/book/${book.id}/chapters/${book.chapters[0].id}#chapter`}>
            <Button
              onClick={() => startReading.mutate(book.chapters[0].id)}
              size="lg"
              variant="outline3"
              className="rounded-full text-md"
            >
              {t('book.startFromFirstChapter')}
            </Button>
          </NavLink>

          {!book.readingList && (
            <Button variant="ghost" size="lg" className="ml-2 pl-2 rounded-full text-md">
              {t('book.saveForLater')}
            </Button>
          )}
        </>
      )
    } else if (book.readingList.status === 'reading' || book.readingList.status === 'paused') {
      if (book.readingList.chapterId) {
        const chapter = book.chapters.find((x) => x.id === book.readingList!.chapterId)

        return (
          <NavLink
            to={`/book/${book.id}/chapters/${book.readingList.chapterId}${book.readingList.status === 'paused' ? '?event=reading_list_unpause' : ''}#chapter`}
          >
            <Button size="lg" variant="outline3" className="rounded-full text-md">
              {chapter
                ? t('book.continueReading', { chapter: chapter?.name })
                : t('book.continueReadingWereYouLeftOff')}{' '}
              <ArrowRight />
            </Button>
          </NavLink>
        )
      } else {
        return (
          <NavLink to={`/book/${book.id}/chapters/${book.chapters[0].id}#chapter`}>
            <Button
              onClick={() => updateReadingListStartReading(book.id, book.chapters[0].id)}
              size="lg"
              variant="outline"
              className="rounded-full text-md"
            >
              {t('book.startFromFirstChapter')}
            </Button>
          </NavLink>
        )
      }
    } else if (book.readingList.status === 'read') {
      return (
        <>
          <p>
            <Trans
              i18nKey="book.finishedReadingBook"
              values={{ bookName: book.name }}
              components={{ name: <em /> }}
            />
          </p>
          <NavLink to={`/book/${book.id}/chapters/${book.chapters[0].id}?event=rereading#chapter`}>
            <Button size="lg" variant="outline3" className="mt-4 rounded-full text-md">
              {t('book.startFromFirstChapterAgain')}
            </Button>
          </NavLink>
        </>
      )
    }
  }, [book, t])

  if (book.chapters.length === 0) {
    return null
  }

  return (
    <>
      <Separator className="my-4" />
      <div className="book-reading-list">{memoizedChildren}</div>
    </>
  )
}
