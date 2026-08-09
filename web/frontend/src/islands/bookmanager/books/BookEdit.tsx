import { NavLink, useLoaderData } from 'react-router'
import { bookRouteLoader } from './Book'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { useState } from 'react'
import { RichTextInput } from '@/components/rte'
import { useOLEditor } from '@/components/rte/RichTextEditor'
import { useMutation } from '@tanstack/react-query'
import { BMBookAPI } from '@/api/bm/book'
import { FormControl } from '@/components/FormControl'
import Switch from '@/components/Switch'
import TagsInput from '@/components/TagsInput'
import AgeRatingInput from '@/components/AgeRatingInput'

export default function BookEdit() {
  const { bookResponse } = useLoaderData<Awaited<ReturnType<typeof bookRouteLoader>>>()

  const [name, setName] = useState(bookResponse.data.name)
  const [ageRating, setAgeRating] = useState(bookResponse.data.ageRating)
  const [isAdult, setIsAdult] = useState(bookResponse.data.adult)
  const [isPubliclyVisible, setIsPubliclyVisible] = useState(bookResponse.data.isPubliclyVisible)
  const [tags, setTags] = useState(bookResponse.data.tags)

  const summaryEditor = useOLEditor({ content: bookResponse.data.summary })

  const saveMutation = useMutation({
    mutationFn: async () => {
      return await BMBookAPI.getInstance().updateBook(bookResponse.data.id, {
        name,
        summary: summaryEditor.getHTML(),
        isAdult,
        ageRating,
        isPubliclyVisible,
        tags: tags.map((x) => x.id),
      })
    },
    onError(error) {
      window.toast.error(error)
    },
  })

  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader
        title={
          <div className="flex items-center">
            <NavLink
              className="btn btn--icon btn--primary mr-4"
              to={`/books/${bookResponse.data.id}`}
            >
              <i className="fa-solid fa-arrow-left" />
            </NavLink>
            {bookResponse.data.name}
          </div>
        }
      />
      <div className="card mb-4">
        <form className="space-y-2">
          <FormControl htmlFor="name-input" label={window._('bookManager.edit.name')}>
            <input
              id="name-input"
              type="text"
              className="input"
              value={name}
              onChange={(e) => setName((e.target as HTMLInputElement).value)}
            />
          </FormControl>

          <FormControl label={window._('bookManager.edit.tags')} htmlFor="tags-input">
            <TagsInput tags={tags} onInput={setTags} id="tags-input" />
          </FormControl>

          <FormControl label={window._('bookManager.edit.ageRating')}>
            <AgeRatingInput value={ageRating} onChange={setAgeRating} />
          </FormControl>

          <FormControl label={window._('common.adult')}>
            <Switch value={isAdult} onChange={setIsAdult} />
          </FormControl>

          <FormControl
            label={window._('bookManager.edit.isPubliclyVisible')}
            description={window._('bookManager.edit.isPubliclyVisible_description')}
          >
            <Switch value={isPubliclyVisible} onChange={setIsPubliclyVisible} />
          </FormControl>

          <FormControl label={window._('bookManager.edit.summary')}>
            <RichTextInput editor={summaryEditor} />
          </FormControl>

          <div className="mt-4 flex">
            <button className="btn btn--primary" onClick={() => saveMutation.mutate()}>
              {window._('common.save')}
            </button>
          </div>
        </form>
      </div>

      <div className="card mb-4">
        <a
          className="btn btn--outline"
          href={`/book/${bookResponse.data.id}`}
          target="_blank"
          rel="noreferrer"
        >
          {window._('bookManager.edit.goToPage')}
          <i className="fa-solid fa-up-right-from-square ml-2" />
        </a>
      </div>
    </DashboardContent.Root>
  )
}
