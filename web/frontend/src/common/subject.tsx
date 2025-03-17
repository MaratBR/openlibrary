import { useEffect, useState } from 'preact/hooks'
import { EventLike, OLEvent } from './event'

export class Subject<T> implements EventLike<T> {
  value: T
  private _updatedEvent = new OLEvent<T>()

  constructor(value: T) {
    this.value = value
  }

  subscribe(callback: (event: T) => void): void {
    this._updatedEvent.subscribe(callback)
  }

  unsubscribe(callback: (event: T) => void): void {
    this._updatedEvent.unsubscribe(callback)
  }

  set(value: T) {
    this.value = value
    this._updatedEvent.publish(value)
  }
}

export function useSubject<T>(subject: Subject<T>): T {
  const [value, setValue] = useState<T>(subject.value)

  useEffect(() => {
    subject.subscribe(setValue)
    return () => subject.unsubscribe(setValue)
  }, [])

  return value
}
