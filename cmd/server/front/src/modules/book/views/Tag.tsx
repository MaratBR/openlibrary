import { NavLink } from 'react-router-dom'
import { DefinedTagDto } from '../api'
import { cn } from '@/lib/utils'
import { ExclamationTriangleIcon } from '@radix-ui/react-icons'

export type TagProps = {
  tag: DefinedTagDto
  disableInteractive?: boolean
}

export default function Tag({ tag, disableInteractive = false }: TagProps) {
  return (
    <NavLink
      data-tag-type={tag.category}
      data-adult={tag.isAdult}
      data-defined={true}
      className={cn(
        {
          '!border-red-900/50': tag.isAdult,
        },
        'text-sm bg-muted rounded-sm active:outline outline-2 outline-primary outline-offset-[-1px] inline items-center overflow-hidden border-2 border-muted',
      )}
      to={`/tag/${encodeURIComponent(tag.name)}`}
      onClick={disableInteractive ? (e) => e.preventDefault() : undefined}
    >
      {tag.isSpoiler && <ExclamationTriangleIcon className="inline mx-1" />}

      <span className="mx-0.5 inline whitespace-nowrap">{tag.name}</span>
      {tag.isAdult && <span className="bg-red-900/50 px-1">18+</span>}
    </NavLink>
  )
}
