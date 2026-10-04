import { useRef, useState } from 'react'
import CSRFInput from '@/components/CSRFInput'
import TagsInput from '@/components/TagsInput'
import AgeRatingInput from '@/components/AgeRatingInput'
import { DefinedTagDto } from '@/features/search'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { BackToBooks, usePageTitle, useUnsavedChanges } from '../ui'

export default function NewBookForm() {
  const [stage, setStage] = useState(0)
  const [furthest, setFurthest] = useState(0)
  const [name, setName] = useState('')
  const [rating, setRating] = useState('')
  const [tags, setTags] = useState<DefinedTagDto[]>([])
  const [loading, setLoading] = useState(false)
  const submitting = useRef(false)
  const validName = name.trim().length >= 2
  const validRating = window.__server__.ageRatings.includes(rating)
  const canVisit = (target: number) =>
    target <= furthest && (target === 0 || validName) && (target <= 1 || validRating)
  const next = () => {
    if ((stage === 0 && !validName) || (stage === 1 && !validRating)) return
    const target = Math.min(3, stage + 1)
    setStage(target)
    setFurthest(Math.max(furthest, target))
  }
  const guard = useUnsavedChanges(!loading && (!!name || !!rating || !!tags.length))
  usePageTitle(window._('bookManager.newBook.title'))
  return (
    <DashboardContent.Root>
      <BackToBooks />
      <DashboardContent.StickyHeader title={window._('bookManager.newBook.title')} />
      <div className="Card max-w-200 mx-auto">
        <p role="status" className="text-secondary-foreground mb-4">
          {window._('bookManager.ui.step', { count: String(stage + 1) })}
        </p>
        <ol className="flex flex-wrap gap-2 mb-6">
          {[0, 1, 2, 3].map((i) => (
            <li key={i}>
              <button
                type="button"
                className={`Btn ${stage === i ? 'Btn--primary' : 'Btn--ghost'}`}
                aria-current={stage === i ? 'step' : undefined}
                disabled={!canVisit(i) || loading}
                onClick={() => setStage(i)}
              >
                {window._(`bookManager.newBook.stageLabel${i}`)}
              </button>
            </li>
          ))}
        </ol>
        <form
          action="/books-manager/new"
          method="post"
          className="space-y-6"
          onSubmit={(event) => {
            if (submitting.current) {
              event.preventDefault()
              return
            }
            if (stage !== 3) {
              event.preventDefault()
              next()
              return
            }
            if (!validName || !validRating) {
              event.preventDefault()
              return
            }
            guard.allowUnload()
            submitting.current = true
            setLoading(true)
          }}
          onKeyDown={(event) => {
            // Step navigation shares button validation; tag-search Enter belongs to the tag control.
            if (
              event.key === 'Enter' &&
              stage < 3 &&
              (event.target as HTMLElement).tagName === 'INPUT' &&
              stage !== 2
            ) {
              event.preventDefault()
              next()
            }
          }}
        >
          <p role="status">{loading ? window._('bookManager.ui.creating') : ''}</p>
          <CSRFInput />
          <input type="hidden" name="name" value={name.trim()} />
          <input type="hidden" name="ageRating" value={rating} />
          <input type="hidden" name="tags" value={tags.map((x) => x.id).join(',')} />
          {stage === 0 && (
            <section className="space-y-3">
              <label className="label" htmlFor="new-book-name">
                {window._('bookManager.newBook.bookName')}
              </label>
              <p id="new-name-help" className="text-secondary-foreground">
                {window._('bookManager.ui.nameHelp')}
              </p>
              <input
                autoFocus
                id="new-book-name"
                className="input"
                value={name}
                aria-describedby="new-name-help"
                aria-invalid={!!name && !validName}
                onChange={(e) => setName(e.currentTarget.value)}
              />
            </section>
          )}
          {stage === 1 && (
            <fieldset className="space-y-3">
              <legend className="label">{window._('bookManager.newBook.ageRating')}</legend>
              <p className="text-secondary-foreground">
                {window._('bookManager.newBook.selectRating')}
              </p>
              <AgeRatingInput value={rating} onChange={setRating} name="rating-choice" />
            </fieldset>
          )}
          {stage === 2 && (
            <section className="space-y-3">
              <label className="label" htmlFor="new-book-tags">
                {window._('bookManager.ui.optionalTags')}
              </label>
              <p className="text-secondary-foreground">
                {window._('bookManager.newBook.selectTags')}
              </p>
              <TagsInput id="new-book-tags" tags={tags} onInput={setTags} />
            </section>
          )}
          {stage === 3 && (
            <section className="space-y-4">
              <h2 className="text-xl">{window._('bookManager.ui.review')}</h2>
              <p>{window._('bookManager.newBook.pleaseReview')}</p>
              <p className="bg-secondary p-4 rounded-xl">
                {window._('bookManager.ui.creationVisibility')}
              </p>
              <dl className="space-y-2">
                <dt className="font-semibold">{window._('bookManager.newBook.bookName')}</dt>
                <dd>{name}</dd>
                <dt className="font-semibold">{window._('bookManager.newBook.ageRating')}</dt>
                <dd>{rating}</dd>
                <dt className="font-semibold">{window._('bookManager.newBook.tags')}</dt>
                <dd className="flex flex-wrap gap-2">
                  {tags.length
                    ? tags.map((tag) => (
                        <span className="Tag" key={tag.id}>
                          {tag.name}
                        </span>
                      ))
                    : window._('bookManager.ui.noTags')}
                </dd>
              </dl>
            </section>
          )}
          <div className="flex justify-between gap-3">
            {stage > 0 ? (
              <button
                type="button"
                className="Btn Btn--outline"
                disabled={loading}
                onClick={() => setStage(stage - 1)}
              >
                {window._('bookManager.ui.previous')}
              </button>
            ) : (
              <span />
            )}
            {stage < 3 ? (
              <button
                key="next"
                type="button"
                className="Btn Btn--primary"
                disabled={(stage === 0 && !validName) || (stage === 1 && !validRating)}
                onClick={next}
              >
                {window._('bookManager.newBook.next')}
              </button>
            ) : (
              <button
                key="create"
                type="submit"
                className="Btn Btn--primary"
                disabled={loading || !validName || !validRating}
              >
                {window._(loading ? 'bookManager.ui.creating' : 'bookManager.ui.create')}
              </button>
            )}
          </div>
        </form>
      </div>
      {guard.prompt}
    </DashboardContent.Root>
  )
}
