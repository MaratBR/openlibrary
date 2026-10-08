import { executeAfterDOMIsReady } from './dom'
import { initflashes } from './flashes'

export function initActiveLinks() {
  // React Router owns hash-routed workspace links. A document-wide pathname
  // match would select every workspace route if its island mounts before us.
  const links = document.querySelectorAll<HTMLAnchorElement>('a.nav-link, a.NavSubmenu-link')
  const currentUrl = new URL(window.location.href)

  for (const link of links) {
    const href = new URL(link.href)
    const active =
      href.origin === currentUrl.origin && href.pathname === currentUrl.pathname && !href.hash
    link.classList.toggle('active', active)
  }
}

function init() {
  initActiveLinks()
  initflashes()
}

export function initAfterDOMReady() {
  executeAfterDOMIsReady(() => requestAnimationFrame(init))
}
