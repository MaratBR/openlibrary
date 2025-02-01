export function generateStrongPassword(length = 18) {
  const upperCase = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ'
  const lowerCase = 'abcdefghijklmnopqrstuvwxyz'
  const numbers = '0123456789'
  const specialChars = '!@#$%^&*()_+[]{}|;:,.<>?'
  const allChars = upperCase + lowerCase + numbers + specialChars

  if (length < 8) {
    throw new Error('Password length should be at least 8 characters for better security.')
  }

  // Ensure at least one character from each category
  const password = [
    upperCase[Math.floor(Math.random() * upperCase.length)],
    lowerCase[Math.floor(Math.random() * lowerCase.length)],
    numbers[Math.floor(Math.random() * numbers.length)],
    specialChars[Math.floor(Math.random() * specialChars.length)],
  ]

  // Fill the rest of the password with random characters from all categories
  for (let i = password.length; i < length; i++) {
    password.push(allChars[Math.floor(Math.random() * allChars.length)])
  }

  // Shuffle the password to ensure randomness
  return password.sort(() => Math.random() - 0.5).join('')
}

export type PasswordStrength = 'Weak' | 'OK' | 'Strong' | 'VeryStrong'

export function evaluatePasswordStrength(password: string): {
  score: number
  strength: PasswordStrength
} {
  let score = 0

  // Criteria 1: Password length
  if (password.length < 8) score -= 6 // punish short passwords
  if (password.length >= 12) score += 2 // Bonus for longer passwords

  // Criteria 2: Character variety
  if (/[A-Z]/.test(password)) score += 2 // Uppercase letters
  if (/[a-z]/.test(password)) score += 2 // Lowercase letters
  if (/[0-9]/.test(password)) score += 2 // Numbers
  if (/[^A-Za-z0-9]/.test(password)) score += 2 // Special characters

  // Criteria 3: Entropy (bonus for less predictability)
  if (!/(.)\1{2,}/.test(password)) score += 1 // No repeated characters like "aaa"
  if (!/^[A-Za-z0-9]+$/.test(password)) score += 1 // No purely alphanumeric strings

  // Cap score to a max of 10
  score = Math.min(score, 10)

  // Determine human-readable strength
  let strength: PasswordStrength = 'Weak'
  if (score <= 4) strength = 'Weak'
  else if (score <= 6) strength = 'OK'
  else if (score <= 8) strength = 'Strong'
  else strength = 'VeryStrong'

  return { score, strength }
}
