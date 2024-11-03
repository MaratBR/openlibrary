import { Navigate, useParams } from 'react-router'
import ChapterEditor from './ChapterEditor'

export function CreateChapterPage() {
  const { bookId } = useParams<{ bookId: string }>()

  if (!bookId) return <Navigate to="/home" />

  return <ChapterEditor bookId={bookId} chapterId={null} />
}
