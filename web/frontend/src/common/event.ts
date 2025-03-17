type Callback<T> = (event: T) => void

export interface EventLike<TEvent> {
  subscribe(callback: Callback<TEvent>): void
  unsubscribe(callback: Callback<TEvent>): void
}

export class OLEvent<TEvent> implements EventLike<TEvent> {
  private readonly _callbacks: Callback<TEvent>[] = []

  unsubscribe(callback: Callback<TEvent>): void {
    const idx = this._callbacks.indexOf(callback)
    if (idx !== -1) {
      this._callbacks.splice(idx, 1)
    }
  }

  subscribe(callback: Callback<TEvent>): void {
    if (this._callbacks.indexOf(callback) === -1) {
      this._callbacks.push(callback)
    }
  }

  publish(event: TEvent) {
    for (const callback of this._callbacks) {
      callback(event)
    }
  }
}
