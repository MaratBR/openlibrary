import { DashboardContent } from '@/components/dashboard-layout-components'
import { RenderLazy } from '@/components/RenderLazy'
import Tabs from '@/components/Tabs'
import { createEnumParameter } from '@/lib/parameters'
import { LoaderFunctionArgs, NavLink, useLoaderData } from 'react-router'
import z from 'zod'
import { BookGeneral } from './book-general'
import { BookChapters } from './book-chapters'
import { BMBookAPI } from '@/api/bm/book'

export const bookRouteLoader = async ({ params, request }: LoaderFunctionArgs) => {
  const { bookId } = z.object({ bookId: z.string().nonempty() }).parse(params)
  const resp = await BMBookAPI.getInstance().getBook(bookId)

  return {
    bookResponse: resp,
    bookId,
  }
}

const useTab = createEnumParameter('t', ['general', 'analytics', 'chapters'])

export function Book() {
  const { bookResponse } = useLoaderData<Awaited<ReturnType<typeof bookRouteLoader>>>()

  const [tab, setTab] = useTab()

  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader
        title={
          <div class="flex items-center">
            <NavLink className="btn btn--icon mr-4" to="/books">
              <i class="fa-solid fa-arrow-left" />
            </NavLink>
            {bookResponse.data.name}
          </div>
        }
      />

      <Tabs.Root value={tab || 'general'} onChange={setTab}>
        <Tabs.List>
          <Tabs.Tab value="general">{window._('bookManager.edit.generalInformation')}</Tabs.Tab>
          <Tabs.Tab value="chapters">{window._('bookManager.edit.chapters')}</Tabs.Tab>
          <Tabs.Tab value="analytics">{window._('bookManager.edit.analytics')}</Tabs.Tab>
        </Tabs.List>

        <Tabs.Body>
          <RenderLazy show={tab === 'general'}>
            <BookGeneral book={bookResponse.data} />
          </RenderLazy>

          <RenderLazy show={tab === 'chapters'}>
            <BookChapters book={bookResponse.data} />
          </RenderLazy>
        </Tabs.Body>
      </Tabs.Root>
    </DashboardContent.Root>
  )
}
