import { useLoaderData, useNavigate, useRevalidator } from 'react-router'
import { bookRouteLoader } from './Book'
import { DashboardContent } from '@/components/dashboard-layout-components'
import { useEffect, useState } from 'react'
import { RichTextInput } from '@/components/rte'
import { useOLEditor } from '@/components/rte/RichTextEditor'
import { useMutation } from '@tanstack/react-query'
import { BMBookAPI, ManagerBookDetailsDto } from '@/api/bm/book'
import { FormControl } from '@/components/FormControl'
import Switch from '@/components/Switch'
import TagsInput from '@/components/TagsInput'
import AgeRatingInput from '@/components/AgeRatingInput'
import {
  keepFormFocusVisible,
  BackToBooks,
  BookNavigation,
  Counts,
  CoverImage,
  RequestError,
  usePageTitle,
  useUnsavedChanges,
} from '../ui'

export default function BookEdit() {
  const { bookResponse } = useLoaderData<Awaited<ReturnType<typeof bookRouteLoader>>>()
  return <BookEditForm key={bookResponse.data.id} book={bookResponse.data} />
}

function BookEditForm({ book }: { book: ManagerBookDetailsDto }) {
  const [name, setName] = useState(book.name)
  const [ageRating, setAgeRating] = useState(book.ageRating)
  const [isAdult, setIsAdult] = useState(book.adult)
  const [isPubliclyVisible, setIsPubliclyVisible] = useState(book.isPubliclyVisible)
  const [tags, setTags] = useState(book.tags)
  const [, setSummary] = useState(book.summary)
  const summaryEditor = useOLEditor({
    content: book.summary,
    onUpdate: ({ editor }) => setSummary(editor.getHTML()),
    editorProps: {
      attributes: {
        'aria-label': window._('bookManager.edit.summary'),
        role: 'textbox',
        'aria-multiline': 'true',
      },
    },
  })
  const snapshot = () =>
    JSON.stringify({
      name,
      ageRating,
      isAdult,
      isPubliclyVisible,
      tags: tags.map((t) => t.id),
      summary: summaryEditor.getHTML(),
    })
  const [saved, setSaved] = useState(snapshot)
  const dirty = snapshot() !== saved
  const navigate = useNavigate()
  const revalidator = useRevalidator()
  usePageTitle(`${book.name} · ${window._('bookManager.ui.details')}`)
  const [unconfirmed, setUnconfirmed] = useState(false)
  const saveMutation = useMutation({
    mutationFn: async () => {
      setUnconfirmed(false)
      const submitted = snapshot()
      const response = await BMBookAPI.getInstance().updateBook(book.id, {
        name: name.trim(),
        summary: summaryEditor.getHTML(),
        isAdult,
        ageRating,
        isPubliclyVisible,
        tags: tags.map((x) => x.id),
      })
      response.throwIfError()
      if (response.data.adult !== isAdult) {
        setUnconfirmed(true)
        throw new Error('Book response did not confirm the requested adult setting')
      }
      return submitted
    },
    onSuccess: async (submitted) => {
      setSaved(submitted)
      await revalidator.revalidate()
    },
  })
  useEffect(() => {
    summaryEditor.setEditable(!saveMutation.isPending)
  }, [summaryEditor, saveMutation.isPending])
  useEffect(() => {
    const onResize = () => keepFormFocusVisible(document.activeElement as HTMLElement | null)
    window.visualViewport?.addEventListener('resize', onResize)
    return () => window.visualViewport?.removeEventListener('resize', onResize)
  }, [])
  const guard = useUnsavedChanges(dirty, saveMutation.isPending)
  return (
    <DashboardContent.Root>
      <BackToBooks />
      <p className="text-secondary-foreground">{book.name}</p>
      <DashboardContent.StickyHeader title={window._('bookManager.ui.bookDetails')}>
        <a className="Link" href={`/book/${book.id}`} target="_blank" rel="noreferrer">
          {window._('bookManager.ui.publicPage')} ↗
        </a>
      </DashboardContent.StickyHeader>
      <BookNavigation id={book.id} active="details" />
      <form
        onFocusCapture={(e) => keepFormFocusVisible(e.target as HTMLElement)}
        onSubmit={(e) => {
          e.preventDefault()
          if (name.trim() && !saveMutation.isPending) saveMutation.mutate()
        }}
      >
        <fieldset disabled={saveMutation.isPending}>
          <div className="BM-workspace">
            <div className="Card space-y-4">
              <h2 className="text-xl">{window._('bookManager.ui.theBook')}</h2>
              <FormControl
                htmlFor="name-input"
                label={window._('bookManager.edit.name')}
                error={!name.trim() ? window._('bookManager.ui.nameRequired') : undefined}
              >
                <input
                  id="name-input"
                  required
                  className="input"
                  value={name}
                  aria-invalid={!name.trim()}
                  aria-describedby={!name.trim() ? 'name-input-error' : undefined}
                  onChange={(e) => setName(e.currentTarget.value)}
                />
              </FormControl>
              <FormControl label={window._('bookManager.edit.tags')} htmlFor="tags-input">
                <TagsInput tags={tags} onInput={setTags} id="tags-input" />
              </FormControl>
              <FormControl label={window._('bookManager.edit.summary')}>
                <RichTextInput editor={summaryEditor} />
              </FormControl>
            </div>
            <aside className="Card space-y-4">
              <CoverImage cover={book.cover} />
              <p className="text-sm text-secondary-foreground">
                <Counts chapters={book.chapters.length} words={book.words} />
              </p>
              <h2 className="text-xl">{window._('bookManager.ui.audience')}</h2>
              <FormControl label={window._('bookManager.edit.ageRating')}>
                <AgeRatingInput value={ageRating} onChange={setAgeRating} name="rating" />
              </FormControl>
              <FormControl label={window._('common.adult')}>
                <Switch
                  value={isAdult}
                  onChange={setIsAdult}
                  slotProps={{ input: { 'aria-label': window._('common.adult') } }}
                />
              </FormControl>
              <FormControl
                label={window._('bookManager.edit.isPubliclyVisible')}
                description={window._('bookManager.edit.isPubliclyVisible_description')}
              >
                <Switch
                  value={isPubliclyVisible}
                  onChange={setIsPubliclyVisible}
                  slotProps={{
                    input: { 'aria-label': window._('bookManager.edit.isPubliclyVisible') },
                  }}
                />
              </FormControl>
            </aside>
          </div>
        </fieldset>
        {saveMutation.isError &&
          (unconfirmed ? (
            <p role="alert" className="text-destructive my-3">
              {window._('bookManager.ui.unconfirmedAdult')}
            </p>
          ) : (
            <RequestError />
          ))}
        <footer className="BM-actionBar">
          <p role="status">
            {window._(
              saveMutation.isPending
                ? 'bookManager.ui.saving'
                : dirty
                  ? 'bookManager.ui.unsaved'
                  : saveMutation.isSuccess
                    ? 'bookManager.ui.saved'
                    : 'bookManager.ui.upToDate',
            )}
          </p>
          <div className="flex gap-2">
            <button
              type="button"
              className="Btn Btn--outline"
              disabled={saveMutation.isPending}
              onClick={() => navigate(`/books/${book.id}`)}
            >
              {window._('common.cancel')}
            </button>
            <button
              type="submit"
              className="Btn Btn--primary"
              disabled={!name.trim() || saveMutation.isPending}
            >
              {window._('bookManager.ui.saveChanges')}
            </button>
          </div>
        </footer>
      </form>
      {guard.prompt}
    </DashboardContent.Root>
  )
}
