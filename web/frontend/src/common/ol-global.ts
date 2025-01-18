declare global {
  // eslint-disable-next-line @typescript-eslint/no-empty-object-type
  interface OLGlobal {}

  interface Window {
    OL: OLGlobal
    __server__?: Record<string, unknown>
  }
}

// eslint-disable-next-line no-undef
window.OL = {} as OLGlobal
