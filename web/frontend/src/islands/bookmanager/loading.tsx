import { useLocation } from 'react-router'

function Placeholder({ className = '' }: { className?: string }) {
  return <div className={`skeleton BM-placeholder ${className}`} />
}

function Lines() {
  return (
    <div className="space-y-3">
      <Placeholder />
      <Placeholder className="BM-placeholder--short" />
    </div>
  )
}

export function ManagerSkeleton({ location }: { location?: { pathname: string; search: string } }) {
  const current = useLocation()
  const target = location ?? current
  const library = /^\/books\/?$/.test(target.pathname)
  const details = target.pathname.endsWith('/edit')
  const chapters = new URLSearchParams(target.search).get('t') === 'chapters'
  return (
    <section className="BM-loading" aria-busy="true">
      <p role="status" className="sr-only">
        {window._('bookManager.ui.loading')}
      </p>
      <div aria-hidden="true">
        {!library && <Placeholder className="BM-placeholder--back" />}
        <div className="BM-loadingHeader">
          <Placeholder className="BM-placeholder--heading" />
        </div>
        {library ? (
          <>
            <Placeholder className="BM-placeholder--description" />
            <div className="BM-grid">
              {Array.from({ length: 6 }, (_, i) => (
                <div className="BM-bookCard" key={i}>
                  <div className="BM-bookArt">
                    <Placeholder className="BM-placeholder--cover" />
                  </div>
                  <div className="BM-bookBody space-y-4">
                    <Placeholder className="BM-placeholder--short" />
                    <Placeholder className="BM-placeholder--title" />
                    <Placeholder />
                  </div>
                  <div className="BM-cardFooter">
                    <Placeholder />
                  </div>
                </div>
              ))}
            </div>
          </>
        ) : (
          <>
            <div className="BM-loadingTabs">
              <Placeholder />
              <Placeholder />
              <Placeholder />
            </div>
            {details ? (
              <div className="BM-workspace">
                <div className="Card space-y-6">
                  <Placeholder className="BM-placeholder--title" />
                  {[0, 1].map((i) => (
                    <div className="space-y-2" key={i}>
                      <Placeholder className="BM-placeholder--short" />
                      <Placeholder className="BM-placeholder--field" />
                    </div>
                  ))}
                  <Placeholder className="BM-placeholder--short" />
                  <Placeholder className="BM-placeholder--editor" />
                </div>
                <div className="Card space-y-6">
                  <Placeholder className="BM-placeholder--cover" />
                  <Lines />
                  <Placeholder className="BM-placeholder--field" />
                  <Lines />
                  <Lines />
                </div>
              </div>
            ) : chapters ? (
              <div className="space-y-4">
                <Placeholder className="BM-placeholder--button" />
                {[0, 1, 2, 3].map((i) => (
                  <div className="Card" key={i}>
                    <Lines />
                  </div>
                ))}
              </div>
            ) : (
              <div className="Card space-y-6">
                <div className="flex flex-wrap gap-6 items-center">
                  <Placeholder className="BM-placeholder--cover" />
                  <div className="BM-loadingMetadata">
                    <Lines />
                  </div>
                </div>
                <div className="bg-secondary rounded-xl p-4">
                  <Lines />
                </div>
                <Placeholder className="BM-placeholder--title" />
                <Lines />
                <Lines />
              </div>
            )}
          </>
        )}
      </div>
    </section>
  )
}
