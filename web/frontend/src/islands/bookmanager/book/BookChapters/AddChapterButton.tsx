import { httpCreateChapter } from '@/api/bm'
import { Collapsible } from '@/components/Collapsible'
import { useMutation } from '@tanstack/react-query'
import { useRef, useState } from 'preact/hooks'

export default function AddChapterButton({ bookId }: { bookId: string }) {
  const [open, setOpen] = useState(false)
  const [name, setName] = useState('')
  const inputRef = useRef<HTMLInputElement | null>(null)

  const createChapter = useMutation({
    mutationFn: async () => {
      const normalizedName = name.trim()
      if (name.length === 0 || name.length > 255) return

      const response = await httpCreateChapter(bookId, {
        name: normalizedName,
        summary: '',
        isAdultOverride: false,
        content: '',
      })
      const chapterId = response.data
      window.location.href = `/books-manager/book/${bookId}/chapter/${chapterId}`
    },
  })

  return (
    <div>
      <button
        style={{ display: !open ? undefined : 'none' }}
        class="flex items-center gap-4 hover:bg-highlight w-full transition-colors p-2 rounded-lg mb-4"
        onClick={() => {
          setOpen(true)
          requestAnimationFrame(() => {
            inputRef.current?.focus()
          })
        }}
      >
        <div class="bg-muted text-2xl rounded-full size-16 flex items-center justify-center">
          <i class="fa-solid fa-plus" />
        </div>
        <span class="font-medium text-lg">{window._('bookManager.edit.addChapter')}</span>
      </button>
      <div
        style={{ display: open ? undefined : 'none' }}
        class="shadow-sm rounded-lg border-none relative p-4 mt-1 mb-6"
      >
        <input
          ref={inputRef}
          name="chapterName"
          value={name}
          maxLength={255}
          onChange={(e) => setName((e.target as HTMLInputElement).value)}
          placeholder={window._('bookManager.edit.chapterNamePlaceholder')}
          class="w-full block text-3xl font-medium !outline-none bg-transparent"
        />
        <Collapsible duration={130} in={name.trim().length > 0}>
          <button class="btn btn--lg btn--outline mt-4" onClick={() => createChapter.mutate()}>
            {createChapter.isPending ? (
              <span class="loader loader--dark" />
            ) : (
              window._('bookManager.edit.addChapter')
            )}
          </button>
        </Collapsible>
      </div>
    </div>
  )
}
