export function isAttrTrue(v: string | null) {
  if (v === '' || v === 'true') {
    return true
  }
  return false
}

export function getNextElement(
  element: Element,
  condition: (element: Element) => boolean,
): Element | null {
  let next = element.nextSibling

  while (next && (!(next instanceof Element) || !condition(next))) {
    next = next.nextSibling
  }

  return next
}

export function getPrevElement(
  element: Element,
  condition: (element: Element) => boolean,
): Element | null {
  let next = element.previousSibling

  while (next && (!(next instanceof Element) || !condition(next))) {
    next = next.previousSibling
  }

  return next
}
