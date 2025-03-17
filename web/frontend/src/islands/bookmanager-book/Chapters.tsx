import { useMemo } from 'preact/hooks'
import { PreactIslandProps } from '../common'
import { managerBookDetailsSchema } from './api'

interface Chapter {
  id: string
  name: string
  words: number
  order: number
  createdAt: string
  summary: string
}

export default function Chapters({ data: dataUnknown }: PreactIslandProps) {
  const data = useMemo(() => managerBookDetailsSchema.parse(dataUnknown), [dataUnknown])

  return (
    <div class="chapters-container">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-bold">{window._('bookManager.edit.chapters')}</h2>
        <button class="ol-btn ol-btn--primary">{window._('bookManager.edit.addChapter')}</button>
      </div>

      {data.chapters?.length ? (
        <div class="chapters-list">
          {data.chapters.map((chapter: Chapter, index: number) => (
            <div key={chapter.id} class="chapter-item ol-card p-4 mb-2">
              <div class="flex justify-between items-center">
                <div>
                  <h3 class="font-medium">{chapter.name}</h3>
                  <p class="text-sm text-gray-600">{chapter.summary}</p>
                  <div class="text-xs text-gray-500 mt-1">
                    {window._('bookManager.edit.words')}: {chapter.words}
                  </div>
                </div>
                <div class="flex gap-2">
                  <button class="ol-btn ol-btn--secondary">
                    {window._('bookManager.edit.edit')}
                  </button>
                  <button class="ol-btn ol-btn--danger">
                    {window._('bookManager.edit.delete')}
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div class="text-center py-8 text-gray-500">{window._('bookManager.edit.noChapters')}</div>
      )}
    </div>
  )
}
