import { Editor, Extension } from '@tiptap/core'
import Suggestion, { SuggestionProps, Trigger } from '@tiptap/suggestion'
import { keymap } from '@tiptap/pm/keymap'
import { JSX } from 'preact/jsx-runtime'
// import { slashCommands } from './slashCommands'

export type SlashCommandItem = {
  name: string
  icon?: JSX.Element
  description?: string
  command(editor: Editor): void
}

export interface SlashCommandDisplayAdapter {
  show(props: SuggestionProps<SlashCommandItem, SlashCommandItem>): void
  update(props: SuggestionProps<SlashCommandItem, SlashCommandItem>): void
  hide(): void
  focus(key: 'ArrowDown' | 'ArrowUp'): void
}

export type SlashCommandDisplayFocusKey = 'ArrowUp' | 'ArrowDown' | null

export type SlashCommandOptions = {
  commands: SlashCommandItem[]
  suggestionClass: string
  displayAdapter: SlashCommandDisplayAdapter
}

type SuggestsionsStorage = {
  open: boolean
}

export const SlashCommand = Extension.create<SlashCommandOptions, SuggestsionsStorage>({
  name: 'slash-command',

  addStorage() {
    return {
      open: false,
      activeKey: null,
    }
  },

  addProseMirrorPlugins() {
    return [
      Suggestion<SlashCommandItem, SlashCommandItem>({
        editor: this.editor,
        char: '/',
        startOfLine: true,
        decorationClass: this.options.suggestionClass,
        findSuggestionMatch: findSlashLineMatch,

        items: ({ query }) => {
          return this.options.commands.filter((item) =>
            item.name.toLowerCase().includes(query.toLowerCase()),
          )
        },

        command: ({ editor, range, props }) => {
          // delete the "/query"
          editor.chain().focus().deleteRange(range).run()
          props.command(editor)
        },

        render: () => {
          return {
            onStart: (props) => {
              this.options.displayAdapter.show(props)
              this.storage.open = props.items.length > 0
            },

            onUpdate: (props) => {
              this.options.displayAdapter.update(props)
              this.storage.open = props.items.length > 0
            },

            onExit: () => {
              this.options.displayAdapter.hide()
              this.storage.open = false
            },
          }
        },
      }),
      keymap({
        ArrowDown: () => {
          if (this.storage.open) {
            this.options.displayAdapter.focus('ArrowDown')
            return true
          }

          return false
        },
        ArrowUp: () => {
          if (this.storage.open) {
            this.options.displayAdapter.focus('ArrowUp')
            return true
          }

          return false
        },
      }),
    ]
  },
})

function findSlashLineMatch({ $position }: Trigger) {
  const parent = $position.parent
  const text = parent.textContent

  if (!text.startsWith('/')) {
    return null
  }

  const from = $position.start()
  const to = from + parent.nodeSize - 2 // text only

  return {
    range: { from, to },
    query: text.slice(1), // everything after "/"
    text,
  }
}
