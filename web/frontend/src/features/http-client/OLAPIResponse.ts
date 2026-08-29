import { KyResponse } from 'ky'
import { Schema } from 'effect'
import type { jsonErrorResponse, Notification } from '@/backend-types'

export type OLNotification = Notification

const NO_BODY_SCHEMA = Schema.Literal('ok')

export class OLAPIResponse<T> {
  private readonly response: Response
  private _notifications?: OLNotification[] = undefined
  private _data?: T
  private _error?: jsonErrorResponse
  private readonly _schema: Schema.Codec<T, unknown>

  public static async create<T>(
    response: Response,
    schema: Schema.Codec<T, unknown> = Schema.Unknown as Schema.Codec<T, unknown>,
  ): Promise<OLAPIResponse<T>> {
    const resp = new OLAPIResponse<T>(response, schema)
    await resp._loadData()
    return resp
  }

  public static async createNoBody(
    response: Response,
  ): Promise<OLAPIResponse<Schema.Schema.Type<typeof NO_BODY_SCHEMA>>> {
    return OLAPIResponse.create(response, NO_BODY_SCHEMA)
  }

  private constructor(response: Response, schema: Schema.Codec<T, unknown>) {
    this.response = response
    this._schema = schema
  }

  get status() {
    return this.response.status
  }

  get ok() {
    return this.response.ok
  }

  get error() {
    return this._error
  }

  get notifications() {
    if (this._notifications === undefined) {
      this._notifications = OLAPIResponse.parseNotifications(this.response)
    }

    return this._notifications
  }

  get data() {
    if (this._data === undefined) throw new Error('internal _data property not initialized')

    return this._data
  }

  throwIfError() {
    if (this._error) {
      throw new OLAPIError(this)
    }
  }

  private async _loadData() {
    if (this._data !== undefined) return

    const json = await this.response.json()

    if (this.response.status >= 400 && this.response.status <= 599) {
      this._error = json as jsonErrorResponse
    }

    this._data = Schema.decodeUnknownSync(this._schema)(json)
  }

  private static parseNotifications(response: KyResponse): OLNotification[] {
    const flashes = response.headers.get('x-flash')
    if (!flashes) return []

    try {
      const json = JSON.parse(flashes)
      return Array.isArray(json) ? (json as OLNotification[]) : []
    } catch (e: unknown) {
      console.warn('failed to parse x-flash header contents', e)
      return []
    }
  }
}

export class OLAPIError extends Error {
  private readonly _response: OLAPIResponse<unknown>

  get error() {
    return this._response.error
  }

  constructor(response: OLAPIResponse<unknown>) {
    const responseErrorMessage = response.error ? response.error.message : 'no error'
    super(`OLAPI error: ${responseErrorMessage}`)
    this._response = response
  }
}
