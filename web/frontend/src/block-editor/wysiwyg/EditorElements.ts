export class EditorElements {
  public readonly iframe: HTMLIFrameElement
  public readonly content: HTMLElement
  public readonly contentWrapper: HTMLElement

  constructor(iframe: HTMLIFrameElement) {
    this.iframe = iframe
    if (!iframe.contentDocument)
      throw new Error('iframe failed to load, contentDocument is not available yet')

    const editorWrapElement = iframe.contentDocument.getElementById('BlockEditorWrap')
    const contentElement = iframe.contentDocument.getElementById('ChapterContent')

    if (!editorWrapElement) throw new Error('cannot find element #BlockEditorWrap')
    if (!contentElement) throw new Error('cannot find element #ChapterContent')

    this.content = contentElement
    this.contentWrapper = editorWrapElement
  }
}
