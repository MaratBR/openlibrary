import { InputHTMLAttributes } from 'react'

export type CheckboxProps = InputHTMLAttributes<HTMLInputElement>

export default function Checkbox({ checked, className = '', children, ...props }: CheckboxProps) {
  return (
    <label className={`checkbox ${className}`}>
      {checked === false ? (
        <i className="fa-regular fa-square cursor-pointer" />
      ) : checked === true ? (
        <i className="fa-regular fa-square-check cursor-pointer" />
      ) : (
        <i className="fa-regular fa-square-minus cursor-pointer" />
      )}
      <input type="checkbox" checked={checked} {...props} />
      {children}
    </label>
  )
}
