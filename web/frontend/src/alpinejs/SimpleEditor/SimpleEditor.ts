import Alpine from 'alpinejs'

import type { SimpleEditor } from './SimpleEditorImpl'

let id = 0

const loadEditor = () => import('./SimpleEditorImpl')

Alpine.data('SimpleEditor', () => ({
  _id: `${id++}`,

  init() {
    loadEditor().then((module) => {
      this.$store[`SimpleEditor:${this._id}`] = new module.SimpleEditor(this.$el)
    })
  },

  destroy() {
    const editor = this.$store[`SimpleEditor:${this._id}`] as SimpleEditor | undefined
    if (editor) {
      editor.destroy()
    }
  },
}))
