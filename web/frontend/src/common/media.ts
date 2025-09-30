import { Dispose, Subject, Subscribable, WithValue } from './rx'

export class MatchesMedia implements Subscribable<boolean>, WithValue<boolean> {
  private _subject: Subject<boolean>
  private _mql?: MediaQueryList

  constructor(request: string, fallback: boolean) {
    if (!window.matchMedia) {
      this._subject = new Subject(fallback)
      return
    }

    this._mql = window.matchMedia(request)
    this._subject = new Subject(this._mql.matches)
    this._mql.addEventListener('change', (ev) => {
      this._subject.set(ev.matches)
    })
  }
  subscribe(callback: (event: boolean) => void): Dispose {
    return this._subject.subscribe(callback)
  }
  getValue(): boolean {
    return this._subject.getValue()
  }
}

export const PREFERS_DARK_MODE = new MatchesMedia('(prefers-color-scheme: dark)', false)
