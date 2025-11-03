import { Derived, Subject } from './rx'
import { debounce } from './util/fn'

function findAllNavigationContents(parent: HTMLElement): {
  element: HTMLTemplateElement
  name: string
}[] {
  const contentTemplates: {
    element: HTMLTemplateElement
    name: string
  }[] = []
  const templates = parent.querySelectorAll('template[data-nav-content]')

  for (const template of templates) {
    if (template instanceof HTMLTemplateElement) {
      contentTemplates.push({
        element: template,
        name: template.dataset.navContent || '',
      })
    } else {
      console.warn('Element with data-nav-content is not a template', template)
    }
  }

  return contentTemplates
}

function findAllNavigationTriggers(parent: HTMLElement): {
  element: HTMLElement
  name: string
}[] {
  const elements: {
    element: HTMLElement
    name: string
  }[] = []
  const triggers = parent.querySelectorAll('[data-nav-trigger]')

  for (const trigger of triggers) {
    if (trigger instanceof HTMLElement) {
      elements.push({
        element: trigger,
        name: trigger.dataset.navTrigger || '',
      })
    } else {
      console.warn('Element with data-nav-trigger is not an HTMLElement', trigger)
    }
  }

  return elements
}

export function setupNavigation(parent: HTMLElement) {
  const portal = parent.querySelector('[data-nav-portal]')
  if (!portal) {
    console.warn('No navigation portal found')
    return
  }
  if (!(portal instanceof HTMLElement)) {
    console.warn('Navigation portal is not an HTMLElement')
    return
  }

  const active = new Subject('')
  const show = new Subject(false)

  const contentTemplates = findAllNavigationContents(parent)
  const findAllTriggers = findAllNavigationTriggers(parent)

  const setShowDebounced = debounce(show.set.bind(show), 350)

  for (const triggerInfo of findAllTriggers) {
    const onMouseEnter = debounce(() => {
      if (active.get() !== triggerInfo.name || show.get() === false) {
        active.set(triggerInfo.name)
        setShowDebounced.cancel()
        show.set(true)
      }
    }, 400)

    triggerInfo.element.addEventListener('mouseenter', onMouseEnter)

    triggerInfo.element.addEventListener('mouseleave', () => {
      onMouseEnter.cancel()
      if (active.get() === triggerInfo.name) {
        setShowDebounced(false)
      }
    })
  }

  portal.addEventListener('mouseenter', () => {
    setShowDebounced.cancel()
  })

  portal.addEventListener('mouseleave', () => {
    setShowDebounced(false)
  })

  const portalContentReady = new Subject(false)

  new Derived([portalContentReady, show], (portalReady, show) => portalReady && show).subscribe(
    (visible) => {
      portal.setAttribute('data-open', visible ? 'true' : 'false')
    },
  )

  active.subscribe((name) => {
    portal.innerHTML = ''
    const contentInfo = contentTemplates.find((content) => content.name === name)
    if (!contentInfo) {
      return
    }
    const trigger = findAllTriggers.find((t) => t.name === name)
    const content = contentInfo.element.content.cloneNode(true)
    portal.appendChild(content)

    if (trigger) {
      const { left, bottom, top, right } = trigger.element.getBoundingClientRect()

      if (!portalContentReady.get()) {
        portal.classList.add('no-transition')
      }

      window.requestAnimationFrame(() => {
        portal.style.setProperty('--nav-trigger-left', `${left}px`)
        portal.style.setProperty('--nav-trigger-bottom', `${bottom}px`)
        portal.style.setProperty('--nav-trigger-top', `${top}px`)
        portal.style.setProperty('--nav-trigger-right', `${right}px`)

        if (!portalContentReady.get()) {
          window.requestAnimationFrame(() => {
            portal.classList.remove('no-transition')
            portalContentReady.set(true)
          })
        }
      })
    }
  })
}

document.addEventListener('DOMContentLoaded', () => {
  setupNavigation(document.body)
})
