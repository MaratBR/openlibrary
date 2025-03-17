import { KyResponse } from 'ky'
import { z, ZodSchema } from 'zod'

const notificationSchema = z.object({
  type: z.enum(['info', 'error']),
  text: z.string(),
})

const anyArray = z.array(z.any())

export type OLNotification = z.infer<typeof notificationSchema>

export class OLAPIResponse<T> {
  private readonly response: Response
  private _notifications?: OLNotification[] = undefined
  private _data?: T
  private readonly _schema: ZodSchema<T>

  public static async create<T>(
    response: Response,
    schema: ZodSchema<T>,
  ): Promise<OLAPIResponse<T>> {
    const resp = new OLAPIResponse<T>(response, schema)
    await resp._loadData()
    return resp
  }

  private constructor(response: Response, schema: ZodSchema<T>) {
    this.response = response
    this._schema = schema
  }

  get status() {
    return this.response.status
  }

  get ok() {
    return this.response.ok
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

  private async _loadData() {
    if (this._data !== undefined) return

    this._data = await this._schema.parseAsync(await this.response.json())
  }

  private static parseNotifications(response: KyResponse): OLNotification[] {
    const flashes = response.headers.get('x-flash')
    if (!flashes) return []

    try {
      const json = JSON.parse(flashes)
      const arr = anyArray.parse(json)
      const notifications: OLNotification[] = []

      for (let i = 0; i < arr.length; i++) {
        const el = arr[i]
        const result = notificationSchema.safeParse(el)
        if (result.success) {
          notifications.push(result.data)
        } else {
          console.warn(
            `failed to parse value as server notification at position ${i}`,
            result.error,
          )
        }
      }

      return notifications
    } catch (e: unknown) {
      console.warn('failed to parse x-flash header contents', e)
      return []
    }
  }
}
