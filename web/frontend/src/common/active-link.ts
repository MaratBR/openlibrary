function init() {
  const links = document.querySelectorAll('a')

  const activeLinks = []

  for (const link of links) {
    if (isActive(link)) {
      activeLinks.push(link)
    }
  }

  for (const activeLink of activeLinks) {
    activeLink.dataset.active = 'true'
  }
}

function isActive(link: HTMLAnchorElement): boolean {
  try {
    const href = new URL(link.href)
    const currentUrl = new URL(window.location.href)
    const currentPath = currentUrl.pathname

    return href.pathname === currentPath && href.host === currentUrl.host
  } catch {
    return false
  }
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
