import { initFlashMessages } from './flashes'

function initActiveLinks() {
  const links = document.querySelectorAll('a')

  const activeLinks: {
    link: HTMLAnchorElement
    activeType: 'path' | 'full'
  }[] = []

  for (const link of links) {
    const activeType = isActive(link)
    if (activeType) {
      activeLinks.push({
        link,
        activeType,
      })
    }
  }

  for (const { link, activeType } of activeLinks) {
    link.dataset.linkActive = 'true'
    if (activeType === 'full') {
      link.dataset.linkActiveFull = 'true'
    }
  }
}

function isActive(link: HTMLAnchorElement): 'full' | 'path' | false {
  try {
    const href = new URL(link.href)
    const currentUrl = new URL(window.location.href)
    const currentPath = currentUrl.pathname

    if (href.pathname === currentPath && href.host === currentUrl.host) {
      if (href.search === currentUrl.search) {
        return 'full'
      }
      return 'path'
    }
    return false
  } catch {
    return false
  }
}

function init() {
  initActiveLinks()
  initFlashMessages()
}

export function initAfterDOMReady() {
  if (document.readyState === 'complete') {
    requestAnimationFrame(init)
  } else {
    document.addEventListener('DOMContentLoaded', () => {
      requestAnimationFrame(init)
    })
  }
}
