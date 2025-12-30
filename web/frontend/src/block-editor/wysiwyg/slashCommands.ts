import { SlashCommandItem } from './Suggestions'

export const slashCommands: () => SlashCommandItem[] = () => [
  {
    name: window._('editor.h2'),
    command: (editor) => editor.chain().focus().toggleHeading({ level: 2 }).run(),
  },
  {
    name: window._('editor.h3'),
    command: (editor) => editor.chain().focus().toggleHeading({ level: 3 }).run(),
  },
  {
    name: window._('editor.ul'),
    description: 'Create a bullet list',
    command: (editor) => editor.chain().focus().toggleBulletList().run(),
  },
  {
    name: window._('editor.ol'),
    description: 'Create an ordered list',
    command: (editor) => editor.chain().focus().toggleOrderedList().run(),
  },
]
