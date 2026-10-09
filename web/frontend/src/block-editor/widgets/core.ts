import { JSX } from 'react'
import { ChapterContentEditor } from '../wysiwyg/editor'

export interface Widget {
  name: string
  availableInToolbar?: boolean
  description?: string
  icon?: JSX.Element
  apply(editor: ChapterContentEditor): void
}
