import { useBEState } from './state'

export function CenterHeader() {
  const scheduledAt = useBEState((state) => state.draft?.scheduledAt)

  return (
    <div className="flex justify-center items-center h-full">
      {scheduledAt && (
        <span className="text-sm text-secondary-foreground">
          {window._('editor.publishingScheduledFor', {
            date: new Date(scheduledAt).toLocaleString(),
          })}
        </span>
      )}
    </div>
  )
}
