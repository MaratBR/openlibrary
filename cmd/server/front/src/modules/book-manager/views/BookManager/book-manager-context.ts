import React, { startTransition } from 'react'
import { httpUpdateBook, ManagerBookDetailsDto, UpdateBookRequest } from '../../api'
import { useMutation, useQueryClient } from '@tanstack/react-query'

export type BookManagerContext = {
  book: ManagerBookDetailsDto
}

export const BookManagerContext = React.createContext<BookManagerContext | null>(null)

export function useBookManager() {
  const ctx = React.useContext(BookManagerContext)
  if (ctx === null) throw new Error('useBookManager must be used within a BookManager')
  return ctx
}

export function useBookManagerUpdateMutation() {
  const { book } = useBookManager()
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (req: UpdateBookRequest) => {
      return httpUpdateBook(book.id!, req)
        .then((r) => r.book)
        .then((book) => {
          startTransition(() => {
            queryClient.setQueryData(['manager', 'books', book.id!], book)
          })
        })
    },
  })
}
