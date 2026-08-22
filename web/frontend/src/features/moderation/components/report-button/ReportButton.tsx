import type { ReactIslandProps } from '@/islands/common/react-island'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components'
import { FormEvent, MouseEvent, useEffect, useRef, useState } from 'react'

export type ReportButtonConfig = {
  targetType: 'user' | 'book' | 'comment'
  targetId: string
  chapterId?: number
  captureSelection?: boolean
  buttonClass?: string
  label: string
  iconOnly?: boolean
  authenticated: boolean
  triggerId?: string
  labels: {
    title: string
    reason: string
    description: string
    submit: string
    cancel: string
    success: string
  }
}

export function ReportButtonIsland({ data }: ReactIslandProps) {
  const config = data as ReportButtonConfig
  const [open, setOpen] = useState(false)
  const excerpt = useRef('')

  useEffect(() => {
    const trigger = config.triggerId ? document.getElementById(config.triggerId) : null
    if (!trigger) return
    const show = () => {
      if (!config.authenticated) {
        location.href = `/login?next=${encodeURIComponent(location.pathname + location.search)}`
        return
      }
      excerpt.current = config.captureSelection
        ? window.getSelection()?.toString().trim().slice(0, 2000) || ''
        : ''
      setOpen(true)
    }
    trigger.addEventListener('click', show)
    return () => trigger.removeEventListener('click', show)
  }, [config])

  return (
    <ReportPopup
      config={config}
      open={open}
      excerpt={excerpt.current}
      onClose={() => setOpen(false)}
    />
  )
}

export function ReportPopup({
  config,
  open,
  excerpt,
  onClose,
}: {
  config: ReportButtonConfig
  open: boolean
  excerpt: string
  onClose: () => void
}) {
  const [submitting, setSubmitting] = useState(false)
  const [reasons, setReasons] = useState<string[]>([])
  const [reason, setReason] = useState('')
  const [dialog, setDialog] = useState<HTMLDialogElement | null>(null)

  useEffect(() => {
    if (open && dialog && !dialog.open) dialog.showModal()
  }, [open, dialog])
  useEffect(() => {
    if (!open) return
    let active = true
    fetch(`/_api/reports/reasons?targetType=${encodeURIComponent(config.targetType)}`)
      .then((response) => window.OLAPIResponse.create<string[]>(response))
      .then((response) => {
        response.throwIfError()
        if (active) {
          setReasons(response.data)
          setReason(response.data[0] || '')
        }
      })
      .catch((error) => window.toast.error(error))
    return () => {
      active = false
    }
  }, [open, config.targetType])

  function close() {
    dialog?.close()
    onClose()
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const form = event.currentTarget
    const data = new FormData(form)
    setSubmitting(true)
    try {
      const response = await fetch('/_api/reports', {
        method: 'POST',
        headers: { 'content-type': 'application/json' },
        body: JSON.stringify({
          targetType: config.targetType,
          targetId: config.targetId,
          reason,
          description: data.get('description') || '',
          chapterId: config.chapterId ? String(config.chapterId) : null,
          excerpt,
        }),
      })
      const result = await window.OLAPIResponse.create(response)
      result.throwIfError()
      form.reset()
      close()
      window.toast({ type: 'success', text: config.labels.success })
    } catch (error) {
      window.toast.error(error)
    } finally {
      setSubmitting(false)
    }
  }

  function backdropClick(event: MouseEvent<HTMLDialogElement>) {
    if (event.target === event.currentTarget) close()
  }

  return open ? (
    <dialog ref={setDialog} className="modal" onCancel={close} onClick={backdropClick}>
      <form className="modal__content grid gap-4" onSubmit={submit}>
        <div className="flex items-center justify-between gap-4">
          <h2 className="font-title text-2xl font-semibold">{config.labels.title}</h2>
          <button
            type="button"
            className="btn btn--icon btn--ghost"
            aria-label={config.labels.cancel}
            onClick={close}
          >
            <i className="fa-solid fa-xmark" />
          </button>
        </div>
        <label className="grid gap-2">
          <span className="font-medium">{config.labels.reason}</span>
          <Select value={reason} onValueChange={setReason}>
            <SelectTrigger className="w-full">
              <SelectValue />
            </SelectTrigger>
            <SelectContent portalContainer={dialog}>
              {reasons.map((value) => (
                <SelectItem key={value} value={value}>
                  {value}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </label>
        <label className="grid gap-2">
          <span className="font-medium">{config.labels.description}</span>
          <textarea className="textarea min-h-28" name="description" maxLength={2000} />
        </label>
        <div className="flex justify-end gap-2">
          <button type="button" className="btn btn--outline" onClick={close}>
            {config.labels.cancel}
          </button>
          <button className="btn btn--primary" type="submit" disabled={submitting || !reason}>
            {config.labels.submit}
          </button>
        </div>
      </form>
    </dialog>
  ) : null
}
