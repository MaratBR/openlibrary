import * as React from 'react'
import type { Editor } from '@tiptap/core'
import type { Content, UseEditorOptions } from '@tiptap/react'
import { StarterKit } from '@tiptap/starter-kit'
import { useEditor } from '@tiptap/react'
import { Typography } from '@tiptap/extension-typography'
import { Placeholder } from '@tiptap/extension-placeholder'
import { TextStyle } from '@tiptap/extension-text-style'
import {
  Image,
  HorizontalRule,
  CodeBlockLowlight,
  Selection,
  Color,
  UnsetAllMarks,
  ResetMarksOnEnter,
  FileHandler,
} from '../extensions'
import { cn } from '@/lib/utils'
import { blobUrlToBase64, getOutput } from '../utils'
import { useThrottle } from '../hooks/use-throttle'

export type UseMinimalTiptapEditorProps = Omit<UseEditorOptions, 'extensions'> & {
  value?: Content
  output?: 'html' | 'json' | 'text'
  placeholder?: string
  editorClassName?: string
  throttleDelay?: number
  onUpdate?: (content: Content) => void
  onBlur?: (content: Content) => void
  extensions?: TiptapExtensionsOptions
}

export type TiptapExtensionsOptions = {
  disableImage?: boolean
  disableLink?: boolean
  disabledColor?: boolean
  disableHeadings?: boolean
}

const createExtensions = (placeholder: string, options: TiptapExtensionsOptions) => {
  const extensions = [
    StarterKit.configure({
      horizontalRule: false,
      codeBlock: false,
      // paragraph: { HTMLAttributes: { class: 'text-node' } },
      heading: options.disableHeadings ? false : undefined,
      // blockquote: { HTMLAttributes: { class: 'block-node' } },
      // bulletList: { HTMLAttributes: { class: 'list-node' } },
      // orderedList: { HTMLAttributes: { class: 'list-node' } },
      code: { HTMLAttributes: { class: 'inline', spellcheck: 'false' } },
      dropcursor: { width: 2, class: 'ProseMirror-dropcursor border' },
    }),
    TextStyle,
    Selection,
    Typography,
    UnsetAllMarks,
    HorizontalRule,
    ResetMarksOnEnter,
    CodeBlockLowlight,
    Placeholder.configure({ placeholder: () => placeholder }),
  ]

  if (!options.disabledColor) {
    extensions.push(Color)
  }

  if (!options.disableLink) {
    extensions.push(Image)
  }

  if (!options.disableImage) {
    extensions.push(
      Image.configure({
        allowedMimeTypes: ['image/*'],
        maxFileSize: 5 * 1024 * 1024,
        allowBase64: true,
        uploadFn: async (file) => {
          // wait 5s to simulate a slow upload
          await new Promise((resolve) => setTimeout(resolve, 500))
          const url = await blobUrlToBase64(file)
          return url
        },
        customCopyLink(props, options) {
          console.log('customCopyLink', props, options)
        },
        onValidationError(errors) {
          errors.forEach((error) => {
            console.log('Image validation error', error)
          })
        },
        onActionSuccess(props) {
          console.log('Image action success', props)
        },
        onActionError(error, props) {
          console.error('Image action error', error, props)
        },
      }),
      FileHandler.configure({
        allowBase64: true,
        allowedMimeTypes: ['image/*'],
        maxFileSize: 5 * 1024 * 1024,
        onDrop: (editor, files, pos) => {
          files.forEach((file) =>
            editor.commands.insertContentAt(pos, {
              type: 'image',
              attrs: { src: URL.createObjectURL(file) },
            }),
          )
        },
        onPaste: (editor, files) => {
          files.forEach((file) =>
            editor.commands.insertContent({
              type: 'image',
              attrs: { src: URL.createObjectURL(file) },
            }),
          )
        },
        onValidationError: (errors) => {
          errors.forEach((error) => {
            console.log('File validation error', error)
          })
        },
      }),
    )
  }

  return extensions
}

export const useMinimalTiptapEditor = ({
  value,
  output = 'html',
  placeholder = '',
  editorClassName,
  throttleDelay = 0,
  onUpdate,
  onBlur,
  extensions,
  ...props
}: UseMinimalTiptapEditorProps) => {
  const throttledSetValue = useThrottle((value: Content) => onUpdate?.(value), throttleDelay)

  const handleUpdate = React.useCallback(
    (editor: Editor) => throttledSetValue(getOutput(editor, output)),
    [output, throttledSetValue],
  )

  const handleCreate = React.useCallback(
    (editor: Editor) => {
      if (value && editor.isEmpty) {
        editor.commands.setContent(value)
      }
    },
    [value],
  )

  const handleBlur = React.useCallback(
    (editor: Editor) => onBlur?.(getOutput(editor, output)),
    [output, onBlur],
  )

  const editor = useEditor({
    extensions: createExtensions(placeholder, extensions ?? {}),
    editorProps: {
      attributes: {
        autocomplete: 'off',
        autocorrect: 'off',
        autocapitalize: 'off',
        class: cn('focus:outline-none', editorClassName),
      },
    },
    onUpdate: ({ editor }) => handleUpdate(editor),
    onCreate: ({ editor }) => handleCreate(editor),
    onBlur: ({ editor }) => handleBlur(editor),
    ...props,
  })

  return editor
}

export default useMinimalTiptapEditor
