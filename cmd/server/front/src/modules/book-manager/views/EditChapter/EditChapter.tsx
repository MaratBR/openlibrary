import { useParams } from 'react-router'
import ChapterEditor from './ChapterEditor'
import { UrlError } from '@/components/errors'

export function EditChapterPage() {
  const { chapterId } = useParams<{ chapterId: string }>()

  if (typeof chapterId !== 'string' || !chapterId) {
    return (
      <UrlError>
        <code>chapterId</code> must be a string
      </UrlError>
    )
  }

  return <ChapterEditor chapterId={chapterId} />
}
