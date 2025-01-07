import { Dexie, Table } from 'dexie'
import { SavedBook, SavedChapter } from './schema'

class BooksDB extends Dexie {
  readonly books!: Table<SavedBook, string>
  readonly chapters!: Table<SavedChapter, string>

  constructor() {
    super('books')

    this.version(1).stores({
      books: '_id',
      chapters: '_id, bookId',
    })
    this.on('populate', () => {
      this.books.clear()
      this.chapters.clear()
    })
  }
}

export const booksDB = new BooksDB()
booksDB.open()
