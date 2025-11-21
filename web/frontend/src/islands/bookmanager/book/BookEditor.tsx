import { useMemo } from 'preact/hooks'
import GeneralInformation from './GeneralInformation'
import { managerBookDetailsSchema } from './api'
import BookCover from './BookCover'
import Chapters from './Chapters'
import { useHashQueryValue } from '@/lib/url-hooks'
import { PreactIslandProps } from '@/islands/common/preact-island'
import Tabs from '@/components/Tabs'

export default function BookEditor({ data: dataUnknown, rootElement: _ }: PreactIslandProps) {
  const data = useMemo(() => managerBookDetailsSchema.parse(dataUnknown), [dataUnknown])

  const [tabValue, setTab] = useHashQueryValue('tab')
  const tab = tabValue || 'general'

  return (
    <>
      <h1 class="page-header inline-block mr-2">{data.name}</h1>
      <div class="page-header__after">
        <a href="/books-manager/books" class="link">
          {window._('bookManager.edit.backToBooksManager')}
        </a>
        &nbsp;|&nbsp;
        <a href={`/book/${data.id}`} class="link">
          {window._('bookManager.edit.goToPage')}
        </a>
      </div>

      <Tabs.Root onChange={setTab} value={tab}>
        <Tabs.List>
          <Tabs.Tab value="general">{window._('bookManager.edit.generalInformation')}</Tabs.Tab>
          <Tabs.Tab value="cover">{window._('bookManager.edit.cover')}</Tabs.Tab>
          <Tabs.Tab value="chapters">{window._('bookManager.edit.chapters')}</Tabs.Tab>
        </Tabs.List>
        <Tabs.Body class="card">
          <div style={{ display: tab === 'general' ? 'block' : 'none' }}>
            <GeneralInformation data={data} />
          </div>

          <div style={{ display: tab === 'cover' ? 'block' : 'none' }}>
            <BookCover book={data} />
          </div>

          {tab === 'chapters' && (
            <div>
              <Chapters book={data} />
            </div>
          )}
        </Tabs.Body>
      </Tabs.Root>
    </>
  )
}
