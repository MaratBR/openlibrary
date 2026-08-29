import { httpClient } from './client'
import { OLAPIResponse, OLNotification } from './OLAPIResponse'

export { httpClient, OLAPIResponse }
export type { OLNotification }

declare global {
  interface Window {
    OLAPIResponse: typeof OLAPIResponse
  }
}

window.OLAPIResponse = OLAPIResponse
