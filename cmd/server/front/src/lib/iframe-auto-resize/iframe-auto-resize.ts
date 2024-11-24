import { ChildMessage, getIframeId, parseChildMessage } from './mesage-types'

export function isInsideIframe(): boolean {
  return window.top !== window.self
}

export function initIframeAgent(): () => void {
  if (isInsideIframe()) {
    return initChild()
  } else {
    return initParent()
  }
}

function initChild(): () => void {
  const iframeId = getIframeId()
  if (!iframeId) {
    return () => {}
  }

  if (!document.body.classList.contains('__iframe-root'))
    document.body.classList.add('__iframe-root')

  const pushUpdate = () => {
    window.requestAnimationFrame(() => {
      publishChildMessage(
        {
          iframeId,
          width: window.document.body.scrollWidth,
          height: window.document.body.scrollHeight,
          type: 'iframe-child-message',
        },
        window.parent,
      )
    })
  }

  const resizeObserver = new ResizeObserver(pushUpdate)

  resizeObserver.observe(window.document.body)

  pushUpdate()

  return () => {
    resizeObserver.disconnect()
  }
}

function initParent(): () => void {
  const unsubscribe = subscribeToChildMessage((event: ChildMessage) => {
    const iframe = document.getElementById(event.iframeId)

    if (iframe instanceof HTMLIFrameElement) {
      iframe.style.height = `${Math.ceil(event.height) + 1}px`
    }
  })

  return () => {
    unsubscribe()
  }
}

function subscribeToChildMessage(callback: (event: ChildMessage) => void) {
  const onMessage = (event: MessageEvent<unknown>) => {
    const childMessage = parseChildMessage(event.data)
    if (childMessage) {
      callback(childMessage)
    }
  }

  window.addEventListener('message', onMessage)

  return () => {
    window.removeEventListener('message', onMessage)
  }
}

function publishChildMessage(message: ChildMessage, wnd: Window) {
  wnd.postMessage(message, '*')
}
