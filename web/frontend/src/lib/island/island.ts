export interface OLIslandMounted {
  dispose(): void
  setData(data: unknown): void
}

export interface OLIsland {
  mount(el: HTMLElement, data: unknown): OLIslandMounted
}

export function validateIslandName(name: string) {
  if (!name) throw new Error('Island name is not specified')
  if (!/^[a-zA-Z0-9_-]+$/.test(name)) throw new Error('Island name is invalid')
}
