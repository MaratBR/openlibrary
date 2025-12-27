export function textSearch(query: string, target: string): boolean {
  return target.toLocaleLowerCase().includes(query.toLocaleLowerCase())
}
