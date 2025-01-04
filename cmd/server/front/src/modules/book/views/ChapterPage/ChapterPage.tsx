import {
  BookDetailsDto,
  ChapterDto,
  preloadBookChapterQuery,
  useBookChapterQuery,
  useBookQuery,
} from '../../api/api'
import { useParams } from 'react-router'
import './ChapterPage.css'
import { Separator } from '@/components/ui/separator'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { NavLink } from 'react-router-dom'
import { useChapterName } from '../../utils'
import { ArrowLeft, ArrowRight, Book, TableOfContents } from 'lucide-react'
import { useQueryClient } from '@tanstack/react-query'
import React from 'react'

import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
} from '@/components/ui/sheet'
import ChapterCard from '../BookPage/ChapterCard'
import { ScrollArea } from '@/components/ui/scroll-area'
import ChapterComments from './ChapterComments'

export default function ChapterPage() {
  const { chapterId, bookId } = useParams<{
    chapterId: string
    bookId: string
  }>()
  const { data: book } = useBookQuery(bookId)
  const { data: chapterData } = useBookChapterQuery(bookId, chapterId)
  const chapterName = useChapterName(
    chapterData?.chapter.name ?? '',
    chapterData?.chapter.order ?? 0,
  )

  return (
    <TOC book={book} activeChapterId={chapterData?.chapter.id}>
      {book && (
        <section className="bg-muted">
          <div className="container-default">
            <header className="page-header">
              <Breadcrumb>
                <BreadcrumbList>
                  <BreadcrumbItem>
                    <NavLink className="link-default" to={`/book/${book.id}`}>
                      {book.name}
                    </NavLink>
                  </BreadcrumbItem>
                  <BreadcrumbSeparator />
                  <BreadcrumbItem>{chapterData && chapterName}</BreadcrumbItem>
                </BreadcrumbList>
              </Breadcrumb>
            </header>
          </div>
        </section>
      )}

      {chapterData && book && bookId && (
        <>
          <div id="chapter" className="pt-[1px]">
            <ChapterControls chapter={chapterData.chapter} bookId={bookId} />
            <ChapterContents chapter={chapterData.chapter} />
            <ChapterControls chapter={chapterData.chapter} bookId={bookId} />
            <ChapterComments chapterId={chapterData.chapter.id} />
          </div>
        </>
      )}
    </TOC>
  )
}

function ChapterContents({ chapter }: { chapter: ChapterDto }) {
  const chapterName = useChapterName(chapter.name, chapter.order)

  return (
    <div className="chapter-content px-6 2xl:px-0">
      <header className="py-5">
        <h2 className="font-semibold text-2xl text-center mb-5">{chapterName}</h2>
        <Separator />
        {chapter.summary && (
          <>
            <div className="py-5">
              <h3 className="font-semibold text-lg">Summary</h3>
              <p>{chapter.summary}</p>
            </div>
            <Separator />
          </>
        )}
      </header>
      <div dangerouslySetInnerHTML={{ __html: chapter.content }}></div>
    </div>
  )
}

function ChapterControls({ chapter, bookId }: { chapter: ChapterDto; bookId: string }) {
  const openTOC = React.useContext(OpenTOCContext)

  const prevChapterName = useChapterName(
    chapter.prevChapter?.name ?? '',
    chapter.prevChapter?.order ?? 0,
  )

  const nextChapterName = useChapterName(
    chapter.nextChapter?.name ?? '',
    chapter.nextChapter?.order ?? 0,
  )

  const queryClient = useQueryClient()

  return (
    <div
      className="my-6 grid grid-rows-3 px-2
      md:grid-rows-1 md:grid-cols-[1fr_auto_1fr] md:gap-16 md:justify-center"
    >
      <div className="flex md:justify-end">
        {chapter.prevChapter ? (
          <NavLink
            to={`/book/${bookId}/chapters/${chapter.prevChapter.id}#chapter`}
            className="link-default link-chapter-control link-chapter-control--prev"
            onMouseEnter={() =>
              preloadBookChapterQuery(queryClient, bookId, chapter.prevChapter!.id)
            }
            onMouseDown={() =>
              preloadBookChapterQuery(queryClient, bookId, chapter.prevChapter!.id)
            }
          >
            <div className="link-chapter-control__label">
              <ArrowLeft className="link-chapter-control__icon" />
              Previous chapter
            </div>
            <span className="link-chapter-control__name">{prevChapterName}</span>
          </NavLink>
        ) : (
          <NavLink
            to={`/book/${bookId}`}
            end
            className="link-default link-chapter-control link-chapter-control--prev"
          >
            <Book />
            Book page
          </NavLink>
        )}
      </div>

      <div>
        <button
          onClick={() => openTOC()}
          className="h-full link-default flex items-center gap-2 hover:bg-muted px-2"
        >
          <TableOfContents /> ToC
        </button>
      </div>

      <div className="flex justify-start">
        {chapter.nextChapter ? (
          <NavLink
            to={`/book/${bookId}/chapters/${chapter.nextChapter.id}#chapter`}
            className="link-default link-chapter-control link-chapter-control--next"
            onMouseEnter={() =>
              preloadBookChapterQuery(queryClient, bookId, chapter.nextChapter!.id)
            }
            onMouseDown={() =>
              preloadBookChapterQuery(queryClient, bookId, chapter.nextChapter!.id)
            }
          >
            <div className="link-chapter-control__label">
              Next chapter
              <ArrowRight className="link-chapter-control__icon" />
            </div>
            <span className="link-chapter-control__name">{nextChapterName}</span>
          </NavLink>
        ) : (
          <NavLink to={`/book/${bookId}`} className="link-default link-chapter-control">
            <Book />
            Book page
          </NavLink>
        )}
      </div>
    </div>
  )
}

const OpenTOCContext = React.createContext<() => void>(() => {})

function TOC({
  children,
  book,
  activeChapterId,
}: React.PropsWithChildren<{
  book: BookDetailsDto | undefined
  activeChapterId?: string
}>) {
  const [tocOpen, setTOCOpen] = React.useState(false)

  const openTOC = React.useCallback(() => {
    setTOCOpen(true)
  }, [])

  return (
    <>
      <Sheet open={tocOpen} onOpenChange={setTOCOpen}>
        {/* <SheetTrigger>Open</SheetTrigger> */}

        <SheetContent className="md:w-[48rem] !max-w-[100vw]">
          <ScrollArea className="h-screen pr-4">
            <SheetHeader>
              <SheetTitle>Chapters list</SheetTitle>
              <SheetDescription>You can chose a chapter to read here.</SheetDescription>
            </SheetHeader>

            <div className="space-y-1 pt-2">
              {book?.chapters.map((chapter) => (
                <ChapterCard
                  key={chapter.id}
                  onClick={() => setTOCOpen(false)}
                  bookId={book.id}
                  chapter={chapter}
                  className={chapter.id === activeChapterId ? 'bg-muted border' : ''}
                />
              ))}
            </div>
          </ScrollArea>
        </SheetContent>
      </Sheet>
      <OpenTOCContext.Provider value={openTOC}>{children}</OpenTOCContext.Provider>
    </>
  )
}
