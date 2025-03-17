import { useMemo, useState } from 'preact/hooks'
import { PreactIslandProps } from '../common'
import { managerBookDetailsSchema } from './api'
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  DragEndEvent,
} from '@dnd-kit/core'
import {
  arrayMove,
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'

interface Chapter {
  id: string
  name: string
  words: number
  order: number
  createdAt: string
  summary: string
}

function SortableChapterItem({
  chapter,
  onEdit,
  isReordering,
}: {
  chapter: Chapter
  onEdit: () => void
  isReordering: boolean
}) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: chapter.id,
    disabled: !isReordering,
  })

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      class={`chapter-item ol-card p-4 mb-2 ${
        isReordering ? (isDragging ? 'cursor-grabbing' : 'cursor-grab') : ''
      } ${isReordering ? 'bg-secondary/50 dark:bg-secondary/20' : ''}`}
      {...(isReordering ? { ...attributes, ...listeners } : {})}
    >
      <div class="flex justify-between items-center">
        <div>
          <h3 class="font-medium">{chapter.name}</h3>
          <p class="text-sm text-gray-600">{chapter.summary}</p>
          <div class="text-xs text-gray-500 mt-1">
            {window._('bookManager.edit.words')}: {chapter.words}
          </div>
        </div>
        <div class="flex gap-2">
          <button class="ol-btn ol-btn--secondary" onClick={onEdit}>
            {window._('bookManager.edit.edit')}
          </button>
          <button class="ol-btn ol-btn--danger">{window._('bookManager.edit.delete')}</button>
        </div>
      </div>
    </div>
  )
}

export default function Chapters({ data: dataUnknown }: PreactIslandProps) {
  const data = useMemo(() => managerBookDetailsSchema.parse(dataUnknown), [dataUnknown])
  const [chapters, setChapters] = useState(data.chapters || [])
  const [isReordering, setIsReordering] = useState(false)

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  )

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event

    if (over && active.id !== over.id) {
      setChapters((items) => {
        const oldIndex = items.findIndex((item) => item.id === active.id)
        const newIndex = items.findIndex((item) => item.id === over.id)
        return arrayMove(items, oldIndex, newIndex)
      })
    }
  }

  const handleSaveOrder = () => {
    // TODO: Implement API call to save new chapter order
    setIsReordering(false)
  }

  return (
    <div class="chapters-container">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-semibold">{window._('bookManager.edit.chapters')}</h2>
        <div class="flex gap-2">
          {!isReordering ? (
            <>
              <button class="ol-btn ol-btn--secondary" onClick={() => setIsReordering(true)}>
                {window._('bookManager.edit.reorder')}
              </button>
              <button class="ol-btn ol-btn--primary">
                {window._('bookManager.edit.addChapter')}
              </button>
            </>
          ) : (
            <>
              <button class="ol-btn ol-btn--secondary" onClick={() => setIsReordering(false)}>
                {window._('bookManager.edit.cancel')}
              </button>
              <button class="ol-btn ol-btn--primary" onClick={handleSaveOrder}>
                {window._('bookManager.edit.save')}
              </button>
            </>
          )}
        </div>
      </div>

      {chapters?.length ? (
        <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
          <SortableContext items={chapters.map((c) => c.id)} strategy={verticalListSortingStrategy}>
            <div class="chapters-list">
              {chapters.map((chapter: Chapter) => (
                <SortableChapterItem
                  key={chapter.id}
                  chapter={chapter}
                  isReordering={isReordering}
                  onEdit={() =>
                    (window.location.href = `/books-manager/book/${data.id}/chapter/${chapter.id}`)
                  }
                />
              ))}
            </div>
          </SortableContext>
        </DndContext>
      ) : (
        <div class="text-center py-8 text-gray-500">{window._('bookManager.edit.noChapters')}</div>
      )}
    </div>
  )
}
