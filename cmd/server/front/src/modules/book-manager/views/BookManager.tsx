import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { useQueryParamDefault } from '@/lib/router-utils'
import BookInfo from './BookInfo'
import BookChapters from './BookChapters'
import BookCoverUploader from './BookCoverUploader'
import { useBookManager } from './book-manager-context'

export default function BookManagerHomePage() {
  const [tab, setTab] = useQueryParamDefault('tab', 'chapters')
  const { book } = useBookManager()

  return (
    <Tabs value={tab ?? ''} onValueChange={setTab}>
      <TabsList>
        <TabsTrigger value="chapters">Chapters</TabsTrigger>
        <TabsTrigger value="info">Book information</TabsTrigger>
        <TabsTrigger value="cover">Cover image</TabsTrigger>
      </TabsList>

      <TabsContent value="chapters">
        <BookChapters />
      </TabsContent>

      <TabsContent value="info">
        <BookInfo />
      </TabsContent>

      <TabsContent value="cover">
        <BookCoverUploader bookId={book.id} maxSize={5 * (1 << 20)} currentImage={book.cover} />
      </TabsContent>
    </Tabs>
  )
}
