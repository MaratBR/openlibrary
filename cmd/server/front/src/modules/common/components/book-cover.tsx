import './book-cover.css'
import bookCoverUrl from './book_cover_1.svg'

export type BookCoverProps = {
  url: string | undefined | null
  name: string
  size?: 'sm' | 'md'
}

export default function BookCover({ url, size = 'md', name }: BookCoverProps) {
  return (
    <div className={`book-cover relative book-cover--${size}`}>
      {url ? <img className="book-cover__img z-[1]" src={url} /> : <GeneratedCover name={name} />}
    </div>
  )
}

import React from 'react'

function GeneratedCover({ name }: { name: string }) {
  return (
    <div className="relative">
      <span className="text-black font-text text-lg left-8 right-8 top-4 bottom-24 flex items-center justify-center absolute break-words text-center">
        {name}
      </span>
      <img width={240} height={360} src={bookCoverUrl} />
    </div>
  )
}
