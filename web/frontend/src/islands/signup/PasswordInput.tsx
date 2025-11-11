import { PasswordRequirements, PasswordValidationResult, validatePassword } from '@/common/password'
import { ChangeEvent } from 'preact/compat'
import { useMemo } from 'preact/hooks'

interface PasswordInputProps {
  passwordRequirements: PasswordRequirements
  name: string
  onPasswordChange: (value: string) => void
  value: string
  id?: string
}

export function PasswordInput({
  passwordRequirements,
  name,
  id,
  onPasswordChange,
  value,
}: PasswordInputProps) {
  const validation = useMemo(() => validatePassword(value, passwordRequirements), [value])

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const value = (e.target as HTMLInputElement).value
    onPasswordChange(value)
  }

  const renderRequirement = (type: keyof PasswordValidationResult, label: string) => {
    const valid = isValid(type, validation)
    return (
      <li aria-valid={valid}>
        <i className="password-input__ok-icon fa-solid fa-circle-check" />
        <i className="password-input__err-icon fa-solid fa-circle-xmark" />
        {label}
      </li>
    )
  }

  return (
    <div className="contents">
      <input
        required
        type="password"
        name={name}
        className="input"
        onChange={handleChange}
        id={id}
      />
      <ul className="password-input__requirements">
        {renderRequirement('symbols', window._('password.mustHaveSymbol'))}
        {renderRequirement('differentCases', window._('password.mustHaveCases'))}
        {renderRequirement('digits', window._('password.mustHaveDigits'))}
        {renderRequirement(
          'minLength',
          window._('password.minLength', { characters: `${passwordRequirements.MinLength}` }),
        )}
      </ul>
    </div>
  )
}

function isValid(type: keyof PasswordValidationResult, result: PasswordValidationResult) {
  switch (type) {
    case 'symbols':
      return result.symbols
    case 'digits':
      return result.digits
    case 'differentCases':
      return result.differentCases
    case 'minLength':
      return result.minLength
    default:
      return false
  }
}
