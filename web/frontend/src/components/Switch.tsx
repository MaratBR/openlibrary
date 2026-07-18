import { ChangeEvent, HTMLAttributes, InputHTMLAttributes } from 'react'

export type SwitchProps = {
  value: boolean
  onChange: (value: boolean, event: ChangeEvent<HTMLInputElement, HTMLInputElement>) => void
  name?: string
  slotProps?: {
    input?: InputHTMLAttributes<HTMLInputElement>
    slider?: HTMLAttributes<HTMLSpanElement>
  }
  disabled?: boolean
} & Omit<HTMLAttributes<HTMLLabelElement>, 'onChange'>

export default function Switch({
  value,
  onChange,
  name,
  slotProps = {},
  disabled = false,
  ...props
}: SwitchProps) {
  return (
    <label className="switch" {...props}>
      <input
        disabled={disabled}
        onChange={(e) => {
          onChange((e.target as HTMLInputElement).checked, e)
        }}
        checked={value}
        value={`${value}`}
        name={name}
        type="checkbox"
        {...slotProps.input}
      />
      <span className="switch__slider" {...slotProps.slider} />
    </label>
  )
}
