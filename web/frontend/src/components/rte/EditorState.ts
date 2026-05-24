import { Editor } from '@tiptap/core'
import { useLayoutEffect, useRef, useState } from 'preact/hooks'

type EditorState = {
  bold: boolean
  italic: boolean
  strikethrough: boolean
  color: string | null
  header: number | null
  font: string | null
  fontSize: string | null
  textAlign: 'left' | 'right' | 'center' | 'justify' | null
}

// eslint-disable-next-line no-redeclare
namespace EditorState {
  export const DEFAULT: EditorState = {
    bold: false,
    italic: false,
    strikethrough: false,
    color: null,
    header: null,
    font: null,
    fontSize: null,
    textAlign: 'left',
  }

  export function get(editor: Editor): EditorState {
    const textStyle = editor.getAttributes('textStyle')

    let textAlign: EditorState['textAlign'] = 'left'

    if (editor.isActive({ textAlign: 'right' })) {
      textAlign = 'right'
    } else if (editor.isActive({ textAlign: 'center' })) {
      textAlign = 'center'
    } else if (editor.isActive({ textAlign: 'justify' })) {
      textAlign = 'justify'
    }

    return {
      bold: editor.isActive('bold'),
      italic: editor.isActive('italic'),
      strikethrough: editor.isActive('strikethrough'),
      color: textStyle.color || null,
      header: editor.isActive('heading') ? editor.getAttributes('heading').level : null,
      font: textStyle.fontFamily || null,
      fontSize: textStyle.fontSize || null,
      textAlign,
    }
  }

  export function useEditorState(editor: Editor): EditorState {
    const stateRef = useRef<EditorState | null>(null)
    const [, fr] = useState(0)
    if (stateRef.current === null) stateRef.current = get(editor)

    useLayoutEffect(() => {
      if (editor.isDestroyed) return

      const onTransaction = () => {
        stateRef.current = get(editor)
        fr((x) => x + 1)
      }

      const cleanup = () => {
        editor.off('transaction', onTransaction)
        editor.off('destroy', onDestroyed)
      }

      const onDestroyed = () => {
        cleanup()
      }

      editor.on('transaction', onTransaction)
      editor.on('destroy', onDestroyed)
    }, [editor])

    return stateRef.current
  }
}

export { EditorState }
