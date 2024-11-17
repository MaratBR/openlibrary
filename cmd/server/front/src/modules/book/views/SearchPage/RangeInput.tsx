import { SliderMulti } from '@/components/slider-dual'
import { NumberRange } from './state'

export function RangeInput({
  value,
  range,
  onChange,
}: {
  value: NumberRange | null
  range: NumberRange
  onChange: (value: NumberRange) => void
}) {
  const valMin = value?.min ?? range.min ?? 0
  const valMax = value?.max ?? range.max ?? 0

  const adjMin = Math.min(valMin, range.min ?? 0)
  const adjMax = Math.max(valMax, range.max ?? 0)

  let step: number
  if (valMax - valMin <= 0) {
    step = 1
  } else {
    step = Math.max(1, Math.floor((valMax - valMin) / 100))
  }

  function handleMinChange(event: React.ChangeEvent<HTMLInputElement>) {
    onChange(updateMin(event.target.valueAsNumber, range))
  }

  function handleMaxChange(event: React.ChangeEvent<HTMLInputElement>) {
    onChange(updateMax(event.target.valueAsNumber, range))
  }

  return (
    <div className="range">
      <div className="pt-10 pb-1">
        <SliderMulti
          min={adjMin}
          max={adjMax}
          minStepsBetweenThumbs={1}
          value={[Number.isNaN(valMin) ? 0 : valMin, Number.isNaN(valMax) ? 0 : valMax]}
          step={step}
          onValueChange={(range) => onChange({ min: range[0], max: range[1] })}
        />
      </div>

      <div className="grid grid-cols-2">
        <input type="number" className="range__input" value={valMin} onChange={handleMinChange} />
        <input type="number" className="range__input" value={valMax} onChange={handleMaxChange} />
      </div>
    </div>
  )
}

function updateMin(min: number, range: NumberRange): NumberRange {
  min = nonNegativeInt32(min)
  const newRange = {
    min,
    max: range.max,
  }

  if (newRange.max !== null && newRange.max < newRange.min) {
    newRange.max = newRange.min
  }

  return newRange
}

function updateMax(max: number, range: NumberRange): NumberRange {
  max = nonNegativeInt32(max)
  const newRange = {
    max,
    min: range.min,
  }

  if (newRange.min !== null && newRange.min > newRange.max) {
    newRange.min = newRange.max
  }

  return newRange
}

function nonNegativeInt32(v: number): number {
  v = Math.round(v)
  if (v < 0) return 0
  if (v > 2147483647) return 2147483647
  return v
}
