import { z } from 'zod'

const iframeChildMessage = z.object({
  iframeId: z.string(),
  width: z.number(),
  height: z.number(),
  type: z.literal('iframe-child-message'),
})

export type ChildMessage = z.infer<typeof iframeChildMessage>

export function parseChildMessage(data: unknown): ChildMessage | null {
  const result = iframeChildMessage.safeParse(data)
  return result.success ? result.data : null
}

export function getIframeId(): string | null {
  if (window.frameElement) {
    return window.frameElement.id
  }
  return null
}
