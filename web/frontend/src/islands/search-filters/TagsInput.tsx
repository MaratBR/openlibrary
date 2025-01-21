import { useEffect, useState } from 'preact/hooks'
import { DefinedTagDto, searchTags } from './api'
import { DropdownCore } from './DropdownCore'
import { _ } from '@/common/i18n'

export type TagsInputProps = {
  tags: DefinedTagDto[]
  // eslint-disable-next-line no-unused-vars
  onInput: (tags: DefinedTagDto[]) => void
}

const defaultTagsPromise = searchTags('')

export default function TagsInput({ tags, onInput }: TagsInputProps) {
  const [searchResults, setSearchResults] = useState<DefinedTagDto[]>([])

  useEffect(() => {
    defaultTagsPromise.then(setSearchResults)
  }, [])

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
                  aria-label={_('search.removeTag')}
                >
                  <span class="material-symbols-outlined !text-[20px]">close</span>
                </button>
              </span>
            ))}
          </div>
        ),
      }}
      slotProps={{
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
