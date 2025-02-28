import { z } from 'zod'
import { PreactIslandProps } from '../common'
import { useMemo, useRef, useState } from 'preact/hooks'
import { definedTagDtoSchema } from '../search-filters/api'
import TextEditor from './TextEditor'

import TagsInput from '../search-filters/TagsInput'
import CSRFInput from '@/components/CSRFInput'
import AgeRatingInput from '@/components/AgeRatingInput'
import Switch from '@/components/Switch'

export default function GeneralInformation({ data: dataUnknown }: PreactIslandProps) {
  const data = useMemo(() => managerBookDetailsSchema.parse(dataUnknown), [dataUnknown])

  const [name, setName] = useState(data.name)
  const [tags, setTags] = useState(data.tags)
  const [rating, setRating] = useState(data.ageRating)
  const [isPubliclyVisible, setPubliclyVisible] = useState(data.isPubliclyVisible)

  const summaryRef = useRef(data.summary)

  const tagsInputRef = useRef<HTMLInputElement | null>(null)
  const summaryInputRef = useRef<HTMLInputElement | null>(null)

  function handleSubmit(e: SubmitEvent) {
    e.preventDefault()
    tagsInputRef.current!.value = tags.map((tag) => tag.id).join(',')
    summaryInputRef.current!.value = summaryRef.current
    if (e.currentTarget instanceof HTMLFormElement) e.currentTarget.submit()
  }

  return (
    <form method="post" action={`/books-manager/book/${data.id}`} onSubmit={handleSubmit}>
      <CSRFInput />
      <div class="form-control form-control--horizontal">
        <div class="form-control__label">
          <label htmlFor="summary" className="label">
            {window._('bookManager.edit.name')}
          </label>
        </div>
        <div class="form-control__value">
          <input
            value={name}
            onInput={(e) => setName((e.target as HTMLInputElement).value)}
            class="input"
            name="name"
          />
        </div>
      </div>

      <div class="form-control form-control--horizontal">
        <div class="form-control__label">
          <label htmlFor="summary" className="label">
            {window._('bookManager.edit.summary')}
          </label>
        </div>
        <div class="form-control__value">
          <input hidden name="summary" ref={summaryInputRef} />
          <TextEditor value={summaryRef} />
        </div>
      </div>

      <div class="form-control form-control--horizontal">
        <div class="form-control__label">
          <label htmlFor="summary" className="label">
            {window._('bookManager.edit.tags')}
          </label>
        </div>
        <div class="form-control__value">
          <input hidden name="tags" ref={tagsInputRef} />
          <TagsInput tags={tags} onInput={setTags} />
        </div>
      </div>

      <div class="form-control form-control--horizontal">
        <div class="form-control__label">
          <label htmlFor="summary" className="label">
            {window._('bookManager.edit.ageRating')}
          </label>
        </div>
        <div class="form-control__value">
          <AgeRatingInput name="rating" value={rating} onChange={setRating} />
        </div>
      </div>

      <div class="form-control form-control--horizontal">
        <div class="form-control__label">
          <label htmlFor="summary" className="label">
            {window._('bookManager.edit.isPubliclyVisible')}
          </label>
          <p class="form-control__hint">
            {window._('bookManager.edit.isPubliclyVisible_description')}
          </p>
        </div>
        <div class="form-control__value">
          <Switch name="isPubliclyVisible" value={isPubliclyVisible} onChange={setPubliclyVisible} />
        </div>
      </div>

      <div class="mt-4">
        <button class="ol-btn ol-btn--primary rounded-full">
          {window._('bookManager.edit.save')}
        </button>
      </div>
    </form>
  )
}

const bookCollectionDtoSchema = z.object({
  id: z.string(),
  name: z.string(),
  position: z.number(),
  size: z.number(),
})

const bookChapterDtoSchema = z.object({
  id: z.string(),
  order: z.number().min(0).int(),
  name: z.string(),
  words: z.number(),
  createdAt: z.string(),
  summary: z.string(),
})

const managerBookDetailsSchema = z.object({
  id: z.string(),
  name: z.string(),
  ageRating: z.string(),
  adult: z.boolean(),
  tags: z.array(definedTagDtoSchema),
  words: z.number(),
  wordsPerChapter: z.number(),
  collections: z.array(bookCollectionDtoSchema),
  chapters: z.array(bookChapterDtoSchema),
  createdAt: z.string(),
  author: z.object({
    id: z.string(),
    name: z.string(),
  }),
  summary: z.string(),
  isPubliclyVisible: z.boolean(),
  isBanned: z.boolean(),
  cover: z.string(),
})
