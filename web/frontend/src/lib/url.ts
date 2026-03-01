export function getPage(url: string): number {
  const u = new URL(url)
  const pageStr = u.searchParams.get('page')
  if (!pageStr) return 1
  const page = Math.floor(+pageStr)
  if (Number.isNaN(page) || page < 1) return 1
  return page
}
