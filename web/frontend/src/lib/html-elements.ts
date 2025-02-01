export function isAttrTrue(v: string | null) {
  if (v === '' || v === 'true') {
    return true
  }
  return false
}
