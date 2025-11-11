import { z } from 'zod'

export const passwordRequirementsSchema = z.object({
  Digits: z.boolean(),
  Symbols: z.string(),
  SymbolsEnabled: z.boolean(),
  DifferentCases: z.boolean(),
  MinLength: z.number(),
})

export type PasswordRequirements = z.infer<typeof passwordRequirementsSchema>

export type PasswordValidationResult = {
  digits: boolean
  symbols: boolean
  differentCases: boolean
  minLength: boolean
  valid: boolean
}

export function validatePassword(
  password: string,
  rules: PasswordRequirements,
): PasswordValidationResult {
  const result: PasswordValidationResult = {
    digits: true,
    symbols: true,
    differentCases: true,
    minLength: true,
    valid: true,
  }

  if (rules.Digits) {
    result.digits = /\d/.test(password)
  }

  if (rules.DifferentCases) {
    result.differentCases = /[a-z]/.test(password) && /[A-Z]/.test(password)
  }

  if (rules.MinLength) {
    result.minLength = password.length >= rules.MinLength
  }

  if (rules.SymbolsEnabled) {
    result.symbols = rules.Symbols.split('').some((symbol) => password.includes(symbol))
  }

  result.valid = result.symbols && result.minLength && result.differentCases && result.digits
  return result
}
