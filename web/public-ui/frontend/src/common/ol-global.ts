declare global {
  interface OLGlobal {}

  interface Window {
    OL: OLGlobal;
  }
}

window.OL = {} as OLGlobal