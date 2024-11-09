import { useBookManager, useBookManagerChaptersQuery } from './book-manager-context'
import { GripVertical, Save } from 'lucide-react'
import Spinner from '@/components/spinner'
import BookChapterCard from './BookChapterCard'
import React from 'react'
import { httpReorderChapters, ManagerBookChapterDto, managerBookChapterDtoSchema } from '../api'
import {
  closestCenter,
  DndContext,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core'
import { CSS } from '@dnd-kit/utilities'
import { cn } from '@/lib/utils'
import {
  SortableContext,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable'
import BackToBookButton from './BackToBookButton'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { useMutation } from '@tanstack/react-query'

export default function BookChaptersReorder() {
  // const { book } = useBookManager()
  const { data, isLoading, refetch } = useBookManagerChaptersQuery()

  function handleReorder() {
    refetch()
  }

  return (
    <section className="page-section">
      <header className="section-header">
        <BackToBookButton />
        <h1 className="section-header-text">Reorder chapters</h1>
      </header>

      {isLoading && <Spinner />}

      {data && (
        <>
          {data.length === 0 && (
            <p className="my-4 text-muted-foreground">No chapters yet. Nothing to reorder.</p>
          )}
          <BookChapterListSortable chapters={data} onReorder={handleReorder} />
        </>
      )}
    </section>
  )
}

function BookChapterListSortable({
  chapters,
  onReorder,
}: {
  chapters: ManagerBookChapterDto[]
  onReorder: () => void
}) {
  const { book } = useBookManager()
  const [draggingChapter, setDraggingChapter] = React.useState<ManagerBookChapterDto | null>(null)
  const [chaptersOrder, setChaptersOrder] = React.useState<ManagerBookChapterDto[]>([])

  const changedChapters = chaptersOrder.reduce(
    (acc, chapter, index) => acc + (chapter.order !== index + 1 ? 1 : 0),
    0,
  )

  const reorderMutation = useMutation({
    mutationFn: async () => {
      await httpReorderChapters(
        book.id,
        chaptersOrder.map((x) => x.id),
      )
      onReorder()
    },
  })

  React.useEffect(() => {
    setChaptersOrder(chapters)
  }, [chapters])

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  )

  function handleDragStart(event: DragStartEvent) {
    const chapterAny = event.active.data.current?.chapter

    if (chapterAny) {
      try {
        const chapter = managerBookChapterDtoSchema.parse(chapterAny)
        setDraggingChapter(chapter)
      } catch (e) {
        console.log(e, chapterAny)
      }
    }
  }

  function handleDragEnd(event: DragEndEvent) {
    setDraggingChapter(null)
    if (!event.over) return

    const id = event.active.id as string
    const newIndex = chaptersOrder.findIndex((x) => x.id === event.over!.id)
    if (newIndex === -1) return
    const oldIndex = chaptersOrder.findIndex((x) => x.id === id)
    if (oldIndex === -1) return
    if (oldIndex === newIndex) return
    setChaptersOrder((chapters) => {
      const copy = [...chapters]
      copy.splice(newIndex, 0, copy.splice(oldIndex, 1)[0])
      return copy
    })
  }

  return (
    <>
      <div className="mb-4">
        <div className="mb-3">
          <Button
            disabled={changedChapters === 0 || reorderMutation.isPending}
            onClick={() => reorderMutation.mutate()}
          >
            <Save /> Save {changedChapters} changes
          </Button>
        </div>

        {changedChapters === 0 && (
          <p className="text-muted-foreground">
            Change order fo the chapters below by dragging them into their new positions.
          </p>
        )}
        {changedChapters !== 0 && (
          <p className="text-muted-foreground">
            You've updated the order of {changedChapters} chapters. "Grip" icon of chapters with
            changed position is highlighted in blue.
          </p>
        )}

        <Separator className="my-4" />
      </div>
      <DndContext
        sensors={sensors}
        collisionDetection={closestCenter}
        onDragStart={handleDragStart}
        onDragEnd={handleDragEnd}
      >
        <SortableContext
          strategy={verticalListSortingStrategy}
          items={chaptersOrder.map((c) => c.id)}
        >
          <div className="space-y-2">
            {chaptersOrder.map((chapter, idx) => (
              <BookChapterCardDraggable key={chapter.id} order={idx + 1} chapter={chapter} />
            ))}
          </div>
          <DragOverlay>
            {draggingChapter && (
              <div className="grid grid-cols-[40px_1fr]">
                <div></div>
                <BookChapterCard chapter={draggingChapter} />
              </div>
            )}
          </DragOverlay>
        </SortableContext>
      </DndContext>
    </>
  )
}

function BookChapterCardDraggable({
  chapter,
  order,
}: {
  chapter: ManagerBookChapterDto
  order: number
}) {
  const { attributes, isDragging, transform, transition, listeners, setNodeRef } = useSortable({
    id: chapter.id,
    data: { chapter },
  })

  const style: React.CSSProperties = {
    transform: CSS.Transform.toString(transform),
    transition,
    visibility: isDragging ? 'hidden' : 'visible',
  }

  return (
    <div ref={setNodeRef} {...attributes} style={style} className="grid grid-cols-[40px_1fr]">
      <div {...listeners} className="mr-1 flex flex-col items-center justify-around">
        <div className="py-2 px-1 hover:bg-muted rounded-sm">
          <GripVertical
            className={cn({
              'text-primary': order !== chapter.order,
            })}
          />
        </div>
      </div>
      <BookChapterCard chapter={chapter} />
    </div>
  )
}
