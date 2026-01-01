export class EditorElements {
  public readonly iframe: HTMLIFrameElement
  public readonly content: HTMLElement
  public readonly contentWrapper: HTMLElement
  public readonly contentWrapperHeader: HTMLElement

  constructor(iframe: HTMLIFrameElement) {
    this.iframe = iframe
    if (!iframe.contentDocument)
      throw new Error('iframe failed to load, contentDocument is not available yet')

    const contentElement = iframe.contentDocument.getElementById('ChapterContent')
    const contentWrapper = iframe.contentDocument.getElementById('BlockEditorWrap')
    const contentWrapperHeader = iframe.contentDocument.getElementById('BlockEditorWrapHeader')

    if (!contentElement) throw new Error('cannot find element #ChapterContent')
    if (!contentWrapper) throw new Error('cannot find element #BlockEditorWrap')
    if (!contentWrapperHeader) throw new Error('cannot find element #BlockEditorWrapHeader')

    this.content = contentElement
    this.contentWrapper = contentWrapper
    this.contentWrapperHeader = contentWrapperHeader
  }
}
