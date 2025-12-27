import { httpGetBookChapters, httpUpdateChaptersOrder, ManagerBookChapterDto } from '@/api/bm'
import { create } from 'zustand/react'

export const useBookChaptersState = create<{
  chapters: ManagerBookChapterDto[]
  loading: boolean
  error: unknown
  bookId: string

  // chapter that is being reordered
  reorderActiveChapter: ManagerBookChapterDto | null
  setReorderActiveChapter(chapter: ManagerBookChapterDto): void

  loadChapters(bookId: string): void
  swapChapters(chapterId1: string, chapterId2: string): Promise<void>
  insertBefore(chapterId1: string, chapterId2: string): Promise<void>
  insertAfter(chapterId1: string, chapterId2: string): Promise<void>
}>((set, get) => ({
  chapters: [],
  loading: false,
  error: undefined,
  bookId: '',

  reorderActiveChapter: null,
  setReorderActiveChapter(chapter: ManagerBookChapterDto): void {
    set({ reorderActiveChapter: chapter })
  },

  loadChapters(bookId) {
    if (get().loading) return
    set({ bookId, loading: true })

    httpGetBookChapters(bookId)
      .then((response) => {
        set({ chapters: response.data, loading: false })
      })
      .catch((err) => {
        set({
          error: err,
          loading: false,
        })
      })
  },

  async swapChapters(chapterId1: string, chapterId2: string): Promise<void> {
    const { bookId, chapters } = get()

    const chapterIdx1 = chapters.findIndex((x) => x.id === chapterId1)
    const chapterIdx2 = chapters.findIndex((x) => x.id === chapterId2)
    if (chapterIdx1 === -1)
      throw new Error(`cannot find chapter ${chapterId1} in a list of chapters`)
    if (chapterIdx2 === -1)
      throw new Error(`cannot find chapter ${chapterId2} in a list of chapters`)

    await httpUpdateChaptersOrder(bookId, {
      modifications: [
        {
          chapterId: chapterId1,
          newIndex: chapterIdx2,
        },
        {
          chapterId: chapterId2,
          newIndex: chapterIdx1,
        },
      ],
    })
    await this.loadChapters(bookId)
  },
  async insertBefore(chapterId1: string, chapterId2: string): Promise<void> {
    const { bookId, chapters } = get()

    const chapterIdx1 = chapters.findIndex((x) => x.id === chapterId1)
    const chapterIdx2 = chapters.findIndex((x) => x.id === chapterId2)
    if (chapterIdx1 === -1)
      throw new Error(`cannot find chapter ${chapterId1} in a list of chapters`)
    if (chapterIdx2 === -1)
      throw new Error(`cannot find chapter ${chapterId2} in a list of chapters`)

    const newIndex = chapterIdx1 < chapterIdx2 ? chapterIdx2 - 1 : chapterIdx2

    await httpUpdateChaptersOrder(bookId, {
      modifications: [
        {
          chapterId: chapterId1,
          newIndex,
        },
      ],
    })
    await this.loadChapters(bookId)
  },
  async insertAfter(chapterId1: string, chapterId2: string): Promise<void> {
    const { bookId, chapters } = get()

    const chapterIdx1 = chapters.findIndex((x) => x.id === chapterId1)
    const chapterIdx2 = chapters.findIndex((x) => x.id === chapterId2)
    if (chapterIdx1 === -1)
      throw new Error(`cannot find chapter ${chapterId1} in a list of chapters`)
    if (chapterIdx2 === -1)
      throw new Error(`cannot find chapter ${chapterId2} in a list of chapters`)

    const newIndex = chapterIdx2

    await httpUpdateChaptersOrder(bookId, {
      modifications: [
        {
          chapterId: chapterId1,
          newIndex,
        },
      ],
    })
    await this.loadChapters(bookId)
  },
}))
