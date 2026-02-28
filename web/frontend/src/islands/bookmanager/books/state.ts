import { httpBmGetBooks, httpBmTrashBook, ManagerBookDto } from '@/api/bm/book'
import { create } from 'zustand/react'

type BooksState = {
  books: ManagerBookDto[]
  loading: boolean
  page: number
  totalPages: number

  load(page: number): Promise<void>
  trash(bookId: string, trash: boolean): Promise<void>
}

export const useBooksState = create<BooksState>((set) => ({
  books: [],
  loading: false,
  page: 1,
  totalPages: 0,

  async load(page) {
    set({ loading: true })
    try {
      const response = await httpBmGetBooks({
        page,
        search: '',
        size: 20,
      })
      response.throwIfError()
      set({ page, books: response.data.books, totalPages: response.data.totalPages })
    } finally {
      set({ loading: false })
    }
  },

  async trash(bookId, trash) {
    const response = await httpBmTrashBook({
      id: bookId,
      trash,
    })
    response.throwIfError()
  },
}))
