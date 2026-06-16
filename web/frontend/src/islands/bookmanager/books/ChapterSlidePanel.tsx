import { ManagerBookChapterDto } from '@/api/bm'
import { SlidePanel } from '@/components/SlidePanel'

export const CHAPTER_SLIDE_OUT_PARAMETER_NAME = 'chapter'

export function ChapterSlidePanel({
  chapter,
  onClose,
}: {
  chapter: ManagerBookChapterDto | null
  onClose: () => void
}) {
  return (
    <SlidePanel.Facade open={!!chapter} onClose={onClose}>
      <SlidePanel.Content>{chapter?.name}</SlidePanel.Content>
    </SlidePanel.Facade>
  )
}
