import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useQueryParamDefault } from '@/lib/router-utils'
import BookInfo from './BookInfo'
import BookChapters from './BookChapters'

export default function BookManagerHomePage() {
  const [tab, setTab] = useQueryParamDefault('tab', 'chapters')

  return (
    <Tabs value={tab ?? ''} onValueChange={setTab}>
      <TabsList>
        <TabsTrigger value="chapters">Chapters</TabsTrigger>
        <TabsTrigger value="info">Book information</TabsTrigger>
      </TabsList>

      <TabsContent value="chapters">
        <BookChapters />
      </TabsContent>

      <TabsContent value="info">
        <BookInfo />
      </TabsContent>
    </Tabs>
  )
}
