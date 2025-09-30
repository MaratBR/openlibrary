import { useMemo } from 'preact/hooks'
import GeneralInformation from './GeneralInformation'
import { managerBookDetailsSchema } from './api'
import BookCover from './BookCover'
import Chapters from './Chapters'
import { useHashQueryValue } from '@/lib/url-hooks'
import { PreactIslandProps } from '@/islands/common/preact-island'
import SideTabs from '@/components/SideTabs'

export default function BookEditor({ data: dataUnknown, rootElement }: PreactIslandProps) {
  const data = useMemo(() => managerBookDetailsSchema.parse(dataUnknown), [dataUnknown])

  const [tabValue, setTab] = useHashQueryValue('tab')
  const tab = tabValue || 'general'

  return (
    <>
      <header class="page-header my-8">
        <h1 class="page-header-text">{data.name}</h1>
        <div class="text-sm">
          <a href="/books-manager?tab=books" class="link">
            {window._('bookManager.edit.backToBooksManager')}
          </a>
          &nbsp;|&nbsp;
          <a href={`/book/${data.id}`} class="link">
            {window._('bookManager.edit.goToPage')}
          </a>
        </div>
      </header>

      <SideTabs.Root onChange={setTab} value={tab}>
        <SideTabs.List>
          <SideTabs.Tab value="general">
            {window._('bookManager.edit.generalInformation')}
          </SideTabs.Tab>
          <SideTabs.Tab value="cover">{window._('bookManager.edit.cover')}</SideTabs.Tab>
          <SideTabs.Tab value="chapters">{window._('bookManager.edit.chapters')}</SideTabs.Tab>
        </SideTabs.List>
        <SideTabs.Body>
          <div class="my-4" style={{ display: tab === 'general' ? 'block' : 'none' }}>
            <GeneralInformation data={data} />
          </div>

          <div class="my-4" style={{ display: tab === 'cover' ? 'block' : 'none' }}>
            <BookCover book={data} />
          </div>

          {tab === 'chapters' && (
            <div class="my-4">
              <Chapters book={data} />
            </div>
          )}
        </SideTabs.Body>
      </SideTabs.Root>
    </>
  )
}
