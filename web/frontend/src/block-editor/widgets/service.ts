import { getBaseWidgets } from './baseWidgets'
import { Widget } from './core'

export class WidgetsService {
  private static _INSTANCE?: WidgetsService

  public static instance() {
    if (!this._INSTANCE) this._INSTANCE = new WidgetsService()
    return this._INSTANCE
  }

  public async loadWidgets(): Promise<Widget[]> {
    return getBaseWidgets()
  }
}
