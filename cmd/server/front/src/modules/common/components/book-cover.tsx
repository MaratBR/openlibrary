import './book-cover.css'

export type BookCoverProps = {
  url: string | undefined | null
  size?: 'sm' | 'md'
}

export default function BookCover({ url, size = 'md' }: BookCoverProps) {
  if (!url) {
    return null
  }

  return (
    <div className={`book-cover relative book-cover--${size}`}>
      <img className="book-cover__img z-[1]" src={url} />
      <img src={url} className="h-full w-full scale-110 absolute blur-lg opacity-70" />
    </div>
  )
}
