import React from 'react'
import MinimalTiptapEditor, { MinimalTiptapProps } from './minimal-tiptap'
import { Editor } from '@tiptap/core'

export type BackupEditorProps = Omit<MinimalTiptapProps, 'editorRef'> & {
  cacheKey?: string
  storage?: {
    setItem: (key: string, value: string) => void
    getItem: (key: string) => string | null | undefined
  }
}

export default class BackupEditor extends React.Component<BackupEditorProps> {
  constructor(props: Readonly<BackupEditorProps>) {
    super(props)
    this._editorRef = this._editorRef.bind(this)
  }

  private _editor: Editor | null = null
  private _backupInterval: number | null = null

  get storage(): BackupEditorProps['storage'] & {} {
    return this.props.storage ?? window.localStorage
  }

  componentDidMount(): void {
    const { cacheKey } = this.props
    if (cacheKey) {
      const value = this.storage.getItem(cacheKey)
      if (value) {
        // this._editor?.commands.setContent(value)
      }
    }
  }

  render(): React.ReactNode {
    const { cacheKey: _, ...props } = this.props

    return <MinimalTiptapEditor editorRef={this._editorRef} {...props} />
  }

  private _editorRef(editor: Editor | null) {
    this._editor = editor

    if (this._backupInterval) clearInterval(this._backupInterval)
    this._backupInterval = window.setInterval(() => {
      this._backup()
    }, 5000)
  }

  private _getStorageKey(): string | undefined {
    const { cacheKey } = this.props
    if (!cacheKey) return undefined

    return `editor-backup-${cacheKey}`
  }

  private _backup() {
    if (!this._editor) {
      return
    }

    const storageKey = this._getStorageKey()
    if (!storageKey) {
      return
    }

    const html = this._editor.getHTML()
    this.storage.setItem(storageKey, html)
  }
}
