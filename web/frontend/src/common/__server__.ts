type __server__ = {
  ageRatings: string[]
  bookId?: string
  user?: unknown | null
}

declare global {
  // eslint-disable-next-line @typescript-eslint/no-empty-object-type
  interface OLGlobal {}

  interface Window {
    OL: OLGlobal
    __server__: __server__
  }
}

// eslint-disable-next-line no-undef
window.OL = {} as OLGlobal
