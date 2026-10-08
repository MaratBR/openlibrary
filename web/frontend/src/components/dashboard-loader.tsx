export function DashboardLoader({ label }: { label: string }) {
  return (
    <div className="DashboardShell-loader" role="progressbar" aria-label={label}>
      <span />
    </div>
  )
}
