export function formatNumber(num: number, digits: number): string {
  if (digits <= 0) return `${Math.floor(num)}`
  const mult = Math.pow(10, digits)
  return `${Math.floor(num * mult) / mult}`
}

export function formatFileSize(size: number) {
  if (size === 0) return '0K'

  if (size >= 1_000_000_000) {
    return `${formatNumber(size / 1_000_000_000, 1)}G`
  }

  if (size >= 1_000_000) {
    return `${formatNumber(size / 1_000_000, 1)}M`
  }

  if (size >= 1_000) {
    return `${formatNumber(size / 1_000, 1)}K`
  }

  return `${size}B`
}
