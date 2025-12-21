import { httpGetBookChapters } from '@/api/bm'
import { PreactIslandProps } from '@/lib/island'
import { useQuery } from '@tanstack/react-query'
import { useMemo } from 'preact/hooks'
import z from 'zod'

const schema = z.object({
  bookId: z.string(),
})

export function BookChapters({ data: dataUnknown }: PreactIslandProps) {
  const data = useMemo(() => schema.parse(dataUnknown), [dataUnknown])

  const { data: chapters } = useQuery({
    queryKey: ['chapters', data.bookId],
    queryFn: () => httpGetBookChapters(data.bookId).then((r) => r.data),
  })

  return (
    <div>
      123
      {chapters?.map((chapter) => {
        return <div key={chapter.id}>{chapter.id}</div>
      })}
    </div>
  )
}
