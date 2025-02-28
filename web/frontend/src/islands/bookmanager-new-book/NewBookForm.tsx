import { useState } from 'preact/hooks'
import { Fragment } from 'preact'
import TagsInput from '../search-filters/TagsInput'
import { DefinedTagDto } from '../search-filters/api'
import clsx from 'clsx'
import CSRFInput from '@/components/CSRFInput'

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
    <form class="anim-appear px-2 space-y-4 md:space-y-0 md:px-0 md:mt-20 md:grid md:grid-cols-[300px_1fr] md:gap-2" action="/books-manager/new" method="post">
      <CSRFInput />

      <ul class="flex flex-col pt-8 gap-2">
        {Array.from({ length: 4 }).map((_v, i) => {
          const canNavigateTo = stage >= i && activeStage !== i

          return (
            <li
              onClick={canNavigateTo ? () => setActiveStage(i) : undefined}
              class={clsx('text-wrap text-muted-foreground', {
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
        <h1 class="page-header mb-8">
          {stage === 0 && <span>{window._('bookManager.newBook.title')}</span>}
          {stage > 0 && name && <span>{name}</span>}
        </h1>

        <fieldset class="w-96" style={activeStage === 0 ? {} : { display: 'none' }}>
          <input
            value={name}
            onInput={(e) => setName((e.target as HTMLInputElement).value)}
            autofocus
            placeholder={window._('bookManager.newBook.namePlaceholder')}
            required
            class="input"
            name="name"
          />

          <button
            disabled={name.trim().length < 2}
            onClick={() => setStage(1)}
            type="button"
            class="mt-8 ol-btn ol-btn--lg ol-btn--primary rounded-full"
          >
            {window._('bookManager.newBook.next')}
          </button>
        </fieldset>

        <fieldset class="w-96" style={activeStage === 1 ? {} : { display: 'none' }}>
          <p class="my-4">{window._('bookManager.newBook.selectRating')}</p>
          <fieldset class="flex gap-2">
            {window.__server__.ageRatings.map((ageRating) => {
              const id = `new-book-${ageRating}`
              return (
                <Fragment key={ageRating}>
                  <input
                    key={ageRating}
                    id={id}
                    class="age-rating-input"
                    name="ageRating"
                    value={ageRating}
                    type="radio"
                    checked={ageRating === rating}
                    onChange={() => setRating(ageRating)}
                  />
                  <label data-rating={ageRating} class="age-rating" for={id}>
                    {ageRating}
                  </label>
                </Fragment>
              )
            })}
          </fieldset>

          <div class="mt-4">
            <button
              disabled={rating === ''}
              onClick={() => setStage(2)}
              type="button"
              class="mt-8 ol-btn ol-btn--lg ol-btn--primary rounded-full"
            >
              {window._('bookManager.newBook.next')}
            </button>
          </div>
        </fieldset>

        <fieldset class="w-[500px]" style={activeStage === 2 ? {} : { display: 'none' }}>
          <p class="mb-4">{window._('bookManager.newBook.selectTags')}</p>
          <TagsInput tags={tags} onInput={setTags} />
          <input hidden name="tags" value={tags.map((x) => x.id).join(',')} />

          <button
            onClick={() => setStage(3)}
            type="button"
            class="mt-8 ol-btn ol-btn--lg ol-btn--primary rounded-full"
          >
            {window._('bookManager.newBook.next')}
          </button>
        </fieldset>

        <div style={activeStage === 3 ? {} : { display: 'none' }}>
          <p>{window._('bookManager.newBook.pleaseReview')}</p>

          <dl class="mt-4 dl">
            <dt>{window._('bookManager.newBook.bookName')}:</dt>
            <dd>{name}</dd>
            <dt>{window._('bookManager.newBook.ageRating')}:</dt>
            <dd>{rating}</dd>
            <dt>{window._('bookManager.newBook.tags')}:</dt>
            <dd class="tabs items-start flex flex-wrap gap-2">
              {tags.map((x) => (
                <a class="tab" key={x.id} href={`/tags/${x.id}`}>
                  {x.name}
                </a>
              ))}
            </dd>
          </dl>

          <button
            onClick={() => {
              setLoading(true)
            }}
            class="mt-8 ol-btn ol-btn--lg ol-btn--primary rounded-full"
          >
            {loading ? <span class="loader loader--dark" /> : window._('bookManager.newBook.create')}
          </button>
        </div>
      </section>
    </form>
  )
}
