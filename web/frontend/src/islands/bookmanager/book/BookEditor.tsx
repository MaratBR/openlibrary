import { useMemo } from 'preact/hooks'
import { Tab, Tabs } from './Tabs'
import GeneralInformation from './GeneralInformation'
import { managerBookDetailsSchema } from './api'
import BookCover from './BookCover'
import Chapters from './Chapters'
import { useHashQueryValue } from '@/lib/url-hooks'
import { PreactIslandProps } from '@/islands/common/preact-island'

export default function BookEditor({ data: dataUnknown, rootElement }: PreactIslandProps) {
  const data = useMemo(() => managerBookDetailsSchema.parse(dataUnknown), [dataUnknown])

  const [tabValue, setTab] = useHashQueryValue('tab')
  const tab = tabValue || 'general'

  return (
    <>
      <header class="page-header my-8">
        <h1 class="page-header-text">{data.name}</h1>
        <a href={`/book/${data.id}`} class="ol-btn ol-btn--lg ol-btn--secondary rounded-full">
          {window._('bookManager.edit.goToPage')}
        </a>
      </header>

      <Tabs onChange={setTab} type="primary" value={tab}>
        <Tab value="general">{window._('bookManager.edit.generalInformation')}</Tab>
        <Tab value="cover">{window._('bookManager.edit.cover')}</Tab>
        <Tab value="chapters">{window._('bookManager.edit.chapters')}</Tab>
      </Tabs>

      <div class="my-4" style={{ display: tab === 'general' ? 'block' : 'none' }}>
        <GeneralInformation data={data} />
      </div>

      <div class="my-4" style={{ display: tab === 'cover' ? 'block' : 'none' }}>
        <BookCover book={data} />
      </div>

      {tab === 'chapters' && (
        <div class="my-4">
          <Chapters data={data} rootElement={rootElement} />
        </div>
      )}
    </>
  )
}
