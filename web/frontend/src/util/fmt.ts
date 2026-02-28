function getLocale(): string {
  return 'en-US' // TODO
}

export function formatNumber(num: number): string {
  return new Intl.NumberFormat(getLocale()).format(num)
}

export function formatNumberK(num: number): string {
  if (num < 100) {
    return formatNumber(num)
  }

  const k = num / 1000

  if (k < 10) {
    return `${formatNumber(Math.round(k * 10) / 10)}k`
  }

  if (k < 1000) {
    return `${formatNumber(Math.round(k))}k`
  }

  const M = k / 1000

  if (M < 10) {
    return `${formatNumber(Math.round(M * 10) / 10)}M`
  }

  return `${formatNumber(Math.round(M))}k`
}
