import { useState } from 'react'
import { Fragment } from 'react'
import clsx from 'clsx'
import CSRFInput from '@/components/CSRFInput'
import TagsInput from '@/components/TagsInput'
import { DefinedTagDto } from '@/api/search'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { NavLink } from 'react-router'

export default function NewBookForm() {
  const [stage, _setStage] = useState(0)
  const [activeStage, setActiveStage] = useState(0)
  const [name, setName] = useState('')
  const [rating, setRating] = useState('')
  const [tags, setTags] = useState<DefinedTagDto[]>([])
  const [loading, setLoading] = useState(false)

  const setStage = (stage: number) => {
    setActiveStage(stage)
    _setStage(stage)
  }

  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader
        title={
          <div className="flex items-center">
            <NavLink className="btn btn--icon btn--primary mr-4" to="/books">
              <i className="fa-solid fa-arrow-left" />
            </NavLink>
            {window._('bookManager.newBook.title')}
          </div>
        }
      />

      <DashboardContent.Card className="p-4">
        <form
          className="space-y-4 md:space-y-0 md:px-0 md:grid md:grid-cols-[150px_1fr] md:gap-2"
          action="/books-manager/new"
          method="post"
        >
          <CSRFInput />

          <ul className="flex flex-col pt-8 gap-2">
            {Array.from({ length: 4 }).map((_v, i) => {
              const canNavigateTo = stage >= i && activeStage !== i

              return (
                <li
                  onClick={canNavigateTo ? () => setActiveStage(i) : undefined}
                  className={clsx('text-wrap text-secondary-foreground', {
                    '!text-foreground hover:underline cursor-pointer': canNavigateTo,
                    'font-[600] !text-foreground': activeStage === i,
                  })}
                  key={i}
                >
                  {window._(`bookManager.newBook.stageLabel${i}`)}
                </li>
              )
            })}
          </ul>
          <section>
            {stage > 0 && name && <h2 className="text-xl font-semibold mb-8">{name}</h2>}

            <fieldset className="w-96" style={activeStage === 0 ? {} : { display: 'none' }}>
              <input
                value={name}
                onInput={(e) => setName((e.target as HTMLInputElement).value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault()
                    e.stopPropagation()
                    setStage(1)
                  }
                }}
                autoFocus
                placeholder={window._('bookManager.newBook.namePlaceholder')}
                required
                className="input"
                name="name"
              />

              <button
                disabled={name.trim().length < 2}
                onClick={() => setStage(1)}
                type="button"
                className="btn btn--lg btn--default mt-8 rounded-full"
              >
                {window._('bookManager.newBook.next')}
              </button>
            </fieldset>

            <fieldset className="w-96" style={activeStage === 1 ? {} : { display: 'none' }}>
              <p className="my-4">{window._('bookManager.newBook.selectRating')}</p>
              <fieldset className="flex gap-2">
                {window.__server__.ageRatings.map((ageRating) => {
                  const id = `new-book-${ageRating}`
                  return (
                    <Fragment key={ageRating}>
                      <input
                        key={ageRating}
                        id={id}
                        className="age-rating-input"
                        name="ageRating"
                        value={ageRating}
                        type="radio"
                        checked={ageRating === rating}
                        onChange={() => setRating(ageRating)}
                      />
                      <label data-rating={ageRating} className="age-rating" htmlFor={id}>
                        {ageRating}
                      </label>
                    </Fragment>
                  )
                })}
              </fieldset>

              <div className="mt-4">
                <button
                  disabled={rating === ''}
                  onClick={() => setStage(2)}
                  type="button"
                  className="btn btn--lg btn--default mt-8 rounded-full"
                >
                  {window._('bookManager.newBook.next')}
                </button>
              </div>
            </fieldset>

            <fieldset className="w-[500px]" style={activeStage === 2 ? {} : { display: 'none' }}>
              <p className="mb-4">{window._('bookManager.newBook.selectTags')}</p>
              <TagsInput tags={tags} onInput={setTags} />
              <input hidden name="tags" value={tags.map((x) => x.id).join(',')} />

              <button
                onClick={() => setStage(3)}
                type="button"
                className="btn btn--lg btn--default mt-8 rounded-full"
              >
                {window._('bookManager.newBook.next')}
              </button>
            </fieldset>

            <div style={activeStage === 3 ? {} : { display: 'none' }}>
              <p>{window._('bookManager.newBook.pleaseReview')}</p>

              <dl className="mt-4 dl">
                <dt>{window._('bookManager.newBook.bookName')}:</dt>
                <dd>{name}</dd>
                <dt>{window._('bookManager.newBook.ageRating')}:</dt>
                <dd>{rating}</dd>
                <dt>{window._('bookManager.newBook.tags')}:</dt>
                <dd className="tags items-start flex flex-wrap gap-2">
                  {tags.map((x) => (
                    <a className="tag" key={x.id} href={`/tags/${x.id}`}>
                      {x.name}
                    </a>
                  ))}
                </dd>
              </dl>

              <button
                onClick={() => {
                  setLoading(true)
                }}
                className="btn btn--lg btn--default mt-8 rounded-full"
              >
                {loading ? (
                  <span className="loader loader--dark" />
                ) : (
                  window._('bookManager.newBook.create')
                )}
              </button>
            </div>
          </section>
        </form>
      </DashboardContent.Card>
    </DashboardContent.Root>
  )
}
