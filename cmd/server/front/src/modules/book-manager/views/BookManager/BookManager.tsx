import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useQueryParamDefault } from '@/lib/router-utils'
import BookManagerLayout from './BookManagerLayout'
import BookInfo from './BookInfo'

export default function BookManager() {
  const [tab, setTab] = useQueryParamDefault('tab', 'chapters')

  return (
    <BookManagerLayout>
      <Tabs value={tab ?? ''} onValueChange={setTab}>
        <TabsList>
          <TabsTrigger value="chapters">Chapters</TabsTrigger>
          <TabsTrigger value="info">Book information</TabsTrigger>
        </TabsList>

        <TabsContent value="info">
          <BookInfo />
        </TabsContent>
      </Tabs>
    </BookManagerLayout>
  )
}
