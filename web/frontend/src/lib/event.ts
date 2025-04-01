export type Callback<T> = (event: T) => void

export interface InternalEvent<T> {
  subscribe(callback: Callback<T>): void
  unsubscribe(callback: Callback<T>): void
  fire(event: T): void
}

class EventSourceImpl<T> implements InternalEvent<T> {
  fire(event: T): void {
    for (const callback of this._callbacks) {
      callback(event)
    }
  }
  private readonly _callbacks: Callback<T>[] = []
  subscribe(callback: Callback<T>): void {
    if (this._callbacks.indexOf(callback) === -1) {
      this._callbacks.push(callback)
    }
  }
  unsubscribe(callback: Callback<T>): void {
    const idx = this._callbacks.indexOf(callback)
    if (idx !== -1) {
      this._callbacks.splice(idx, 1)
    }
  }
}

export function createEvent<T>(): InternalEvent<T> {
  return new EventSourceImpl<T>()
}
