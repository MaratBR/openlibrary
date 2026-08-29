import { useState } from 'react'
import { DropdownCore } from '../islands/search-filters/DropdownCore'
import { DefinedTagDto, useTagsSearch } from '@/features/search'

export type TagsInputProps = {
  tags?: DefinedTagDto[]
  onInput?: (tags: DefinedTagDto[]) => void
  id?: string
}

export default function TagsInput({ tags = [], onInput, id }: TagsInputProps) {
  const [searchQuery, setSearchQuery] = useState('')

  const query = useTagsSearch({
    query: searchQuery,
    fetchDefault: true,
  })

  const searchResults = query.data ?? []

  function add(tag: DefinedTagDto) {
    if (tags.some((x) => x.id === tag.id)) return
    window.requestAnimationFrame(() => {
      onInput?.([...tags, tag])
    })
  }

  function remove(tag: DefinedTagDto) {
    if (!tags.some((x) => x.id === tag.id)) return
    window.requestAnimationFrame(() => {
      onInput?.(tags.filter((x) => x.id !== tag.id))
    })
  }

  return (
    <DropdownCore
      slots={{
        beforeInput: (
          <div className="dropdown__chips">
            {tags.map((tag) => (
              <span key={tag.id} className="tag">
                {tag.name}

                {tag.adult && <span className="tag__adult">&nbsp;18+</span>}

                <button
                  onClick={(e) => {
                    e.preventDefault()
                    remove(tag)
                  }}
                  className="chip__close"
                  aria-label={window._('search.removeTag')}
                >
                  <i className="fa-solid fa-xmark !text-[20px]" />
                </button>
              </span>
            ))}
          </div>
        ),
      }}
      slotProps={{
        input: {
          onInput: (e) => setSearchQuery((e.target as HTMLInputElement).value),
          id,
        },
        menu: {
          className: 'max-h-[300px] overflow-y-auto',
          children: (
            <ul>
              {searchResults.map((tag) =>
                tags.some((x) => x.id === tag.id) ? null : (
                  <li key={tag.id} onClick={() => add(tag)} role="button" className="list-item">
                    {tag.name}
                  </li>
                ),
              )}
            </ul>
          ),
        },
      }}
    />
  )
}
