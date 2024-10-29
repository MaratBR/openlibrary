import {
  AgeRating,
  AuthorBookDto,
  BookChapterDto,
  BookCollectionDto,
  BookDetailsDto,
  TagDto,
} from "@/modules/book/api";
import { httpClient } from "@/modules/common/api";

export type CreateBookRequest = {
  name: string;
  ageRating: AgeRating;
  tags: string[];
};

export type CreateBookResponse = {
  id: string;
};

export function httpCreateBook(
  req: CreateBookRequest
): Promise<CreateBookResponse> {
  return httpClient
    .post("/api/manager/books", { json: req })
    .then((r) => r.json());
}

export type CreateBookChapterRequest = {
  content: string;
  isAdultOverride: boolean;
  name: string;
  summary: string;
};

export type CreateBookChapterResponse = {
  id: string;
};

export function httpCreateBookChapter(
  bookId: string,
  req: CreateBookChapterRequest
): Promise<CreateBookChapterResponse> {
  return httpClient
    .post(`/api/manager/books/${bookId}/chapters`, { json: req })
    .then((r) => r.json());
}

export type ManagerBookDetailsDto = {
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

export type ManagerGetBookResponse = ManagerBookDetailsDto;

export function httpManagerGetBook(
  id: string
): Promise<ManagerGetBookResponse> {
  return httpClient.get(`/api/manager/books/${id}`).then((r) => r.json());
}

export function httpGetMyBooks(): Promise<{ books: AuthorBookDto[] }> {
  return httpClient.get("/api/manager/books/my-books").then((r) => r.json());
}
