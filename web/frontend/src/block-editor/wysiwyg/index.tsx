import { useEffect } from 'react'
import { useAtomValue, useSetAtom } from 'jotai'
import { Effect } from 'effect'
import './style.scss'
import { ChapterContentEditor, ChapterContentEditorOptions } from './editor'
import { EditorContent } from '@tiptap/react'
import { EditorBubbleMenu } from './EditorBubbleMenu'
import EditorFloatingMenu from './EditorFloatingMenu'
import { mountWysiwygEditorAtom, wysiwygEditorAtom } from './state'

export function WYSIWYGEditor({ editorOptions }: { editorOptions: ChapterContentEditorOptions }) {
  const editor = useAtomValue(wysiwygEditorAtom)
  const mountEditor = useSetAtom(mountWysiwygEditorAtom)

  useEffect(() => {
    const editor = new ChapterContentEditor(editorOptions)
    return Effect.runSync(mountEditor(editor))
  }, [editorOptions, mountEditor])

  if (!editor) return null

  return (
    <>
      <EditorContent editor={editor} />
      <EditorBubbleMenu editor={editor} appendTo={editorOptions.contentWrapperElement} />
      <EditorFloatingMenu editor={editor} />
    </>
  )
}
