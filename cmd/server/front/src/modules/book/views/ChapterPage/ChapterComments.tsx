import VisibilityTrigger from '@/modules/common/components/visibility-trigger'
import { useQuery } from '@tanstack/react-query'

export default function ChapterComments({ chapterId }: { chapterId: string }) {
  const { data, isLoading, error } = useQuery({
    queryKey: ['comments', 'chapter', chapterId],
  })

  function handleAppear() {}

  return (
    <VisibilityTrigger onAppear={handleAppear} className="chapter-comments container-default">
      Coments will be here
    </VisibilityTrigger>
  )
}
