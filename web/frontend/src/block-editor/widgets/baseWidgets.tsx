import { Widget } from './core'

export const getBaseWidgets: () => Widget[] = () => [
  {
    name: window._('editor.p'),
    icon: <i className="fa-solid fa-paragraph" />,
    apply: (editor) => editor.chain().focus().setParagraph().run(),
  },
  ...([1, 2, 3, 4, 5, 6] as const).map((level) => ({
    name: window._(`editor.h${level}`),
    icon: <i className="fa-solid fa-heading" />,
    apply: (editor: Parameters<Widget['apply']>[0]) =>
      editor.chain().focus().toggleHeading({ level }).run(),
  })),
  {
    name: window._('editor.bold'),
    icon: <i className="fa-solid fa-bold" />,
    apply: (editor) => editor.chain().focus().toggleBold().run(),
  },
  {
    name: window._('editor.italic'),
    icon: <i className="fa-solid fa-italic" />,
    apply: (editor) => editor.chain().focus().toggleItalic().run(),
  },
  {
    name: window._('editor.strike'),
    icon: <i className="fa-solid fa-strikethrough" />,
    apply: (editor) => editor.chain().focus().toggleStrike().run(),
  },
  {
    name: window._('editor.underline'),
    icon: <i className="fa-solid fa-underline" />,
    apply: (editor) => editor.chain().focus().toggleUnderline().run(),
  },
  {
    name: window._('editor.code'),
    icon: <i className="fa-solid fa-code" />,
    apply: (editor) => editor.chain().focus().toggleCode().run(),
  },
  {
    name: window._('editor.codeBlock'),
    icon: <i className="fa-solid fa-file-code" />,
    apply: (editor) => editor.chain().focus().toggleCodeBlock().run(),
  },
  {
    name: window._('editor.blockquote'),
    icon: <i className="fa-solid fa-quote-left" />,
    apply: (editor) => editor.chain().focus().toggleBlockquote().run(),
  },
  {
    name: window._('editor.ul'),
    icon: <i className="fa-solid fa-list-ul" />,
    apply: (editor) => editor.chain().focus().toggleBulletList().run(),
  },
  {
    name: window._('editor.ol'),
    icon: <i className="fa-solid fa-list-ol" />,
    apply: (editor) => editor.chain().focus().toggleOrderedList().run(),
  },
  {
    name: window._('editor.horizontalRule'),
    icon: <i className="fa-solid fa-minus" />,
    apply: (editor) => editor.chain().focus().setHorizontalRule().run(),
  },
]
