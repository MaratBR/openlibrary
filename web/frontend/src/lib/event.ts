export type Callback<T> = (event: T) => void

export interface InternalEvent<T> {
  subscribe(callback: Callback<T>): Unsubscribe
  unsubscribe(callback: Callback<T>): void
  fire(event: T): void
}

export type Unsubscribe = () => void

class EventSourceImpl<T> implements InternalEvent<T> {
  fire(event: T): void {
    for (const callback of this._callbacks) {
      callback(event)
    }
  }
  private readonly _callbacks: Callback<T>[] = []
  subscribe(callback: Callback<T>): Unsubscribe {
    if (this._callbacks.indexOf(callback) === -1) {
      this._callbacks.push(callback)
    }
    return () => this.unsubscribe(callback)
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
