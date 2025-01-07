import { NavLink } from 'react-router-dom'
import { DefinedTagDto } from '../api/api'
import { cn } from '@/lib/utils'
import { ExclamationTriangleIcon } from '@radix-ui/react-icons'
import './Tag.css'

export type TagProps = {
  tag: DefinedTagDto
  disableInteractive?: boolean
}

export default function Tag({ tag, disableInteractive = false }: TagProps) {
  return (
    <NavLink
      data-tag-type={tag.cat}
      data-adult={tag.adult}
      data-defined={true}
      className={cn(
        {
          'tag--adult': tag.adult,
        },
        'tag',
      )}
      to={`/search?it=${tag.id}`}
      onClick={disableInteractive ? (e) => e.preventDefault() : undefined}
    >
      {tag.spoiler && <ExclamationTriangleIcon className="inline mx-1" />}

      <span className="mx-0.5 inline">{tag.name}</span>
      {tag.adult && <span className="bg-red-900/50 px-1">18+</span>}
    </NavLink>
  )
}
