import { z } from 'zod'
import { ReactIslandProps } from '../common/react-island'
import { passwordRequirementsSchema, validatePassword } from '@/common/password'
import { useLayoutEffect, useMemo, useRef, useState } from 'react'
import { PasswordInput } from './PasswordInput'
import { animate } from 'popmotion'

const signUpFormDataSchema = z.object({
  PasswordRequirements: passwordRequirementsSchema,
  PrefilledUsername: z.string(),
  PrefilledEmail: z.string(),
})

function isValidUsername(username: string) {
  return username.trim().length >= 2
}

function isValidEmail(value: string) {
  return /.+@.+/.test(value)
}

function normalizeUsername(value: string): string {
  return value.replace(/[^0-9a-zA-Z_]/g, '')
}

export default function SignUpForm({ data: dataParam }: ReactIslandProps) {
  const { PasswordRequirements } = useMemo(() => signUpFormDataSchema.parse(dataParam), [dataParam])

  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [repeatPassword, setRepeatPassword] = useState('')
  const [submitting, setSubmitting] = useState(false)

  const isValid = useMemo(
    () =>
      isValidUsername(username) &&
      validatePassword(password, PasswordRequirements).valid &&
      password === repeatPassword &&
      password.length > 0 &&
      isValidEmail(email),
    [username, password, PasswordRequirements, repeatPassword, email],
  )

  const rootElementRef = useRef<HTMLDivElement | null>(null)

  useLayoutEffect(() => {
    const { current: rootElement } = rootElementRef
    if (!rootElement) return
    rootElement.style.overflow = 'hidden'
    rootElement.style.maxHeight = '0px'

    animate({
      duration: 120,
      from: 0,
      to: 500,
      onUpdate: (latest) => {
        rootElement.style.maxHeight = `${latest}px`
        rootElement.style.opacity = `${latest / 500}`
      },
      onComplete() {
        rootElement.style.overflow = ''
        rootElement.style.opacity = ''
        rootElement.style.maxHeight = ''
      },
    })
  }, [])

  return (
    <div ref={rootElementRef} className="overflow-auto" style={{ scrollbarWidth: 'thin' }}>
      <div className="form-control">
        <div className="form-control__label">
          <label htmlFor="username" className="label">
            {window._('login.username')}
          </label>
        </div>
        <div className="form-control__value">
          <input
            required
            name="username"
            id="username"
            className="input"
            type="text"
            value={username}
            onChange={(e) => setUsername(normalizeUsername((e.target as HTMLInputElement).value))}
          />
        </div>
      </div>

      <div className="form-control">
        <div className="form-control__label">
          <label htmlFor="email" className="label">
            {window._('login.email')}
          </label>
        </div>
        <div className="form-control__value">
          <input
            required
            name="email"
            id="email"
            className="input"
            type="email"
            value={email}
            onChange={(e) => setEmail((e.target as HTMLInputElement).value)}
          />
        </div>
      </div>

      <div className="form-control">
        <div className="form-control__label">
          <label htmlFor="password" className="label">
            {window._('login.password')}
          </label>
        </div>
        <div className="form-control__value">
          <PasswordInput
            id="password"
            name="password"
            passwordRequirements={PasswordRequirements}
            onPasswordChange={setPassword}
            value={password}
          />
        </div>
      </div>

      <div className="form-control">
        <div className="form-control__label">
          <label htmlFor="repeatPassword" className="label">
            {window._('signup.repeatPassword')}
          </label>
        </div>
        <div className="form-control__value">
          <input
            required
            onChange={(e) => setRepeatPassword((e.target as HTMLInputElement).value)}
            name="repeatPassword"
            id="repeatPassword"
            className="input"
            type="password"
            value={repeatPassword}
          />
        </div>
      </div>

      <div className="flex gap-2 justify-center">
        <button
          disabled={!isValid}
          type="submit"
          onClick={() => {
            setSubmitting(true)
          }}
          className="btn btn--lg btn--default"
        >
          {submitting ? <span className="loader" /> : window._('signup.signup')}
        </button>
      </div>
    </div>
  )
}
