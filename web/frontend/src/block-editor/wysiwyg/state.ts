import { create } from 'zustand/react'
import { ChapterContentEditor } from './editor'

export type WYSIWYGState = {
  editor: ChapterContentEditor
}

export const useWYSIWYG = create<WYSIWYGState>((get) => ({
  editor: new ChapterContentEditor(),
}))
