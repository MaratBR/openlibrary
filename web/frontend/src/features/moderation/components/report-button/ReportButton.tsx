import type { ReactIslandProps } from '@/islands/common/react-island'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components'
import { FormEvent, MouseEvent, useEffect, useId, useRef, useState } from 'react'
import './report-dialog.scss'

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
    excerpt: string
    selectionHint: string
    submitting: string
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
  const [excerpt, setExcerpt] = useState('')
  const pendingSelection = useRef<string | null>(null)

  useEffect(() => {
    const trigger = config.triggerId ? document.getElementById(config.triggerId) : null
    if (!trigger) return
    const readSelection = () => window.getSelection()?.toString().trim().slice(0, 2000) || ''
    // Pointer focus can collapse the selection before the click is dispatched.
    const capture = () => {
      pendingSelection.current = config.captureSelection ? readSelection() : ''
    }
    const show = () => {
      if (!config.authenticated) {
        location.href = `/login?next=${encodeURIComponent(location.pathname + location.search)}`
        return
      }
      setExcerpt(config.captureSelection ? (pendingSelection.current ?? readSelection()) : '')
      pendingSelection.current = null
      setOpen(true)
    }
    trigger.addEventListener('pointerdown', capture)
    trigger.addEventListener('click', show)
    return () => {
      trigger.removeEventListener('pointerdown', capture)
      trigger.removeEventListener('click', show)
    }
  }, [config])

  return (
    <ReportPopup config={config} open={open} excerpt={excerpt} onClose={() => setOpen(false)} />
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
  const titleId = useId()
  const reasonId = useId()
  const excerptId = useId()
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
    <dialog
      ref={setDialog}
      className="Dialog ReportDialog"
      aria-labelledby={titleId}
      onCancel={close}
      onClick={backdropClick}
    >
      <form className="ReportDialog-form" onSubmit={submit} aria-busy={submitting}>
        <header className="ReportDialog-header">
          <h2 id={titleId}>{config.labels.title}</h2>
          <button
            type="button"
            className="Btn Btn--icon ReportDialog-close"
            aria-label={config.labels.cancel}
            onClick={close}
          >
            <i className="fa-solid fa-xmark" aria-hidden="true" />
          </button>
        </header>
        <div className="ReportDialog-body">
          {excerpt ? (
            <section className="ReportDialog-selection" aria-labelledby={excerptId}>
              <h3 id={excerptId}>{config.labels.excerpt}</h3>
              <blockquote>{excerpt}</blockquote>
            </section>
          ) : config.captureSelection ? (
            <p className="ReportDialog-hint">{config.labels.selectionHint}</p>
          ) : null}
          <div className="ReportDialog-field">
            <label id={reasonId}>{config.labels.reason}</label>
            <Select value={reason} onValueChange={setReason} disabled={!reasons.length}>
              <SelectTrigger className="w-full" aria-labelledby={reasonId}>
                <SelectValue placeholder={config.labels.reason} />
              </SelectTrigger>
              <SelectContent portalContainer={dialog}>
                {reasons.map((value) => (
                  <SelectItem key={value} value={value}>
                    {value}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <label className="ReportDialog-field">
            <span>{config.labels.description}</span>
            <textarea className="textarea" name="description" rows={4} maxLength={2000} />
          </label>
        </div>
        <footer className="ReportDialog-actions">
          <button type="button" className="Btn Btn--outline" onClick={close}>
            {config.labels.cancel}
          </button>
          <button className="Btn Btn--primary" type="submit" disabled={submitting || !reason}>
            {submitting ? config.labels.submitting : config.labels.submit}
          </button>
        </footer>
      </form>
    </dialog>
  ) : null
}
