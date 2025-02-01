export interface OLIsland {
  // eslint-disable-next-line no-unused-vars
  mount(el: HTMLElement, data: unknown): () => void
}

export function validateIslandName(name: string) {
  if (!name) throw new Error('Island name is not specified')
  if (!/^[a-zA-Z0-9_-]+$/.test(name)) throw new Error('Island name is invalid')
}
