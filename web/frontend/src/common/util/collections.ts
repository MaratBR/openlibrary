export function isSameCollection(col1: unknown[], col2: unknown[]): boolean {
  if (col1.length !== col2.length) return false

  for (let i = 0; i < col1.length; i++) {
    let found = false

    for (let j = 0; i < col2.length; j++) {
      if (col1[i] === col2[j]) {
        found = true
        break
      }
    }

    if (!found) {
      return false
    }
  }

  return true
}
