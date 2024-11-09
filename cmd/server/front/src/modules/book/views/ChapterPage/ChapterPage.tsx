import { ChapterDto, useBookChapterQuery, useBookQuery } from '../../api'
import { useParams } from 'react-router'
import './ChapterPage.css'
import BookInfoCard from '../BookPage/BookInfoCard'
import { Separator } from '@/components/ui/separator'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'
import { NavLink } from 'react-router-dom'
import { useChapterName } from '../../utils'
import { ArrowLeft, ArrowRight, Book } from 'lucide-react'
import React from 'react'

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
    <>
      {book && (
        <section className="bg-muted pb-8">
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
            <BookInfoCard book={book} />
          </div>
        </section>
      )}

      {chapterData && book && bookId && (
        <>
          <div id="chapter" className="pt-[1px]">
            <ChapterControls chapter={chapterData.chapter} bookId={bookId} />
            <ChapterContents chapter={chapterData.chapter} />
            <ChapterControls chapter={chapterData.chapter} bookId={bookId} />
          </div>
        </>
      )}
    </>
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
  const prevChapterName = useChapterName(
    chapter.prevChapter?.name ?? '',
    chapter.prevChapter?.order ?? 0,
  )

  const nextChapterName = useChapterName(
    chapter.nextChapter?.name ?? '',
    chapter.nextChapter?.order ?? 0,
  )

  return (
    <div className="my-6 flex justify-center gap-16 container-default h-12">
      {chapter.prevChapter ? (
        <NavLink
          to={`/book/${bookId}/chapters/${chapter.prevChapter.id}#chapter`}
          className="link-default link-chapter-control link-chapter-control--prev"
        >
          <div className="link-chapter-control__label">
            <ArrowLeft className="link-chapter-control__icon" />
            Previous chapter
          </div>
          <span className="link-chapter-control__name">{prevChapterName}</span>
        </NavLink>
      ) : (
        <NavLink to={`/book/${bookId}`} className="link-default link-chapter-control">
          <Book />
          Book page
        </NavLink>
      )}

      {chapter.nextChapter && (
        <NavLink
          to={`/book/${bookId}/chapters/${chapter.nextChapter.id}#chapter`}
          className="link-default link-chapter-control link-chapter-control--next"
        >
          <div className="link-chapter-control__label">
            Next chapter
            <ArrowRight className="link-chapter-control__icon" />
          </div>
          <span className="link-chapter-control__name">{nextChapterName}</span>
        </NavLink>
      )}
    </div>
  )
}
