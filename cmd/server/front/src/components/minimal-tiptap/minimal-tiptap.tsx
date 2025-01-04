import * as React from 'react'
import './styles/index.css'
import './minimal-tiptap.css'

import type { Content, Editor } from '@tiptap/react'
import type { UseMinimalTiptapEditorProps } from './hooks/use-minimal-tiptap'
import { EditorContent } from '@tiptap/react'
import { cn } from '@/lib/utils'
import { SectionOne } from './components/section/one'
import { SectionTwo } from './components/section/two'
import { SectionThree } from './components/section/three'
import { SectionFour } from './components/section/four'
import { SectionFive } from './components/section/five'
import { LinkBubbleMenu } from './components/bubble-menu/link-bubble-menu'
import { useMinimalTiptapEditor } from './hooks/use-minimal-tiptap'

export interface MinimalTiptapProps extends Omit<UseMinimalTiptapEditorProps, 'onUpdate'> {
  value?: Content
  onChange?: (value: Content) => void
  className?: string
  editorContentClassName?: string
}

const Toolbar = ({ editor }: { editor: Editor }) => (
  <div className="shrink-0 overflow-x-auto border-b border-border p-2">
    <div className="editor-toolbar">
      <div className="editor-toolbar__section">
        <SectionOne editor={editor} activeLevels={[1, 2, 3, 4, 5, 6]} />
      </div>

      <div className="editor-toolbar__section">
        <SectionTwo editor={editor} />
      </div>

      <div className="editor-toolbar__section">
        <SectionThree editor={editor} />
      </div>

      <div className="editor-toolbar__section">
        <SectionFour
          editor={editor}
          activeActions={['orderedList', 'bulletList']}
          mainActionCount={0}
        />
      </div>

      <div className="editor-toolbar__section">
        <SectionFive
          editor={editor}
          activeActions={['codeBlock', 'blockquote', 'horizontalRule']}
          mainActionCount={0}
        />
      </div>
    </div>
  </div>
)

export function useMinimalTiptapEditorComponent({
  value,
  onChange,
  className,
  editorContentClassName,
  ...props
}: MinimalTiptapProps) {
  const editor = useMinimalTiptapEditor({
    value,
    onUpdate: onChange,
    ...props,
  })

  return {
    editor,
    editorElement:
      editor === null ? null : (
        <div
          className={cn(
            'flex h-auto min-h-72 w-full flex-col rounded-md border border-input shadow-sm focus-within:border-primary focus-within:outline focus-within:outline-primary focus-within:outline-1',
            className,
          )}
        >
          <Toolbar editor={editor} />
          <EditorContent
            editor={editor}
            className={cn('minimal-tiptap-editor', editorContentClassName)}
          />
          <LinkBubbleMenu editor={editor} />
        </div>
      ),
  }
}
