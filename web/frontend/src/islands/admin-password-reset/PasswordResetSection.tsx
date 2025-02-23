import { useEffect, useMemo, useState } from 'preact/hooks'
import { PreactIslandProps } from '../common'
import { evaluatePasswordStrength, generateStrongPassword } from '@/lib/password'

import clsx from 'clsx'

export default function PasswordReset({ rootElement }: PreactIslandProps) {
  const [value, setValue] = useState('')

  useEffect(() => {
    setValue(generateStrongPassword())

    const callback = () => {
      setValue(generateStrongPassword())
    }

    rootElement.addEventListener('island:custom:generate-password', callback)

    return () => {
      rootElement.removeEventListener('island:custom:generate-password', callback)
    }
  }, [rootElement])

  const pwdEval = useMemo(() => evaluatePasswordStrength(value), [value])

  function handleCancel() {
    rootElement.dispatchEvent(new CustomEvent('island:request-destroy'))
  }

  return (
    <div class="flex">
      <div class="relative mb-10 mr-2">
        <input
          name="password"
          class="input"
          onInput={(e) => setValue((e.target as HTMLInputElement).value)}
          value={value}
        />
        <div
          class={clsx('absolute top-full w-full left-0 py-2 flex justify-center', {
            'text-red-800 bg-red-200 border-red-600 border': pwdEval.strength === 'Weak',
            'text-green-800 bg-green-200 border-green-600 border':
              pwdEval.strength === 'Strong' || pwdEval.strength === 'VeryStrong',
            'text-yellow-800 bg-yellow-200 border-yellow-600 border': pwdEval.strength === 'OK',
          })}
        >
          {window._(`passwordStrength.${pwdEval.strength}`)}
        </div>
      </div>
      <button onClick={handleCancel} class="ol-btn ol-btn--secondary ol-btn-sm">
        {window._('common.cancel')}
      </button>
    </div>
  )
}
