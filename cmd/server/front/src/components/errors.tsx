export function UrlError({ children }: React.PropsWithChildren) {
  return (
    <div className="error-container error-container--url">
      <p className="text-red-500">Invalid URL</p>
      <p className="text-red-500">{children}</p>
    </div>
  )
}
