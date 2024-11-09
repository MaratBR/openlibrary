import React from 'react'
import { ManagerBookChapterDto } from '../api'
import BookChapterCard from './BookChapterCard'

export type BookChaptersListProps = {
  chapters: ManagerBookChapterDto[]
}

export default function BookChaptersList({ chapters }: BookChaptersListProps) {
  return (
    <div className="space-y-2">
      {chapters.map((chapter) => (
        <BookChapterCard key={chapter.id} chapter={chapter} />
      ))}
    </div>
  )
}
