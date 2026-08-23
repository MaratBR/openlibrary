import { httpClient } from '@/http-client'

export type Font = {
  name: string
  source: string
  license: string
  includeCss: string
  externalLink: string
}

export namespace FontsApi {
  export function getGoogleFonts(): Promise<Font[]> {
    return httpClient.get('/_api/fonts/google').then((r) => r.json())
  }
}
