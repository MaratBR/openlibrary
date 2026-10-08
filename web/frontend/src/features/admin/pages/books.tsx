import sanitizeHtml from 'sanitize-html'
import { Link, LoaderFunctionArgs, redirect, useLoaderData, useSearchParams } from 'react-router'
import { loadBook, loadBooks } from '@/features/admin/api'
import { Card, Field, PageHeader, Pager, t } from '@/features/admin/components/ui'

export async function booksLoader({ request }: LoaderFunctionArgs) {
  const data = await loadBooks(new URL(request.url).search, request.signal)
  if ('redirect' in data) return redirect(data.redirect)
  return data
}
export function bookLoader({ params, request }: LoaderFunctionArgs) {
  return loadBook(params.id || '', request.signal)
}
function Status({ reasons }: { reasons: readonly string[] }) {
  return (
    <span>
      {reasons.length
        ? reasons.map((reason) => t(`admin.books.reason.${reason}`)).join(' ')
        : t('admin.books.available')}
    </span>
  )
}
export function Books() {
  const data = useLoaderData<typeof booksLoader>()
  const [params, setParams] = useSearchParams()
  return (
    <>
      <PageHeader
        title={t('admin.books.title')}
        actions={
          <a className="Btn Btn--outline" href="/search?admin.link=1">
            {t('admin.books.regularSearch')}
          </a>
        }
      />
      <Card className="admin-card--table">
        <form
          key={params.toString()}
          className="Admin-toolbar"
          onSubmit={(event) => {
            event.preventDefault()
            const q = String(new FormData(event.currentTarget).get('q') || '').trim()
            setParams(q ? { q } : {})
          }}
        >
          <Field id="book-search" label={t('admin.books.search')}>
            <input
              className="input"
              id="book-search"
              name="q"
              type="search"
              defaultValue={params.get('q') || ''}
            />
          </Field>
          <button className="Btn Btn--primary" type="submit">
            {t('common.search')}
          </button>
        </form>
        <div className="Admin-resultCount">
          {t('admin.ui.results', { count: data.total.toLocaleString() })}
        </div>
        {data.books.length === 0 ? (
          <p className="p-6">{t('admin.books.noResults')}</p>
        ) : (
          <div
            className="Admin-tableScroll"
            tabIndex={0}
            role="region"
            aria-label={t('admin.books.title')}
          >
            <table className="table">
              <thead>
                <tr>
                  <th>{t('admin.ui.bookID')}</th>
                  <th>{t('admin.users.name')}</th>
                  <th>{t('admin.books.author')}</th>
                  <th>{t('admin.ui.status')}</th>
                  <th>{t('admin.books.words')}</th>
                  <th>{t('admin.books.chapters')}</th>
                </tr>
              </thead>
              <tbody>
                {data.books.map((book) => (
                  <tr key={book.id}>
                    <td>
                      <Link className="Link" to={`/books/${book.id}`}>
                        {book.id}
                      </Link>
                    </td>
                    <td>
                      <Link className="Link" to={`/books/${book.id}`}>
                        {book.name}
                      </Link>
                    </td>
                    <td>{book.author}</td>
                    <td>
                      <Status reasons={book.visibilityReasons} />
                    </td>
                    <td>{book.words.toLocaleString()}</td>
                    <td>{book.chapters.toLocaleString()}</td>
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
export function BookDetails() {
  const book = useLoaderData<typeof bookLoader>()
  return (
    <>
      <PageHeader
        title={book.name}
        actions={
          <>
            <Link className="Btn Btn--outline" to="/books">
              {t('admin.books.title')}
            </Link>
            <a className="Btn Btn--outline" href={`/moderation#/books/${book.id}`}>
              {t('admin.books.moderation')}
            </a>
            <a className="Btn Btn--primary" href={`/book/${book.id}?admin.override=1`}>
              {t('admin.books.publicPage')}
            </a>
          </>
        }
      />
      <Card>
        <div className="flex flex-wrap gap-6">
          {book.cover.url && (
            <img src={book.cover.url} alt="" className="w-40 object-contain self-start" />
          )}
          <dl className="grid grid-cols-2 gap-4">
            <dt>{t('admin.ui.bookID')}</dt>
            <dd>{book.id}</dd>
            <dt>{t('admin.books.author')}</dt>
            <dd>
              <Link className="Link" to={`/users/${book.author.id}`}>
                {book.author.name}
              </Link>
            </dd>
            <dt>{t('admin.ui.status')}</dt>
            <dd>
              <Status reasons={book.visibilityReasons} />
            </dd>
            <dt>{t('admin.books.ageRating')}</dt>
            <dd>{book.ageRating}</dd>
            <dt>{t('admin.books.words')}</dt>
            <dd>{book.words.toLocaleString()}</dd>
            <dt>{t('admin.books.chapters')}</dt>
            <dd>{book.chapters.toLocaleString()}</dd>
            <dt>{t('admin.books.createdAt')}</dt>
            <dd>{new Date(book.createdAt).toLocaleString()}</dd>
          </dl>
        </div>
      </Card>
      <Card title={t('admin.books.summary')}>
        <div
          className="user-content"
          dangerouslySetInnerHTML={{ __html: sanitizeHtml(book.summary) }}
        />
      </Card>
    </>
  )
}
