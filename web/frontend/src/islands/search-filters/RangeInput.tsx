type NumberRange = {
  max: number | null
  min: number | null
}

export type RangeInputProps = {
  value: NumberRange
  // eslint-disable-next-line no-unused-vars
  onInput: (value: NumberRange) => void

  disableNegative?: boolean
}

function getValue(event: InputEvent): number | null | undefined {
  if (event.target instanceof HTMLInputElement) {
    const value = event.target.value
    if (value === '') return null
    const numValue = +value
    if (Number.isNaN(numValue)) return undefined
    return numValue
  }
  return undefined
}

function normalizeValue(
  value: number | null,
  options: { disableNegative: boolean },
): number | null {
  if (value === null) return value

  if (options.disableNegative && value < 0) {
    return 0
  }

  return value
}

export default function RangeInput({ value, onInput, disableNegative = false }: RangeInputProps) {
  function handleKeyDown(event: KeyboardEvent) {
    if (
      (event.key === '-' && disableNegative) ||
      event.key === '.' ||
      event.key === 'e' ||
      event.key === 'E'
    ) {
      event.preventDefault()
    }
  }

  function handleMinChange(event: InputEvent) {
    const min = getValue(event)
    if (min !== undefined) {
      onInput({
        max: value.max,
        min: normalizeValue(min, { disableNegative }),
      })
    }
  }

  function handleMaxChange(event: InputEvent) {
    const max = getValue(event)
    if (max !== undefined) {
      onInput({
        min: value.min,
        max: normalizeValue(max, { disableNegative }),
      })
    }
  }

  return (
    <div class="grid grid-cols-2 max-w-80">
      <input
        class="input rounded-r-none border-r-0 hover:border-primary hover:ring-1 hover:ring-primary !outline-none transition-all"
        type="number"
        value={value.min ?? ''}
        onInput={handleMinChange}
        onKeyDown={handleKeyDown}
      />
      <input
        class="input rounded-l-none hover:border-primary hover:ring-1 hover:ring-primary !outline-none transition-all"
        type="number"
        value={value.max ?? ''}
        onInput={handleMaxChange}
        onKeyDown={handleKeyDown}
      />
    </div>
  )
}
