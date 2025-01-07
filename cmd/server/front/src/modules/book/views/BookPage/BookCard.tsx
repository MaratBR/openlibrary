import BookCover from '@/modules/common/components/book-cover'
import {
  BookDetailsDto,
  bookDetailsDtoSchema,
  ReadingListStatus,
  readingListStatusSchema,
} from '../../api'
import { Button } from '@/components/ui/button'
import { Pencil1Icon } from '@radix-ui/react-icons'
import { useTranslation } from 'react-i18next'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { httpUpdateReadingListStatus } from '../../api/reading-list'
import { ButtonSpinner } from '@/components/spinner'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'

export default function BookCard({ book }: { book: BookDetailsDto }) {
  return (
    <div className="flex flex-col items-center gap-4 md:sticky md:top-24">
      <BookCover name={book.name} url={book.cover} />
      <ReadingListStatusSection readingList={book.readingList} bookId={book.id} />
    </div>
  )
}

function ReadingListStatusSection({
  bookId,
  readingList,
}: {
  bookId: string
  readingList: BookDetailsDto['readingList']
}) {
  const { t } = useTranslation()
  const queryClient = useQueryClient()

  const updateState = useMutation({
    mutationFn: async (status: ReadingListStatus) => {
      const readingList = await httpUpdateReadingListStatus(bookId, status)
      const bookDetails = queryClient.getQueryData(['book', bookId])
      const parseResult = bookDetailsDtoSchema.safeParse(bookDetails)
      if (parseResult.success) {
        queryClient.setQueryData(['book', bookId], {
          ...parseResult.data,
          readingList,
        })
      }
    },
  })

  if (!readingList) {
    return (
      <Button
        onClick={() => updateState.mutate('want_to_read')}
        className="rounded-full text-md w-[75%] h-[3em]"
        size="lg"
      >
        {updateState.isPending && <ButtonSpinner />}
        Want to read
      </Button>
    )
  }

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="outline3" className="rounded-full text-md w-[75%] h-[3em]" size="lg">
          {updateState.isPending ? <ButtonSpinner /> : <Pencil1Icon />}
          {t(`readingList.status.${readingList.status}`)}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="rounded-2xl !animate-none shadow-none">
        <div className="flex flex-col items-stretch gap-2">
          {readingListStatusSchema.options
            .filter((x) => x !== readingList.status)
            .map((option) => (
              <Button
                key={option}
                disabled={updateState.isPending}
                className="rounded-full"
                variant="outline"
                onClick={() => updateState.mutate(option)}
              >
                {t(`readingList.status.${option}`)}
              </Button>
            ))}
        </div>
      </PopoverContent>
    </Popover>
  )
}
