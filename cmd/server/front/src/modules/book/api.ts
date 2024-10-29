import { useQuery } from "@tanstack/react-query";
import { httpClient } from "../common/api";

export type AuthorBookDto = {
  id: string;
  name: string;
  createdAt: string;
  ageRating: AgeRating;
  words: number;
  wordsPerChapter: number;
  chapters: number;
  tags: TagDto[];
  collections: BookCollectionDto[];
};

export type BookCollectionDto = {
  id: string;
  name: string;
  position: number;
  size: number;
};

export type TagsCategory =
  | "other"
  | "warning"
  | "fandom"
  | "relationship"
  | "relationshipType"
  | "unknown";

export type TagDto = {
  id: string;
  name: string;
  isAdult: boolean;
  isSpoiler: boolean;
  category: TagsCategory;
  isDefined: boolean;
};

export type AgeRating = "?" | "G" | "PG" | "PG-13" | "R" | "NC-17";

export const AGE_RATINGS_LIST: AgeRating[] = [
  "?",
  "G",
  "PG",
  "PG-13",
  "R",
  "NC-17",
];

export type BookChapterDto = {
  id: string;
  order: number;
  name: string;
  words: number;
  createdAt: string;
};

export type BookDetailsDto = {
  id: string;
  name: string;
  ageRating: AgeRating;
  isAdult: boolean;
  tags: TagDto[];
  words: number;
  wordsPerChapter: number;
  collections: BookCollectionDto[];
  chapters: BookChapterDto[];
  createdAt: string;
  author: {
    id: string;
    name: string;
  };
  permissions: {
    canEdit: boolean;
  };
};

export type GetBookResponse = BookDetailsDto;

export function httpGetBook(id: string): Promise<GetBookResponse> {
  return httpClient.get(`/api/books/${id}`).then((r) => r.json());
}

export function useBookQuery(bookId: string | undefined) {
  return useQuery({
    queryKey: ["book", bookId],
    enabled: !!bookId,
    queryFn: () => httpGetBook(bookId!),
    staleTime: 0,
    gcTime: 60000,
  });
}

export type ChapterDto = {
  id: string;
  name: string;
  words: number;
  content: string;
  isAdultOverride: boolean;
  createdAt: string;
  order: number;
  summary: string;
};

export type GetBookChapterResponse = {
  chapter: ChapterDto;
};

export function httpGetBookChapter(
  bookId: string,
  chapterId: string
): Promise<GetBookChapterResponse> {
  return httpClient
    .get(`/api/books/${bookId}/chapters/${chapterId}`)
    .then((r) => r.json());
}

export function useBookChapterQuery(
  bookId: string | undefined,
  chapterId: string | undefined
) {
  return useQuery({
    queryKey: ["book", bookId, "chapter", chapterId],
    enabled: !!bookId && !!chapterId,
    queryFn: () => httpGetBookChapter(bookId!, chapterId!),
    staleTime: 0,
    gcTime: 60000,
  });
}
