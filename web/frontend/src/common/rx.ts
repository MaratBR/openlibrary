import { useEffect, useState } from 'preact/hooks'

type Callback<T> = (event: T) => void

export type Dispose = () => void

export interface Subscribable<TEvent> {
  subscribe(callback: Callback<TEvent>): Dispose
}

export interface Writeable<T> {
  set(value: T): void
}

export interface WithValue<T> {
  get(): T
}

export class OLEvent<TEvent> implements Subscribable<TEvent>, Writeable<TEvent> {
  private readonly _callbacks: Callback<TEvent>[] = []

  subscribe(callback: Callback<TEvent>): Dispose {
    if (this._callbacks.indexOf(callback) === -1) {
      this._callbacks.push(callback)
    }

    return () => {
      const idx = this._callbacks.indexOf(callback)
      if (idx !== -1) {
        this._callbacks.splice(idx, 1)
      }
    }
  }

  set(event: TEvent) {
    for (const callback of this._callbacks) {
      callback(event)
    }
  }

  dispose() {
    this._callbacks.splice(0, this._callbacks.length)
  }
}

export class Subject<T> extends OLEvent<T> implements WithValue<T> {
  private _value: T

  constructor(value: T) {
    super()
    this._value = value

    this.subscribe((value) => {
      this._value = value
    })
  }

  get(): T {
    return this._value
  }
}

export function useSubject<T>(subject: Subscribable<T> & WithValue<T>): T {
  const [value, setValue] = useState<T>(subject.get())

  useEffect(() => {
    return subject.subscribe(setValue)
  }, [subject])

  return value
}

/* eslint-disable @typescript-eslint/no-explicit-any */
export class Derived<
    TList extends (Subscribable<any> & WithValue<any>)[], // tuple of Subscribables
    Result,
  >
  implements Subscribable<Result>, WithValue<Result>
{
  private values: { [K in keyof TList]: TList[K] extends Subscribable<infer E> ? E : never }
  private result: Result
  private subscribers = new Set<Callback<Result>>()
  private unsubscribers: Dispose[] = []

  constructor(
    private sources: [...TList],
    private compute: (
      ...args: { [K in keyof TList]: TList[K] extends Subscribable<infer E> ? E : never }
    ) => Result,
  ) {
    // initialize values with undefined assertions (we’ll fill them)
    this.values = new Array(this.sources.length) as any
    this.sources.forEach((source, index) => {
      this.values[index] = source.get()
      const unsub = source.subscribe((value: any) => {
        this.values[index] = value
        this.recompute()
      })
      this.unsubscribers.push(unsub)
    })

    // initial computation (in case all sources emit immediately)
    this.result = this.compute(...(this.values as any))
  }

  private recompute() {
    const newResult = this.compute(...(this.values as any))
    this.result = newResult
    this.subscribers.forEach((cb) => cb(newResult))
  }

  subscribe(callback: Callback<Result>): Dispose {
    this.subscribers.add(callback)
    callback(this.result) // fire immediately with current value
    return () => {
      this.subscribers.delete(callback)
    }
  }

  get(): Result {
    return this.result
  }

  disposeAll(): void {
    this.unsubscribers.forEach((dispose) => dispose())
    this.unsubscribers = []
    this.subscribers.clear()
  }
}
/* eslint-enable @typescript-eslint/no-explicit-any */

export class State<TObject extends object> extends Subject<TObject> {
  updateState(update: Partial<TObject>) {
    this.set({
      ...this.get(),
      ...update,
    })
  }
}
