import {
  PasswordRequirements,
  passwordRequirementsSchema,
  PasswordValidationResult,
  validatePassword,
} from '@/common/password'
import { Subject } from '@/common/rx'
import Alpine from 'alpinejs'

Alpine.data('passwordInput', () => ({
  id: `passwordInputRX_${Math.random().toString(16).substring(2)}`,
  requirements: null as PasswordRequirements | null,

  init() {
    this.requirements = passwordRequirementsSchema.parse(
      JSON.parse(this.$el.getAttribute('data-password-requirements') || ''),
    )
    const rx = new Subject<PasswordValidationResult>({
      valid: false,
      symbols: false,
      minLength: false,
      differentCases: false,
      digits: false,
    })
    this.$store[this.id] = rx
  },

  destroy() {},

  Input: {
    'x-init'() {
      if (!(this.$el instanceof HTMLInputElement))
        throw new Error('x-bind="Input" must be set on an input element')
      this.$el.addEventListener('input', (event: Event) => {
        if (!(event.target instanceof HTMLInputElement)) return
        const value = event.target.value

        const requirements = this.requirements
        if (!requirements) return

        const subject = this.$store[this.id] as Subject<PasswordValidationResult> | undefined
        if (!subject) return
        const validationResult = validatePassword(value, requirements)
        subject.set(validationResult)
      })
    },
  },

  Requirement: {
    'x-init'() {
      const requirement = this.$el.dataset.requirement
      const subject = this.$store[this.id] as Subject<PasswordValidationResult> | undefined
      if (!subject) return

      const validate = (result: PasswordValidationResult) => {
        if (!this.requirements) return
        if (typeof requirement !== 'string') return

        const valid = isValid(requirement, result)
        this.$el.setAttribute('aria-valid', `${valid}`)
      }
      validate(subject.get())
      subject.subscribe(validate)
    },
  },
}))

function isValid(type: string, validationResult: PasswordValidationResult) {
  switch (type) {
    case 'Symbols':
      return validationResult.symbols
    case 'Digits':
      return validationResult.digits
    case 'DifferentCases':
      return validationResult.differentCases
    case 'MinLength':
      return validationResult.minLength
    default:
      return false
  }
}
