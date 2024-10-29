import React from "react";
import { ManagerBookDetailsDto } from "../../api";

export type BookManagerContext = {
  book: ManagerBookDetailsDto;
};

export const BookManagerContext =
  React.createContext<BookManagerContext | null>(null);

export function useBookManager() {
  const ctx = React.useContext(BookManagerContext);
  if (ctx === null)
    throw new Error("useBookManager must be used within a BookManager");
  return ctx;
}
