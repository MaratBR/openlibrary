import { useMemo, useRef, useState } from 'preact/hooks'
import TextEditor from './TextEditor'

import AgeRatingInput from '@/components/AgeRatingInput'
import Switch from '@/components/Switch'
import { httpUpdateBook, managerBookDetailsSchema } from './api'
import TagsInput from '@/components/TagsInput'
import { PreactIslandProps } from '@/lib/island'

export default function GeneralInformation({ data: dataUnknown }: PreactIslandProps) {
  const initialData = useMemo(() => managerBookDetailsSchema.parse(dataUnknown), [dataUnknown])
  const [data, setData] = useState(initialData)
  const [name, setName] = useState(data.name)
  const [tags, setTags] = useState(data.tags)
  const [rating, setRating] = useState(data.ageRating)
  const [isPubliclyVisible, setPubliclyVisible] = useState(data.isPubliclyVisible)
  const summaryRef = useRef(data.summary)
  const tagsInputRef = useRef<HTMLInputElement | null>(null)
  const summaryInputRef = useRef<HTMLInputElement | null>(null)

  const [loading, setLoading] = useState(false)

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault()
    setLoading(true)
    try {
      const response = await httpUpdateBook(data.id, {
        name,
        tags: tags.map((x) => x.id),
        rating,
        summary: summaryRef.current,
        isPubliclyVisible,
      })
      setData(response.data)
      if (response.notifications) {
        response.notifications.forEach(window.flash)
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div class="card">
      <form onSubmit={handleSubmit}>
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
            <Switch
              name="isPubliclyVisible"
              value={isPubliclyVisible}
              onChange={setPubliclyVisible}
            />
          </div>
        </div>

        <div class="mt-4">
          <button class="btn btn--primary">
            {loading ? <span class="loader loader--dark" /> : window._('bookManager.edit.save')}
          </button>
        </div>
      </form>
    </div>
  )
}
