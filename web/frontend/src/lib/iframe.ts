type VirtualElement = {
  getBoundingClientRect: () => DOMRect
}

/**
 * Translates a virtual element inside an iframe into
 * a virtual element relative to the main window.
 */
export function wrapVirtualElement(
  iframe: HTMLIFrameElement,
  innerVirtual: VirtualElement,
): VirtualElement {
  return {
    getBoundingClientRect(): DOMRect {
      const iframeRect = iframe.getBoundingClientRect()
      const innerRect = innerVirtual.getBoundingClientRect()

      return new DOMRect(
        iframeRect.left + innerRect.left,
        iframeRect.top + innerRect.top,
        innerRect.width,
        innerRect.height,
      )
    },
  }
}
