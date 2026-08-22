import { HTMLAttributes } from 'react'

export type EditorToggleButtonProps = {
  active: boolean
} & HTMLAttributes<HTMLButtonElement>

export default function EditorToggleButton({ active, ...props }: EditorToggleButtonProps) {
  return (
    <button
      className="be-listitem be-listitem--toggle"
      aria-selected={active ? 'true' : 'false'}
      {...props}
    />
  )
}
