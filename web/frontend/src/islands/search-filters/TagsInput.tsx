import { useState } from 'preact/hooks'
import { DropdownCore } from './DropdownCore'
import { DefinedTagDto, useTagsSearch } from '@/api/search'

export type TagsInputProps = {
  tags: DefinedTagDto[]

  onInput: (tags: DefinedTagDto[]) => void
}

export default function TagsInput({ tags, onInput }: TagsInputProps) {
  const [searchQuery, setSearchQuery] = useState('')

  const query = useTagsSearch({
    query: searchQuery,
    fetchDefault: true,
  })

  const searchResults = query.data ?? []

  function add(tag: DefinedTagDto) {
    if (tags.some((x) => x.id === tag.id)) return
    window.requestAnimationFrame(() => {
      onInput([...tags, tag])
    })
  }

  function remove(tag: DefinedTagDto) {
    if (!tags.some((x) => x.id === tag.id)) return
    window.requestAnimationFrame(() => {
      onInput(tags.filter((x) => x.id !== tag.id))
    })
  }

  return (
    <DropdownCore
      slots={{
        beforeInput: (
          <div class="flex flex-wrap items-center gap-1 m-2 empty:hidden">
            {tags.map((tag) => (
              <span
                key={tag.id}
                class="text-sm whitespace-nowrap inline-flex items-center p-0.5 bg-muted"
              >
                {tag.name}

                <button
                  onClick={(e) => {
                    e.preventDefault()
                    remove(tag)
                  }}
                  class="h-5 hover:text-rose-600"
                  aria-label={window._('search.removeTag')}
                >
                  <i class="fa-solid fa-xmark !text-[20px]" />
                </button>
              </span>
            ))}
          </div>
        ),
      }}
      slotProps={{
        input: {
          onInput: (e) => setSearchQuery((e.target as HTMLInputElement).value),
        },
        menu: {
          className: 'max-h-[300px] overflow-y-auto',
          children: (
            <ul>
              {searchResults.map((tag) =>
                tags.some((x) => x.id === tag.id) ? null : (
                  <li
                    key={tag.id}
                    onClick={() => add(tag)}
                    role="button"
                    class="p-2 cursor-pointer hover:bg-muted hover:text-primary"
                  >
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
