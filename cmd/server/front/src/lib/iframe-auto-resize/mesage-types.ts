import { z } from 'zod'

export const iframeChildMessage = z.object({
  iframeId: z.string(),
  width: z.number(),
  height: z.number(),
  type: z.literal('iframeC'),
})

export const iframeParentMessage = z.object({
  type: z.literal('iframePWndSize'),
  width: z.number(),
  height: z.number(),
})

export type IframeParentMessage = z.infer<typeof iframeParentMessage>
export type IframeChildMessage = z.infer<typeof iframeChildMessage>

export function getIframeId(): string | null {
  if (window.frameElement) {
    return window.frameElement.id
  }
  return null
}
