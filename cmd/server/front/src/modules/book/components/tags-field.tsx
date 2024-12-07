import { DefinedTagDto, httpTagsSearch } from '../api'
import Tag from '../views/Tag'
import React, { useMemo } from 'react'
import { AutoComplete } from '@/components/autocomplete'
import { useQuery } from '@tanstack/react-query'
import { X } from 'lucide-react'

export type TagsFieldProps = {
  value: DefinedTagDto[]
  onChange: (value: DefinedTagDto[]) => void
  disabled?: boolean
}

export default function TagsField({ value, onChange, disabled = false }: TagsFieldProps) {
  function handleTagSelected(tag: DefinedTagDto) {
    if (value.some((x) => x.id === tag.id)) return
    onChange([...value, tag])
  }

  const removeTag = (tag: DefinedTagDto) => {
    onChange(value.filter((x) => x.id !== tag.id))
  }

  return (
    <div className="">
      <div className="flex flex-wrap gap-1 mb-3 pt-2 pb-1">
        {value.map((tag) => (
          <div key={tag.id} className="flex items-center">
            <Tag tag={tag} disableInteractive />
            <button
              disabled={disabled}
              className="rounded-full size-[18px] p-[1px] transition-colors ml-1 hover:bg-muted"
              onClick={(e) => {
                e.preventDefault()
                removeTag(tag)
              }}
            >
              <X size={16} />
            </button>
          </div>
        ))}
        {value.length === 0 && <p className="text-muted-foreground pl-3">No tags selected</p>}
      </div>

      <TagsAutocomplete disabled={disabled} selectedTags={value} onTagAdded={handleTagSelected} />
    </div>
  )
}

export type TagsAutocompleteProps = {
  onTagAdded: (value: DefinedTagDto) => void
  selectedTags: DefinedTagDto[]
  disabled?: boolean
  fetchWhenCollapsed?: boolean
}

export function TagsAutocomplete({
  onTagAdded,
  selectedTags,
  disabled,
  fetchWhenCollapsed = false,
}: TagsAutocompleteProps) {
  const [query, setQuery] = React.useState('')
  const debouncedSetQuery = React.useMemo(() => debounce(setQuery, 500), [])
  const [isOpen, setOpen] = React.useState(false)

  const { data } = useQuery({
    queryKey: ['tags', query],
    queryFn: () => httpTagsSearch(query),
    enabled: fetchWhenCollapsed || isOpen,
  })

  const options = useMemo(() => {
    const tags = data?.tags ?? []
    if (tags.length > 0 && selectedTags.length > 0) {
      return tags.filter((t) => !selectedTags.some((st) => st.id === t.id))
    } else {
      return tags
    }
  }, [data, selectedTags])

  function handleOpen() {
    if (!fetchWhenCollapsed) {
      setOpen(true)
    }
  }

  function handleClose() {
    if (!fetchWhenCollapsed) {
      setOpen(false)
    }
  }

  return (
    <AutoComplete<DefinedTagDto>
      getKey={(t) => t.id}
      getLabel={(t) => t.name}
      placeholder="Start typing tag name"
      emptyMessage="No tags found"
      disabled={disabled}
      onInputValueChange={debouncedSetQuery}
      options={options}
      onValueChange={onTagAdded}
      itemComponent={DefinedTagListItem}
      onOpen={handleOpen}
      onClosed={handleClose}
    />
  )
}

function DefinedTagListItem({ value }: { value: DefinedTagDto }) {
  return (
    <div>
      <div>
        <span>{value.name}</span>

        <span className="ml-2 inline-block p-1 rounded-sm bg-muted">{value.cat}</span>
      </div>

      <div className="text-muted-foreground text-sm">{value.desc}</div>
    </div>
  )
}

function debounce<Args extends unknown[]>(fn: (...args: Args) => void, delay: number) {
  let timer: number
  return (...args: Args) => {
    window.clearTimeout(timer)
    timer = window.setTimeout(() => fn(...args), delay)
  }
}
