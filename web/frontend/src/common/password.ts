import { Schema } from 'effect'

export const passwordRequirementsSchema = Schema.Struct({
  Digits: Schema.Boolean,
  Symbols: Schema.String,
  SymbolsEnabled: Schema.Boolean,
  DifferentCases: Schema.Boolean,
  MinLength: Schema.Number,
})

export type PasswordRequirements = Schema.Schema.Type<typeof passwordRequirementsSchema>

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
