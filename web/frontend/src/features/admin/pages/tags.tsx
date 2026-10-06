import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Link, LoaderFunctionArgs, useLoaderData, useLocation, useSearchParams } from 'react-router'
import { AdminTag, categories, loadTag, loadTags, save } from '@/features/admin/api'
import {
  Card,
  Dialog,
  Empty,
  Field,
  Flag,
  Notice,
  PageHeader,
  Pager,
  SaveBar,
  t,
  useDirtyForm,
  useSave,
} from '@/features/admin/components/ui'

export function tagsLoader({ request }: LoaderFunctionArgs) {
  return loadTags(new URL(request.url).search, request.signal)
}
export function tagLoader({ params, request }: LoaderFunctionArgs) {
  return loadTag(params.id || '', request.signal)
}
export function Tags() {
  const data = useLoaderData<typeof tagsLoader>()
  const [params, setParams] = useSearchParams()
  const location = useLocation()
  const state = { back: `/tags${location.search}` }
  return (
    <>
      <PageHeader title={t('admin.sidebar.tags')} description={t('admin.ui.tagsDescription')} />
      <Card className="admin-card--table">
        <form
          className="Admin-toolbar"
          key={params.toString()}
          onSubmit={(event) => {
            event.preventDefault()
            const fields = new FormData(event.currentTarget)
            const q = new URLSearchParams()
            const search = String(fields.get('q') || '').trim()
            if (search) q.set('q', search)
            for (const key of ['onlyParentTags', 'onlyAdultTags'])
              if (fields.has(key)) q.set(key, 'true')
            setParams(q)
          }}
        >
          <Field id="tag-search" label={t('admin.ui.searchTags')}>
            <input
              className="input"
              id="tag-search"
              name="q"
              type="search"
              defaultValue={params.get('q') || ''}
              placeholder={t('common.search')}
            />
          </Field>
          <div className="Admin-filterChecks">
            <label>
              <input
                className="checkbox"
                type="checkbox"
                name="onlyParentTags"
                defaultChecked={
                  params.get('onlyParentTags') === 'true' || params.get('onlyParentTags') === 'on'
                }
              />
              {t('admin.tags.onlyParentTags')}
            </label>
            <label>
              <input
                className="checkbox"
                type="checkbox"
                name="onlyAdultTags"
                defaultChecked={
                  params.get('onlyAdultTags') === 'true' || params.get('onlyAdultTags') === 'on'
                }
              />
              {t('admin.tags.onlyAdultTags')}
            </label>
          </div>
          <button className="Btn Btn--primary" type="submit">
            {t('common.apply')}
          </button>
        </form>
        <div className="Admin-resultCount">
          {t('admin.ui.results', { count: data.total.toLocaleString() })}
        </div>
        {data.tags.length === 0 ? (
          <Empty kind="tags" />
        ) : (
          <div
            className="Admin-tableScroll"
            tabIndex={0}
            role="region"
            aria-label={t('admin.sidebar.tags')}
          >
            <table className="table Admin-tagTable">
              <thead>
                <tr>
                  <th>{t('admin.tags.name')}</th>
                  <th>{t('admin.tags.synonym')}</th>
                  <th>{t('admin.tags.adult')}</th>
                  <th>{t('admin.tags.spoiler')}</th>
                  <th>{t('admin.tags.createdAt')}</th>
                  <th>{t('admin.users.actions')}</th>
                </tr>
              </thead>
              <tbody>
                {data.tags.map((tag) => (
                  <tr key={tag.id}>
                    <td>
                      <Link className="Link font-semibold" state={state} to={`/tags/${tag.id}`}>
                        {tag.name}
                      </Link>
                      <small className="block text-secondary-foreground">
                        {t(`tag.${tag.category}`)}
                        {tag.isDefault && ` · ${t('admin.tags.default')}`}
                      </small>
                    </td>
                    <td>
                      {tag.synonymOf ? (
                        <Link to={`/tags/${tag.synonymOf.id}`} state={state} className="Link">
                          {tag.synonymOf.name}
                        </Link>
                      ) : (
                        t('common.none')
                      )}
                    </td>
                    <td>
                      <Flag value={tag.adult} />
                    </td>
                    <td>
                      <Flag value={tag.spoiler} />
                    </td>
                    <td>
                      <time dateTime={tag.createdAt}>{tag.createdAt}</time>
                    </td>
                    <td>
                      <Link
                        to={`/tags/${tag.id}/edit`}
                        state={state}
                        className="Btn Btn--outline Btn--sm"
                      >
                        {t('common.edit')}
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>
      <Pager page={data.page} total={data.totalPages} />
    </>
  )
}
function useBackToTags() {
  const location = useLocation()
  return typeof location.state?.back === 'string' && location.state.back.startsWith('/tags')
    ? location.state.back
    : '/tags'
}
export function TagDetails() {
  const tag = useLoaderData<typeof tagLoader>()
  const back = useBackToTags()
  return (
    <>
      <Link to={back} className="Admin-back">
        ← {t('admin.ui.backTags')}
      </Link>
      <PageHeader
        title={tag.name}
        description={t('admin.ui.tagDetails')}
        actions={
          <Link className="Btn Btn--primary" state={{ back }} to={`/tags/${tag.id}/edit`}>
            {t('common.edit')}
          </Link>
        }
      />
      <div className="Admin-workspace">
        <Card title={t('admin.ui.tagDescription')}>
          <p className={tag.description ? 'whitespace-pre-wrap' : 'text-secondary-foreground'}>
            {tag.description || t('admin.tags.noDescription')}
          </p>
          <div className="Admin-actions mt-6">
            <a className="Btn Btn--outline" href={`/tag/${tag.id}`}>
              {t('common.view')}
            </a>
            <a className="Btn Btn--outline" href={`/search?it=${tag.id}`}>
              {t('admin.tags.books')}
            </a>
          </div>
        </Card>
        <Card title={t('admin.ui.tagDetails')}>
          <dl className="Admin-properties">
            <dt>{t('admin.tags.type')}</dt>
            <dd>{t(`tag.${tag.category}`)}</dd>
            <dt>{t('admin.tags.adult')}</dt>
            <dd>
              <Flag value={tag.adult} />
            </dd>
            <dt>{t('admin.tags.spoiler')}</dt>
            <dd>
              <Flag value={tag.spoiler} />
            </dd>
            <dt>{t('admin.ui.parentTag')}</dt>
            <dd>
              {tag.synonymOf ? (
                <Link className="Link" to={`/tags/${tag.synonymOf.id}`}>
                  {tag.synonymOf.name}
                </Link>
              ) : (
                t('admin.ui.noParent')
              )}
            </dd>
            <dt>{t('admin.tags.createdAt')}</dt>
            <dd>
              <time dateTime={tag.createdAt}>{tag.createdAt}</time>
            </dd>
          </dl>
        </Card>
      </div>
    </>
  )
}
export function TagEdit() {
  const tag = useLoaderData<typeof tagLoader>()
  const back = useBackToTags()
  return <TagForm key={tag.id} tag={tag} back={back} />
}
function TagForm({ tag, back }: { tag: AdminTag; back: string }) {
  const initial = {
    name: tag.name,
    description: tag.description,
    type: tag.category,
    adult: tag.adult,
    spoiler: tag.spoiler,
    parent: tag.synonymOf,
  }
  const [values, setValues] = useState(initial)
  const [baseline, setBaseline] = useState(initial)
  const [picker, setPicker] = useState(false)
  const state = useSave()
  const prompt = useDirtyForm(JSON.stringify(values) !== JSON.stringify(baseline), state.pending)
  const change = (next: Partial<typeof values>) => {
    setValues((v) => ({ ...v, ...next }))
    state.changed()
  }
  return (
    <>
      <Link className="Admin-back" to={`/tags/${tag.id}`} state={{ back }}>
        ← {tag.name}
      </Link>
      <PageHeader title={t('admin.ui.editTag')} description={tag.name} />
      <form
        onSubmit={(event) =>
          state.submit(event, async () => {
            await save(`/tags/${tag.id}`, {
              name: values.name,
              description: values.description,
              type: values.type,
              adult: values.adult ? 'on' : '',
              spoiler: values.spoiler ? 'on' : '',
              synonymOf: values.parent?.id || '',
            })
            setBaseline(values)
          })
        }
      >
        <fieldset disabled={state.pending} className="Admin-workspace">
          <Card title={t('admin.ui.tagDetails')}>
            <Field id="tag-name" label={t('admin.tags.name')}>
              <input
                id="tag-name"
                className="input"
                required
                maxLength={50}
                value={values.name}
                onChange={(event) => change({ name: event.target.value })}
              />
            </Field>
            <Field id="tag-description" label={t('admin.ui.tagDescription')}>
              <textarea
                id="tag-description"
                className="textarea"
                maxLength={500}
                rows={8}
                value={values.description}
                onChange={(event) => change({ description: event.target.value })}
              />
            </Field>
            <Field id="tag-category" label={t('admin.tags.type')}>
              <select
                className="Select"
                id="tag-category"
                value={values.type}
                onChange={(event) => change({ type: event.target.value as AdminTag['category'] })}
              >
                {categories.map((category) => (
                  <option value={category} key={category}>
                    {t(`tag.${category}`)}
                  </option>
                ))}
              </select>
            </Field>
          </Card>
          <div className="Admin-stack">
            <Card title={t('admin.ui.parentTag')} description={t('admin.ui.parentHelp')}>
              <p className="mb-4">{values.parent?.name || t('admin.ui.noParent')}</p>
              <div className="Admin-actions">
                <button className="Btn Btn--outline" type="button" onClick={() => setPicker(true)}>
                  {t('common.change')}
                </button>
                {values.parent && (
                  <button
                    type="button"
                    className="Btn Btn--ghost"
                    onClick={() => change({ parent: null })}
                  >
                    {t('admin.ui.removeParent')}
                  </button>
                )}
              </div>
            </Card>
            <Card title={t('admin.ui.status')}>
              <label className="Admin-checkField">
                <input
                  type="checkbox"
                  className="checkbox"
                  checked={values.adult}
                  onChange={(event) => change({ adult: event.target.checked })}
                />
                <span>
                  <strong>{t('admin.tags.adult')}</strong>
                  <small>{t('admin.tags.adultDescription')}</small>
                </span>
              </label>
              <label className="Admin-checkField">
                <input
                  type="checkbox"
                  className="checkbox"
                  checked={values.spoiler}
                  onChange={(event) => change({ spoiler: event.target.checked })}
                />
                <span>
                  <strong>{t('admin.tags.spoiler')}</strong>
                  <small>{t('admin.tags.spoilerDescription')}</small>
                </span>
              </label>
            </Card>
          </div>
        </fieldset>
        <SaveBar state={state} cancelTo={`/tags/${tag.id}`} cancelState={{ back }} />
      </form>
      {picker && (
        <ParentPicker
          id={tag.id}
          onClose={() => setPicker(false)}
          onSelect={(parent) => {
            change({ parent })
            setPicker(false)
          }}
        />
      )}
      {prompt}
    </>
  )
}
function ParentPicker({
  id,
  onClose,
  onSelect,
}: {
  id: string
  onClose: () => void
  onSelect: (tag: NonNullable<AdminTag['synonymOf']>) => void
}) {
  const [query, setQuery] = useState('')
  const [search, setSearch] = useState('')
  const result = useQuery({
    queryKey: ['admin-parent-tags', search],
    queryFn: ({ signal }) =>
      loadTags(`?${new URLSearchParams({ q: search, onlyParentTags: 'true' })}`, signal),
    retry: false,
  })
  return (
    <Dialog title={t('admin.ui.chooseParent')} onClose={onClose}>
      <form
        className="Admin-pickerSearch"
        onSubmit={(event) => {
          event.preventDefault()
          setSearch(query.trim())
        }}
      >
        <Field id="parent-search" label={t('admin.ui.searchTags')}>
          <input
            className="input"
            id="parent-search"
            type="search"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
          />
        </Field>
        <button className="Btn Btn--primary" type="submit">
          {t('common.search')}
        </button>
      </form>
      {result.isPending ? (
        <Notice>{t('admin.ui.loading')}</Notice>
      ) : result.isError ? (
        <>
          <Notice error>{t('admin.ui.loadError')}</Notice>
          <button className="Btn Btn--outline" onClick={() => result.refetch()}>
            {t('admin.ui.retry')}
          </button>
        </>
      ) : (
        <ul className="Admin-pickerList">
          {result.data.tags
            .filter((tag) => tag.id !== id)
            .map((tag) => (
              <li key={tag.id}>
                <button type="button" onClick={() => onSelect({ id: tag.id, name: tag.name })}>
                  {tag.name}
                  <span aria-hidden="true">→</span>
                </button>
              </li>
            ))}
          {result.data.tags.filter((tag) => tag.id !== id).length === 0 && (
            <li>{t('admin.ui.emptyTags')}</li>
          )}
        </ul>
      )}
    </Dialog>
  )
}
