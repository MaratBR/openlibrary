import { FormEvent, useState } from 'react'
import { useRevalidator } from 'react-router'
import { addReportNote } from '../../api/workspace'
import { FilterSelect, t } from './ui'
export function ReportCollaboration({
  reportId,
  priority,
}: {
  reportId: string
  priority: string
}) {
  const revalidator = useRevalidator()
  const [pending, setPending] = useState(false)
  const [error, setError] = useState('')
  const [saved, setSaved] = useState(false)
  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    if (pending) return
    const form = event.currentTarget
    const fields = new FormData(form)
    setPending(true)
    setError('')
    setSaved(false)
    try {
      await addReportNote(reportId, {
        note: String(fields.get('note') ?? '').trim(),
        priority: String(fields.get('priority') ?? ''),
        reason: String(fields.get('reason') ?? '').trim(),
      })
      form.reset()
      setSaved(true)
      await revalidator.revalidate()
    } catch (err) {
      setError(err instanceof Error ? err.message : t('saveError'))
    } finally {
      setPending(false)
    }
  }
  return (
    <section className="Card Card--dashboard">
      <h2 className="text-lg font-semibold mb-4">{t('caseNotes')}</h2>
      <form onSubmit={submit} className="grid gap-3">
        <label className="grid gap-2">
          {t('internalNote')}
          <textarea name="note" className="input min-h-24" required maxLength={1000} />
        </label>
        <p className="text-sm text-secondary-foreground">{t('internalOnly')}</p>
        <button className="Btn Btn--outline justify-self-start" disabled={pending}>
          {t('addNote')}
        </button>
      </form>
      <form onSubmit={submit} className="grid gap-3 mt-6">
        <FilterSelect
          key={priority}
          name="priority"
          label={t('priority')}
          value={priority}
          options={['low', 'medium', 'high', 'critical'].map((v) => [v, t(v)])}
        />
        <label className="grid gap-2">
          {t('reason')}
          <textarea className="input" name="reason" required />
        </label>
        <button className="Btn Btn--outline justify-self-start" disabled={pending}>
          {t('changePriority')}
        </button>
      </form>
      {error && (
        <p role="alert" className="mt-3">
          {error}
        </p>
      )}
      {saved && (
        <p role="status" className="mt-3">
          {t('saved')}
        </p>
      )}
    </section>
  )
}
