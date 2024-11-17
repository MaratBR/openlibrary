import React from 'react'

export interface IEvent<T> {
  addListener(listener: (value: T) => void): void
  removeListener(listener: (value: T) => void): void
  emit(value: T): void
}

export class Event<T> implements IEvent<T> {
  private listeners: ((value: T) => void)[] = []

  addListener(listener: (value: T) => void) {
    this.listeners.push(listener)
  }

  removeListener(listener: (value: T) => void) {
    this.listeners = this.listeners.filter((l) => l !== listener)
  }

  emit(value: T) {
    this.listeners.forEach((l) => l(value))
  }
}

export interface ISubject<T> extends IEvent<T> {
  value: T
  set(value: T): void
}

export class Subject<T> extends Event<T> implements ISubject<T> {
  private event = new Event<T>()
  private _value: T

  get value() {
    return this._value
  }

  constructor(value: T) {
    super()
    this._value = value
  }

  set(value: T) {
    if (this._value === value) return
    this._value = value
    this.event.emit(value)
  }
}

export function useSubject<T>(subject: ISubject<T>) {
  const [value, setValue] = React.useState(subject.value)
  React.useEffect(() => {
    subject.addListener(setValue)
    return () => {
      subject.removeListener(setValue)
    }
  }, [subject])
  return value
}

type InferSubjectsResults<TSubjects extends ISubject<unknown>[]> = {
  [K in keyof TSubjects]: TSubjects[K] extends ISubject<infer T> ? T : never
}

export class Derived<TSubjects extends ISubject<unknown>[]> {
  constructor(...subjects: TSubjects) {}
}
