import { createBrowserRouter, Outlet } from "react-router-dom";
import LoginPage from "./modules/auth/views/LoginPage";
import HomePage from "./modules/book/views/HomePage";
import SiteHeader from "./components/site-header";
import MyBooks from "./modules/book-manager/views/MyBooks";
import NewBook from "./modules/book-manager/views/NewBook";
import BookPage from "./modules/book/views/BookPage";
import { CreateChapterPage } from "./modules/book-manager/views/EditChapter";
import ChapterPage from "./modules/book/views/ChapterPage";
import { BookManager } from "./modules/book-manager/views/BookManager";

const router = createBrowserRouter([
  {
    path: "login",
    element: <LoginPage />,
  },
  {
    path: "*",
    element: (
      <>
        <SiteHeader />
        <Outlet />
      </>
    ),
    children: [
      {
        path: "home",
        element: <HomePage />,
      },

      {
        path: "book/:id",
        element: <BookPage />,
      },
      {
        path: "book/:bookId/chapters/:chapterId",
        element: <ChapterPage />,
      },
      {
        path: "my-books",
        element: <MyBooks />,
      },
      {
        path: "manager/book/:bookId",
        element: <BookManager />,
      },
    ],
  },
]);

export default router;
