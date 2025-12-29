import { HTMLAttributes } from 'preact'

export type EditorToggleButtonProps = {
  active: boolean
} & HTMLAttributes<HTMLButtonElement>

export default function EditorToggleButton({ active, ...props }: EditorToggleButtonProps) {
  return (
    <button class="be-listitem be-toggle" aria-selected={active ? 'true' : 'false'} {...props} />
  )
}
