import { Widget } from './core'

export const getBaseWidgets: () => Widget[] = () => [
  {
    name: window._('editor.h2'),
    icon: <i class="fa-solid fa-heading" />,
    apply: (editor) => editor.chain().focus().toggleHeading({ level: 2 }).run(),
  },
  {
    name: window._('editor.h3'),
    icon: <i class="fa-solid fa-heading" />,
    apply: (editor) => editor.chain().focus().toggleHeading({ level: 3 }).run(),
  },
  {
    name: window._('editor.ul'),
    description: 'Create a bullet list',
    apply: (editor) => editor.chain().focus().toggleBulletList().run(),
  },
  {
    name: window._('editor.ol'),
    description: 'Create an ordered list',
    apply: (editor) => editor.chain().focus().toggleOrderedList().run(),
  },
]
