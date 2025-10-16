import { InputHTMLAttributes } from 'preact/compat'

export type CheckboxProps = InputHTMLAttributes<HTMLInputElement>

export default function Checkbox({
  checked,
  class: class_ = '',
  className,
  children,
  ...props
}: CheckboxProps) {
  return (
    <label class={`checkbox ${class_}`} className={className}>
      {checked === false ? (
        <i class="fa-regular fa-square cursor-pointer" />
      ) : checked === true ? (
        <i class="fa-regular fa-square-check cursor-pointer" />
      ) : (
        <i class="fa-regular fa-square-minus cursor-pointer" />
      )}
      <input type="checkbox" checked={checked} {...props} />
      {children}
    </label>
  )
}
