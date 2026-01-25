import { useLayoutEffect } from 'preact/hooks'
import './style.scss'
import { ChapterContentEditor, ChapterContentEditorOptions } from './editor'
import { EditorContent } from '@tiptap/react'
import EditorBubbleMenu from './EditorBubbleMenu'
import EditorFloatingMenu from './EditorFloatingMenu'
import { useWYSIWYG } from './state'

export function WYSIWYGEditor({ editorOptions }: { editorOptions: ChapterContentEditorOptions }) {
  const editor = useWYSIWYG((s) => s.editor)

  useLayoutEffect(() => {
    const editor = new ChapterContentEditor(editorOptions)
    useWYSIWYG.getState().init(editor)

    return () => {
      editor.destroy()
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  if (!editor) return null

  return (
    <>
      <EditorContent editor={editor} />
      <EditorBubbleMenu editor={editor} appendTo={editorOptions.contentWrapperElement} />
      <EditorFloatingMenu editor={editor} />
    </>
  )
}
