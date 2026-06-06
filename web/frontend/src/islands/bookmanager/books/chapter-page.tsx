import { SlidePanel } from '@/components/SlidePanel'
import { useSearchParams } from 'react-router'

export const CHAPTER_SLIDE_OUT_PARAMETER_NAME = 'chapter'

export function ChapterSlidePanel() {
  const [searchParams, setSearchParams] = useSearchParams()

  const chapterId = searchParams.get(CHAPTER_SLIDE_OUT_PARAMETER_NAME)

  return (
    <SlidePanel.Facade
      open={!!chapterId}
      onClose={() => {
        setSearchParams((sp) => {
          sp.delete(CHAPTER_SLIDE_OUT_PARAMETER_NAME)
          return sp
        })
      }}
    >
      <SlidePanel.Content>{chapterId}</SlidePanel.Content>
    </SlidePanel.Facade>
  )
}
