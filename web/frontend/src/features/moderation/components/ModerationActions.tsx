import { FormControl } from '@/components/FormControl'
import Modal from '@/components/Modal'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components'
import { FormEvent, ReactNode, useEffect, useState } from 'react'
import { useRevalidator } from 'react-router'
import clsx from 'clsx'

export type ModerationConfirmation = {
  title: string
  description: string
  destructive?: boolean
  run: () => Promise<unknown>
  onSuccess?: () => void
}

export type ModerationActionOption = { value: string; label: string }

export function ModerationActionSelector({
  value,
  onValueChange,
  options,
  label,
  id,
  className,
}: {
  value: string
  onValueChange: (value: string) => void
  options: ModerationActionOption[]
  label: string
  id: string
  className?: string
}) {
  return (
    <section className={clsx('Card', className)}>
      <FormControl label={label} htmlFor={id}>
        <Select value={value} onValueChange={onValueChange}>
          <SelectTrigger id={id} className="w-full">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            {options.map((option) => (
              <SelectItem key={option.value} value={option.value}>
                {option.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </FormControl>
    </section>
  )
}

export function ModerationConfirmationDialog({
  confirmation,
  pending,
  onClose,
  onConfirm,
  confirmLabel,
}: {
  confirmation?: ModerationConfirmation
  pending: boolean
  onClose: () => void
  onConfirm: () => void
  confirmLabel: string
}) {
  return (
    <Modal open={Boolean(confirmation)} onClose={() => !pending && onClose()}>
      <div className="max-w-128">
        <h2 className="text-xl font-semibold">{confirmation?.title}</h2>
        <p className="my-3 text-secondary-foreground">{confirmation?.description}</p>
        <div className="flex justify-end gap-2 mt-5">
          <button disabled={pending} className="Btn Btn--outline" onClick={onClose}>
            {window._('common.cancel')}
          </button>
          <button
            disabled={pending}
            className={`Btn ${confirmation?.destructive ? 'Btn--destructive' : 'Btn--primary'}`}
            onClick={onConfirm}
          >
            {pending && <span className="circle-loader mr-2" />}
            {confirmLabel}
          </button>
        </div>
      </div>
    </Modal>
  )
}

export function ModerationActionsPanel({
  action,
  onActionChange,
  options,
  selectLabel,
  selectId,
  confirmation,
  onConfirmationChange,
  successMessage,
  confirmLabel,
  children,
}: {
  action: string
  onActionChange: (value: string) => void
  options: ModerationActionOption[]
  selectLabel: string
  selectId: string
  confirmation?: ModerationConfirmation
  onConfirmationChange: (value?: ModerationConfirmation) => void
  successMessage: string
  confirmLabel: string
  children: ReactNode
}) {
  const revalidator = useRevalidator()
  const [pending, setPending] = useState(false)
  const confirm = async () => {
    if (!confirmation) return
    setPending(true)
    try {
      const response = (await confirmation.run()) as { throwIfError?: () => void }
      response.throwIfError?.()
      confirmation.onSuccess?.()
      window.toast({ title: successMessage })
      onConfirmationChange(undefined)
      await revalidator.revalidate()
    } catch (error) {
      window.toast.error(error)
    } finally {
      setPending(false)
    }
  }
  return (
    <>
      <div className="grid lg:grid-cols-[18rem_minmax(0,1fr)] gap-5 items-start">
        <ModerationActionSelector
          value={action}
          onValueChange={onActionChange}
          options={options}
          label={selectLabel}
          id={selectId}
        />
        <div>{children}</div>
      </div>
      <ModerationConfirmationDialog
        confirmation={confirmation}
        pending={pending}
        onClose={() => onConfirmationChange(undefined)}
        onConfirm={() => void confirm()}
        confirmLabel={confirmLabel}
      />
    </>
  )
}

export function ModerationActionCard({
  title,
  description,
  children,
  className,
}: {
  title: string
  description: string
  children: ReactNode
  className?: string
}) {
  return (
    <section className={clsx('Card', className)}>
      <h2 className="text-xl font-semibold">{title}</h2>
      <p className="text-secondary-foreground mt-1 mb-5">{description}</p>
      {children}
    </section>
  )
}

export function ModerationReasonActionCard({
  title,
  description,
  submitLabel,
  reasonLabel,
  destructive = false,
  onSubmit,
}: {
  title: string
  description: string
  submitLabel: string
  reasonLabel: string
  destructive?: boolean
  onSubmit: (reason: string, onSuccess: () => void) => void
}) {
  const [reason, setReason] = useState('')
  const submit = (event: FormEvent) => {
    event.preventDefault()
    if (reason.trim()) onSubmit(reason.trim(), () => setReason(''))
  }
  return (
    <ModerationActionCard title={title} description={description}>
      <form onSubmit={submit}>
        <FormControl label={reasonLabel}>
          <textarea
            className="input min-h-24"
            required
            value={reason}
            onChange={(event) => setReason(event.target.value)}
          />
        </FormControl>
        <button className={`Btn mt-4 ${destructive ? 'Btn--destructive' : 'Btn--primary'}`}>
          {submitLabel}
        </button>
      </form>
    </ModerationActionCard>
  )
}

export function ModerationValueActionCard({
  title,
  description,
  initialValue,
  multiline,
  valueLabel,
  reasonLabel,
  submitLabel,
  control,
  onSubmit,
}: {
  title: string
  description: string
  initialValue: string
  multiline?: boolean
  valueLabel: string
  reasonLabel: string
  submitLabel: string
  control?: (value: string, setValue: (value: string) => void) => ReactNode
  onSubmit: (value: string, reason: string, onSuccess: () => void) => void
}) {
  const [value, setValue] = useState(initialValue)
  const [reason, setReason] = useState('')
  useEffect(() => setValue(initialValue), [initialValue])
  const submit = (event: FormEvent) => {
    event.preventDefault()
    if (value.trim() && reason.trim()) onSubmit(value.trim(), reason.trim(), () => setReason(''))
  }
  const valueControl =
    control?.(value, setValue) ??
    (multiline ? (
      <textarea
        className="input min-h-32"
        required
        value={value}
        onChange={(event) => setValue(event.target.value)}
      />
    ) : (
      <input
        className="input"
        required
        value={value}
        onChange={(event) => setValue(event.target.value)}
      />
    ))
  return (
    <ModerationActionCard title={title} description={description}>
      <form onSubmit={submit} className="grid gap-4">
        <FormControl label={valueLabel}>{valueControl}</FormControl>
        <FormControl label={reasonLabel}>
          <textarea
            className="input min-h-20"
            required
            value={reason}
            onChange={(event) => setReason(event.target.value)}
          />
        </FormControl>
        <button className="Btn Btn--primary justify-self-start">{submitLabel}</button>
      </form>
    </ModerationActionCard>
  )
}

export const ModerationActions = {
  Root: ModerationActionsPanel,
  Selector: ModerationActionSelector,
  Card: ModerationActionCard,
  ConfirmationDialog: ModerationConfirmationDialog,
  ReasonCard: ModerationReasonActionCard,
  ValueCard: ModerationValueActionCard,
}
