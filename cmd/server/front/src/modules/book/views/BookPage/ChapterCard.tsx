import { BookChapterDto } from "../../api";
import { useChapterName } from "../../utils";
import { NavLink } from "react-router-dom";

export default function ChapterCard({
  chapter,
  bookId,
}: {
  chapter: BookChapterDto;
  bookId: string;
}) {
  const chapterName = useChapterName({
    name: chapter.name,
    order: chapter.order,
  });

  return (
    <NavLink
      to={`/book/${bookId}/chapters/${chapter.id}`}
      className="p-2 rounded-lg border bg-card text-card-foreground shadow-sm block w-full hover:bg-muted"
    >
      <span className="font-[500]">{chapterName}</span>
      &nbsp;&nbsp;&bull;&nbsp;&nbsp;
      <span className="text-sm text-muted-foreground">
        {chapter.words} words &nbsp;&nbsp;&bull;&nbsp;&nbsp; published{" "}
        {new Date(chapter.createdAt).toLocaleDateString("en-US")}
      </span>
    </NavLink>
  );
}
