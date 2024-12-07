import { ZodSchema } from 'zod'
import { debounce } from '../utils'
import {
  getIframeId,
  IframeChildMessage,
  iframeChildMessage,
  IframeParentMessage,
  iframeParentMessage,
} from './mesage-types'

const PARAM_HEIGHT = 'iframe-height'
const PARAM_WIDTH = 'iframe-width'

function setIframeWidth(v: number) {
  window.document.body.style.setProperty('--iframe-outer-width', `${v}px`)
}
function setIframeHeight(v: number) {
  window.document.body.style.setProperty('--iframe-outer-height', `${v}px`)
}

export function initIframeAgent(): () => void {
  if (window.top !== window.self) {
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

  const unsubscribe = subscribeToIframeMessage(iframeParentMessage, (message) => {
    if (message.type === 'iframePWndSize') {
      window.requestAnimationFrame(() => {
        setIframeWidth(message.width)
        setIframeHeight(message.height)
      })
    }
  })

  const pushUpdate = () => {
    window.requestAnimationFrame(() => {
      publishChildMessage(
        {
          iframeId,
          width: window.document.body.scrollWidth,
          height: window.document.body.scrollHeight,
          type: 'iframeC',
        },
        window.parent,
      )
    })
  }

  const resizeObserver = new ResizeObserver(pushUpdate)

  resizeObserver.observe(window.document.body)

  pushUpdate()

  const urlParams = new URLSearchParams(window.location.search)

  if (urlParams.has(PARAM_HEIGHT)) {
    const height = +urlParams.get(PARAM_HEIGHT)!
    if (!Number.isNaN(height)) {
      setIframeHeight(height)
    }
  }

  if (urlParams.has(PARAM_WIDTH)) {
    const width = +urlParams.get(PARAM_WIDTH)!
    if (!Number.isNaN(width)) {
      setIframeWidth(width)
    }
  }

  return () => {
    resizeObserver.disconnect()
    unsubscribe()
  }
}

function initParent(): () => void {
  const unsubscribe = subscribeToIframeMessage(iframeChildMessage, (event: IframeChildMessage) => {
    if (event.type === 'iframeC') {
      const iframe = document.getElementById(event.iframeId)

      if (iframe instanceof HTMLIFrameElement) {
        iframe.style.height = `${Math.ceil(event.height)}px`
      }
    }
  })

  const onResize = debounce(() => {
    window.requestAnimationFrame(() => {
      const height = window.innerHeight
      const width = window.innerWidth
      const message: IframeParentMessage = {
        type: 'iframePWndSize',
        width,
        height,
      }
      for (let i = 0; i < window.frames.length; i++) {
        const frame = window.frames[i]
        frame.postMessage(message, '*')
      }
    })
  }, 100)

  window.addEventListener('resize', onResize)

  return () => {
    window.removeEventListener('resize', onResize)
    unsubscribe()
  }
}

function subscribeToIframeMessage<T>(schema: ZodSchema<T>, callback: (event: T) => void) {
  const onMessage = (event: MessageEvent<unknown>) => {
    const childMessageResult = schema.safeParse(event.data)
    if (childMessageResult.success) {
      callback(childMessageResult.data)
    }
  }

  window.addEventListener('message', onMessage)

  return () => {
    window.removeEventListener('message', onMessage)
  }
}

function publishChildMessage(message: IframeChildMessage, wnd: Window) {
  wnd.postMessage(message, '*')
}

export function getIframeUrl(
  urlString: string,
  includeWidth: boolean,
  includeHeight: boolean,
): string {
  const url = new URL(window.location.origin + urlString)

  if (includeHeight) {
    url.searchParams.set(PARAM_HEIGHT, String(window.innerHeight))
  }

  if (includeWidth) {
    url.searchParams.set(PARAM_WIDTH, String(window.innerWidth))
  }

  return url.toString()
}
