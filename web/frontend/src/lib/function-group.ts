export class FunctionGroup {
  private readonly _fns: Array<() => void> = []

  add(fn: () => void) {
    this._fns.push(fn)
  }

  fire() {
    while (this._fns.length) {
      const fn = this._fns.pop()!
      fn()
    }
  }
}
