import { HeartFilledIcon, HeartIcon } from '@radix-ui/react-icons'

export type FavoritesCounterProps = {
  count: number
  isLiked: boolean
  onClick?: React.MouseEventHandler
}

export default function FavoritesCounter({ onClick, count, isLiked }: FavoritesCounterProps) {
  return (
    <button
      onClick={onClick}
      className="flex items-center gap-2 text-lg font-[500] rounded-md hover:bg-highlight p-3 justify-center"
    >
      {isLiked ? (
        <HeartFilledIcon className="text-red-600" width="1.8em" height="1.8em" />
      ) : (
        <HeartIcon width="1.8em" height="1.8em" />
      )}
      <span>{count === 0 && isLiked ? 1 : count}</span>
    </button>
  )
}
