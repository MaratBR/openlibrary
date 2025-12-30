import { BubbleMenu } from '@tiptap/react/menus'

import './EditorBubbleMenu.scss'
import { ChapterContentEditor, useEditorToolbarState } from './editor'
import EditorToggleButton from './EditorToggleButton'
import TextFeatureSelector from './TextFeatureSelector'
import { Editor } from '@tiptap/core'
import { Node as PMNode } from 'prosemirror-model'
import { VirtualElement } from '@floating-ui/react'

export default function EditorBubbleMenu({
  editor,
  appendTo,
}: {
  editor: ChapterContentEditor
  appendTo: HTMLElement
}) {
  const { bold, italic, strikethrough } = useEditorToolbarState(editor)

  return (
    <BubbleMenu
      class="be-bubble-menu"
      // getReferencedVirtualElement={() => {
      //   const textElement = getSelectedTextElement(editor)
      //   return textElement
      // }}
      options={{
        placement: 'top-start',
      }}
      editor={editor}
      appendTo={appendTo}
    >
      <div class="be-toggle-group">
        <EditorToggleButton active={bold} onClick={() => editor.chain().focus().toggleBold().run()}>
          <i class="fa-solid fa-bold" />
        </EditorToggleButton>
        <EditorToggleButton
          active={italic}
          onClick={() => editor.chain().focus().toggleItalic().run()}
        >
          <i class="fa-solid fa-italic" />
        </EditorToggleButton>
        <EditorToggleButton
          active={strikethrough}
          onClick={() => editor.chain().focus().toggleStrike().run()}
        >
          <i class="fa-solid fa-strikethrough" />
        </EditorToggleButton>
      </div>

      <div class="be-bubble-menu__delimiter" />

      <TextFeatureSelector editor={editor} />
    </BubbleMenu>
  )
}

function isValidTextNode(node: PMNode): boolean {
  return node.isTextblock && !node.isAtom
}

export function getSelectedTextElement(editor: Editor): VirtualElement | null {
  const { state, view } = editor
  const { selection, doc } = state
  const { from, to, empty } = selection

  let result: VirtualElement | null = null

  doc.nodesBetween(from, to, (node, pos) => {
    if (result) return false
    if (!isValidTextNode(node)) return

    // cursor only
    if (empty && node.content.size === 0) return

    const dom = view.nodeDOM(pos) as HTMLElement | null
    if (!dom) return

    result = {
      contextElement: dom,
      getBoundingClientRect: () => dom.getBoundingClientRect(),
      getClientRects: () => dom.getClientRects(),
    }

    return false
  })

  return result
}
