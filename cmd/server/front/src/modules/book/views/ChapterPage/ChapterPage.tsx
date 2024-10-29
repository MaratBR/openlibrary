import {
  BookDetailsDto,
  ChapterDto,
  useBookChapterQuery,
  useBookQuery,
} from "../../api";
import { useParams } from "react-router";
import "./ChapterPage.css";
import BookInfoCard from "../BookPage/BookInfoCard";
import { useChapterName } from "../../utils";
import { Separator } from "@/components/ui/separator";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbList,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { NavLink } from "react-router-dom";

export default function ChapterPage() {
  const { chapterId, bookId } = useParams<{
    chapterId: string;
    bookId: string;
  }>();
  const { data: bookData } = useBookQuery(bookId);
  const { data: chapterData } = useBookChapterQuery(bookId, chapterId);

  return (
    <>
      {chapterData && bookData && (
        <>
          <BookInfo book={bookData} chapter={chapterData.chapter} />
          <ChapterContents chapter={chapterData.chapter} />
        </>
      )}
    </>
  );
}

function BookInfo({
  book,
  chapter,
}: {
  book: BookDetailsDto;
  chapter: ChapterDto;
}) {
  const chapterName = useChapterName({
    name: chapter.name,
    order: chapter.order,
  });

  return (
    <section className="container-default">
      <header className="page-header">
        <Breadcrumb>
          <BreadcrumbList>
            <BreadcrumbItem>
              <NavLink className="link-default" to={`/book/${book.id}`}>
                {book.name}
              </NavLink>
            </BreadcrumbItem>
            <BreadcrumbSeparator />
            <BreadcrumbItem>{chapterName}</BreadcrumbItem>
          </BreadcrumbList>
        </Breadcrumb>
      </header>
      <BookInfoCard book={book} />
    </section>
  );
}

function ChapterContents({ chapter }: { chapter: ChapterDto }) {
  const chapterName = useChapterName({
    name: chapter.name,
    order: chapter.order,
  });

  return (
    <div id="chapter-wrapper" className="chapter-content">
      <header className="py-5">
        <h2 className="font-semibold text-2xl text-center mb-5">
          {chapterName}
        </h2>
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
  );
}
